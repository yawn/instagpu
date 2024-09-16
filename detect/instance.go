package detect

import (
	"fmt"
	"log/slog"
)

type Instance struct {
	Arch       string
	ClockSpeed float32
	Count      uint
	GPU        struct {
		Count  uint
		FP32   *float32 // TFLOPS performance
		Memory uint64
		Name   string
		Vendor string
	}
	Memory  uint64
	Name    string
	Network float64
	Region  *Region
	Vendor  string
}

func (i *Instance) MeasureTFLOPS() {

	tag := fmt.Sprintf("%s-%s", i.GPU.Vendor, i.GPU.Name)

	tflops, ok := Devices[tag]

	if !ok {

		slog.Warn("no device data for vendor tag - please consider contributing a pull-request",
			slog.String("tag", tag),
		)

	} else {

		tflops = tflops * float32(i.GPU.Count)
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
