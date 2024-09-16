package detect

// Devices maps known vendor-device tags to FP32 TFLOP
var Devices = map[string]float32{ // NOTE: help with validation always welcome
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
