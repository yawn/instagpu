package aws

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
	"github.com/yawn/spottty/detect"
)

const NAME = "aws"

type AWS struct {
	cfg aws.Config
}

func DefaultConfig(ctx context.Context) (aws.Config, error) {

	var (
		region  = os.Getenv("AWS_REGION")
		profile = os.Getenv("AWS_PROFILE")
	)

	slog.Info("configuring aws provider",
		slog.String("region", region),
		slog.String("profile", profile),
	)

	return config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)

}

func New(ctx context.Context) (*AWS, error) {

	cfg, err := DefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	return NewWithConfig(cfg), nil

}

func NewWithConfig(cfg aws.Config) *AWS {
	return &AWS{
		cfg: cfg,
	}
}

func (a *AWS) clientForRegion(region *detect.Region) *ec2.Client {

	cfg := a.cfg.Copy()
	cfg.Region = region.Name

	return ec2.NewFromConfig(cfg)

}

func (a *AWS) Instances(ctx context.Context, region *detect.Region) ([]*detect.Instance, error) {

	client := a.clientForRegion(region)

	var instances []*detect.Instance

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

			if len(e.NetworkInfo.NetworkCards) == 0 {
				panic("unexpected length of NetworkInfo.NetworkCards")
			}

			instance := &detect.Instance{
				Arch:       string(e.ProcessorInfo.SupportedArchitectures[0]),
				Count:      uint(*e.VCpuInfo.DefaultCores),
				ClockSpeed: float64(*e.ProcessorInfo.SustainedClockSpeedInGhz),
				Memory:     uint64(*e.MemoryInfo.SizeInMiB),
				Name:       string(e.InstanceType),
				Network:    float64(*e.NetworkInfo.NetworkCards[0].PeakBandwidthInGbps),
				Region:     region,
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

func (a *AWS) Name() string {
	return NAME
}

func (a *AWS) Regions(ctx context.Context) ([]*detect.Region, error) {

	client := ec2.NewFromConfig(a.cfg)

	var regions []*detect.Region

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

	for _, region := range res.Regions {

		regions = append(regions, &detect.Region{
			Name:     *region.RegionName,
			Endpoint: *region.Endpoint,
			Provider: a.Name(),
		})

	}

	return regions, nil

}

func (a *AWS) Prices(ctx context.Context, region *detect.Region, instance *detect.Instance) (*detect.Prices, error) {

	// looking back a week
	const WINDOW = -24 * 7

	client := a.clientForRegion(region)

	var (
		azs    = make(map[string]any)
		prices []float64
	)

	window := time.Now().Add(time.Duration(WINDOW) * time.Hour)

	paginator := ec2.NewDescribeSpotPriceHistoryPaginator(client, &ec2.DescribeSpotPriceHistoryInput{
		InstanceTypes: []types.InstanceType{
			types.InstanceType(instance.Name),
		},
		ProductDescriptions: []string{"Linux/UNIX"},
		StartTime:           &window,
	})

	for paginator.HasMorePages() {

		res, err := paginator.NextPage(ctx)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to enumerate instances prices")
		}

		for _, e := range res.SpotPriceHistory {

			azs[*e.AvailabilityZone] = struct{}{}

			p, err := strconv.ParseFloat(*e.SpotPrice, 64)

			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse price %q", *e.SpotPrice)
			}

			prices = append(prices, p)

		}

	}

	if len(prices) > 0 {

		price := &detect.Prices{
			AvailablityZones: uint(len(azs)),
			Instance:         instance,
		}

		var avg float64

		for _, price := range prices {
			avg += price
		}

		price.Avg = avg / float64(len(prices))

		slices.Sort(prices)

		price.Min = prices[0]
		price.Max = prices[len(prices)-1]

		return price, nil

	}

	return nil, nil

}
