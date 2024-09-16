package detect

import (
	"fmt"
	"log/slog"
)

type Instance struct {
	Arch       string  `json:"arch"`
	ClockSpeed float64 `json:"clock_speed"`
	Count      uint    `json:"count"`
	GPU        struct {
		Count  uint     `json:"count"`
		FP32   *float64 `json:"fp32 "` // TFLOPS performance
		Memory uint64   `json:"memory"`
		Name   string   `json:"name"`
		Vendor string   `json:"vendor"`
	} `json:"gpu"`
	Memory  uint64  `json:"memory"`
	Name    string  `json:"name"`
	Network float64 `json:"network"`
	Region  *Region `json:"region"`
	Vendor  string  `json:"vendor"`
}

func (i *Instance) MeasureTFLOPS() {

	tag := fmt.Sprintf("%s-%s", i.GPU.Vendor, i.GPU.Name)

	tflops, ok := Devices[tag]

	if !ok {

		slog.Warn("no device data for vendor tag - please consider contributing a pull-request",
			slog.String("tag", tag),
		)

	} else {

		tflops = tflops * float64(i.GPU.Count)
		i.GPU.FP32 = &tflops

	}

}

func (i *Instance) String() string {
	return fmt.Sprintf("%s %s: CPU %d x %s-%s@%.2fGHz; GPU %dGiB plus %dx%s-%s with %dGB; Network: %.2f GiBs)",
		i.Region.String(),
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
