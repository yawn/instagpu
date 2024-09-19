package database

import (
	"context"

	"github.com/yawn/instagpu/detect"
)

type Provider interface {
	Instances(ctx context.Context, region *detect.Region) ([]*detect.Instance, error)
	Name() string
	Prices(ctx context.Context, region *detect.Region, instance *detect.Instance) (*detect.Prices, error)
	Regions(ctx context.Context) ([]*detect.Region, error)
}
