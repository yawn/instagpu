package database

import "github.com/yawn/spottty/detect"

type Filter func(p *detect.Prices) bool

func GPUVendorName(vendor string) Filter {

	return func(p *detect.Prices) bool {
		return p.Instance.GPU.Vendor == vendor
	}

}

func MaxLatencyInMillis(latency int64) Filter {

	return func(p *detect.Prices) bool {
		return p.Instance.Region.Latency.Avg <= latency
	}

}

func MaxPricePerHourInUSD(price float64) Filter {

	return func(p *detect.Prices) bool {
		return p.Avg <= price
	}

}

func MinComputeMemoryInGiB(memory uint64) Filter {

	return func(p *detect.Prices) bool {
		return p.Instance.Memory >= memory
	}

}
func MinGPUMemoryInGiB(memory uint64) Filter {

	return func(p *detect.Prices) bool {
		return p.Instance.GPU.Memory >= memory
	}

}
