package detect

import (
	"fmt"
)

type Instance struct {
	Arch       string  `json:"arch"`
	ClockSpeed float64 `json:"clock_speed"`
	Count      uint    `json:"count"`
	GPU        *GPU    `json:"gpu"`
	Memory     uint64  `json:"memory"`
	Name       string  `json:"name"`
	Network    float64 `json:"network"`
	Region     *Region `json:"region"`
	Vendor     string  `json:"vendor"`
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s 🏷️ %s ⚙️ %d x %s-%s ⚡ %.2fGHz 🧠 %dGiB %s 🌐 %.2f GiBs",
		i.Region.String(),
		i.Name,
		i.Count,
		i.Vendor,
		i.Arch,
		i.ClockSpeed,
		i.Memory/1024,
		i.GPU.String(),
		i.Network/8,
	)
}
