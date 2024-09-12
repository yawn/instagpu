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

func (a *AWS) Instances(ctx context.Context, region *Region) ([]*Instance, error) {

	cfg := a.cfg.Copy()
	cfg.Region = region.Name

	client := ec2.NewFromConfig(cfg)

	paginator := ec2.NewDescribeInstanceTypesPaginator(client, &ec2.DescribeInstanceTypesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("current-generation"),
				Values: []string{
					"true",
				},
			},
			{
				Name: aws.String("instance-type"),
				Values: []string{
					"g*",
					"p*",
				},
			},
			{
				Name: aws.String("supported-usage-class"),
				Values: []string{
					"spot",
				},
			},
		},
	})

	var instances []*Instance

	for paginator.HasMorePages() {

		res, err := paginator.NextPage(ctx)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to enumerate candidate instances")
		}

		for _, e := range res.InstanceTypes {

			if len(e.ProcessorInfo.SupportedArchitectures) != 1 {
				panic("unexpected length of ProcessorInfo.SupportedArchitectures")
			}

			if e.GpuInfo == nil {
				panic("unexpected empty GpuInfo")
			}

			if len(e.GpuInfo.Gpus) != 1 {
				panic("unexpected length of GpuInfo.Gpus")
			}

			if len(e.NetworkInfo.NetworkCards) != 1 {
				panic("unexpected length of NetworkInfo.NetworkCards")
			}

			instance := &Instance{
				Arch:       string(e.ProcessorInfo.SupportedArchitectures[0]),
				Count:      uint(*e.VCpuInfo.DefaultCores),
				ClockSpeed: float32(*e.ProcessorInfo.SustainedClockSpeedInGhz),
				Memory:     uint64(*e.MemoryInfo.SizeInMiB),
				Name:       string(e.InstanceType),
				Network:    float64(*e.NetworkInfo.NetworkCards[0].PeakBandwidthInGbps),
				Vendor:     *e.ProcessorInfo.Manufacturer,
			}

			instance.GPU.Memory = uint64(*e.GpuInfo.TotalGpuMemoryInMiB)

			gpus := e.GpuInfo.Gpus[0]

			instance.GPU.Count = uint(*gpus.Count)
			instance.GPU.Name = *gpus.Name
			instance.GPU.Vendor = *gpus.Manufacturer

			instances = append(instances, instance)

		}

	}

	return instances, nil

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
