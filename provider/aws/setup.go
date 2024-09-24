package aws

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
)

//go:embed cloudformation.yml
var stack string

func (a *AWS) Setup(ctx context.Context) error {

	var (
		client      = cloudformation.NewFromConfig(a.cfg)
		deadline, _ = ctx.Deadline()
		name        = "InstaGPUv1"
		op          = "create"
		outputs     func(context.Context) (*cloudformation.DescribeStacksOutput, error)
		req         = &cloudformation.DescribeStacksInput{
			StackName: &name,
		}
		timeout = deadline.Sub(time.Now())
		tags    = Tags{ // TODO: support custom tags
			"InstaGPU": "v1",
		}
	)

	_, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: &name,
	})

	if err != nil {

		var apiError smithy.APIError

		if errors.As(err, &apiError) && apiError.ErrorMessage() == fmt.Sprintf("Stack with id %s does not exist", name) {
			op = "create"
		} else {
			return err
		}

	} else {
		op = "update"
	}

	if op == "create" {

		slog.Debug("creating stack")

		outputs = func(ctx context.Context) (*cloudformation.DescribeStacksOutput, error) {
			waiter := cloudformation.NewStackCreateCompleteWaiter(client)
			return waiter.WaitForOutput(ctx, req, timeout)
		}

		_, err := client.CreateStack(ctx, &cloudformation.CreateStackInput{
			Capabilities: []types.Capability{
				types.CapabilityCapabilityNamedIam,
			},
			OnFailure:    types.OnFailureDelete,
			StackName:    &name,
			Tags:         tags.ToCF(),
			TemplateBody: &stack,
		})

		if err != nil {
			return errors.Wrapf(err, "failed to create cloudformation stack")
		}

	} else {

		slog.Debug("updating stack")

		outputs = func(ctx context.Context) (*cloudformation.DescribeStacksOutput, error) {
			waiter := cloudformation.NewStackUpdateCompleteWaiter(client)
			return waiter.WaitForOutput(ctx, req, timeout)
		}

		_, err := client.UpdateStack(ctx, &cloudformation.UpdateStackInput{
			Capabilities: []types.Capability{
				types.CapabilityCapabilityNamedIam,
			},
			StackName:    &name,
			Tags:         tags.ToCF(),
			TemplateBody: &stack,
		})

		if err != nil {

			var apiError smithy.APIError

			if errors.As(err, &apiError) && apiError.ErrorMessage() == "No updates are to be performed." {

				slog.Debug("no updates required, skipping")

				outputs = func(ctx context.Context) (*cloudformation.DescribeStacksOutput, error) {
					return client.DescribeStacks(ctx, req)
				}

			} else {
				return errors.Wrapf(err, "failed to create cloudformation stack")
			}

		}

	}

	res, err := outputs(ctx)

	if err != nil {
		return errors.Wrapf(err, "failed to retrieve stack outputs")
	}

	a.instanceProfileARN = *res.Stacks[0].Outputs[0].OutputValue

	slog.Debug("instance profile identified",
		slog.String("arn", a.instanceProfileARN),
	)

	return nil

}
