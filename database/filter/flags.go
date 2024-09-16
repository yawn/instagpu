package filter

import (
	"github.com/spf13/pflag"
	"github.com/yawn/spottty/detect"
)

type Flag interface {
	Apply() Filter
	Install(*pflag.FlagSet)
	IsSet() bool
	Name() string
}

var Flags []Flag

func init() {

	gpuMemory := &filterFlag[uint64]{
		description: "Filters by minimum GPU memory in GiB (no default)",
		filter: func(memory uint64) Filter {
			return func(p *detect.Prices) bool {
				return p.Instance.GPU.Memory/1024 >= memory
			}
		},
		install: func(flags *pflag.FlagSet) install[uint64] {
			return flags.Uint64Var
		},
		name: "filter-instance-min-vram",
	}

	gpuTFLOPS := &filterFlag[float64]{
		description: "Filters by minimum GPU TFLOPs (no default)",
		filter: func(tflops float64) Filter {
			return func(p *detect.Prices) bool {

				perf := p.Instance.GPU.FP32

				if perf == nil {
					return false
				}

				return *perf >= tflops

			}
		},
		install: func(flags *pflag.FlagSet) install[float64] {
			return flags.Float64Var
		},
		name: "filter-gpu-min-tflops",
	}

	gpuVendor := &filterFlag[string]{
		description: "Filters by GPU vendor name (no default)",
		filter: func(vendor string) Filter {
			return func(p *detect.Prices) bool {
				return p.Instance.GPU.Vendor == vendor
			}
		},
		install: func(flags *pflag.FlagSet) install[string] {
			return flags.StringVar
		},
		name: "filter-gpu-vendor",
	}

	instanceMemory := &filterFlag[uint64]{
		description: "Filters by minimum compute memory in GiB (no default)",
		filter: func(memory uint64) Filter {
			return func(p *detect.Prices) bool {
				return p.Instance.Memory >= memory*1024
			}
		},
		install: func(flags *pflag.FlagSet) install[uint64] {
			return flags.Uint64Var
		},
		name: "filter-instance-min-ram",
	}

	instancePrice := &filterFlag[float64]{
		description: "Filters by maximum average instance price in USD / h (no default)",
		filter: func(price float64) Filter {
			return func(p *detect.Prices) bool {
				return p.Avg <= price
			}
		},
		install: func(flags *pflag.FlagSet) install[float64] {
			return flags.Float64Var
		},
		name: "filter-instance-max-price",
	}

	regionLatency := &filterFlag[uint64]{
		description: "Filters by maximum region latency (no default)",
		filter: func(latency uint64) Filter {
			return func(p *detect.Prices) bool {
				return p.Instance.Region.Latency.Avg <= latency
			}
		},
		install: func(flags *pflag.FlagSet) install[uint64] {
			return flags.Uint64Var
		},
		name: "filter-region-max-latency",
	}

	Flags = []Flag{
		gpuMemory,
		gpuTFLOPS,
		gpuVendor,
		instanceMemory,
		instancePrice,
		regionLatency,
	}

}
