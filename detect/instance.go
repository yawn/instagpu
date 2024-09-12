package detect

import "fmt"

type Instance struct {
	GPU struct {
		Count  uint
		Memory uint64
		Name   string
		Vendor string
	}
	Arch       string
	ClockSpeed float32
	Count      uint
	Memory     uint64
	Name       string
	Network    float64
	Vendor     string
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s: CPU %d x %s-%s@%.2fGHz; GPU %dGiB plus %dx%s-%s with %dGB; Network: %.2f GiBs)",
		i.Name,
		i.Count,
		i.Vendor,
		i.Arch,
		i.ClockSpeed,
		i.Memory/1024,
		i.GPU.Count,
		i.GPU.Vendor,
		i.GPU.Name,
		i.GPU.Memory/1024,
		i.Network/8,
	)
}
