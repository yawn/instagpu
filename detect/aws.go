//xgo:build aws

package detect

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
)

type AWS struct {
	cfg aws.Config
}

func NewAWS(cfg aws.Config) *AWS {
	return &AWS{
		cfg: cfg,
	}
}

func (a *AWS) Regions(ctx context.Context) (Regions, error) {

	const PROVIDER = "aws"

	client := ec2.NewFromConfig(a.cfg)

	res, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		Filters: []types.Filter{
			{
				Name: aws.String("opt-in-status"),
				Values: []string{
					"opt-in-not-required",
					"opted-in",
				},
			},
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to enumerate regions")
	}

	var regions Regions

	for _, region := range res.Regions {

		regions = append(regions, &Region{
			Name:     *region.RegionName,
			endpoint: *region.Endpoint,
			Provider: PROVIDER,
		})

	}

	if err := regions.MeasureLatency(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed to measure region latency")
	}

	return regions, nil

}
