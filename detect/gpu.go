package detect

import (
	"fmt"
	"log/slog"
	"strings"
)

// devices maps known vendor-device tags to FP32 TFLOP
var devices = map[string]float64{ // NOTE: help with validation always welcome
	"AMD-Radeon Pro V520": 7.373, // https://www.techpowerup.com/gpu-specs/radeon-pro-v520.c3755
	"NVIDIA-A100":         19.49, // https://www.techpowerup.com/gpu-specs/a100-sxm4-40-gb.c3506
	"NVIDIA-A10G":         31.52, // https://www.techpowerup.com/gpu-specs/a10g.c3798
	"NVIDIA-H100":         66.91, // https://www.techpowerup.com/gpu-specs/h100-sxm5-80-gb.c3900
	"NVIDIA-K80":          4.113, // https://www.techpowerup.com/gpu-specs/tesla-k80.c2616
	"NVIDIA-L4":           30.29, // https://www.techpowerup.com/gpu-specs/l4.c4091
	"NVIDIA-L40S":         91.61, // https://www.techpowerup.com/gpu-specs/l40s.c4173
	"NVIDIA-M60":          4.825, // https://www.techpowerup.com/gpu-specs/tesla-m60.c2760
	"NVIDIA-T4":           8.141, // https://www.techpowerup.com/gpu-specs/tesla-t4.c3316
	"NVIDIA-T4g":          8.141, // https://www.techpowerup.com/gpu-specs/tesla-t4g.c4134
	"NVIDIA-V100":         16.35, // https://www.techpowerup.com/gpu-specs/tesla-v100-sxm2-16-gb.c3471
}

type GPU struct {
	Count  uint     `json:"count"`
	FP32   *float64 `json:"fp32 "` // TFLOPS performance
	Memory uint64   `json:"memory"`
	Name   string   `json:"name"`
	Vendor string   `json:"vendor"`
}

func (g *GPU) MeasureTFLOPS() {

	tag := fmt.Sprintf("%s-%s", g.Vendor, g.Name)

	tflops, ok := devices[tag]

	if !ok {

		slog.Warn("no device data for vendor tag - please consider contributing a pull-request",
			slog.String("tag", tag),
		)

	} else {

		tflops = tflops * float64(g.Count)
		g.FP32 = &tflops

	}

}

func (g *GPU) String() string {

	var b strings.Builder

	fmt.Fprintf(&b, "🎨 %dx%s-%s", g.Count, g.Vendor, g.Name)

	if g.FP32 != nil {
		fmt.Fprintf(&b, "\t⚡ %.2f fp32", *g.FP32)
	}

	fmt.Fprintf(&b, "\t🧠 %dGiB", g.Memory/1024)

	return b.String()

}
