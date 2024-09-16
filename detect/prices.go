package detect

import (
	"fmt"
	"strings"
)

type Prices struct {
	AvailablityZones uint      `json:"availability_zones"`
	Avg              float64   `json:"avg"`
	Instance         *Instance `json:"instance"`
	Max              float64   `json:"max"`
	Min              float64   `json:"min"`
}

func (p *Prices) PTGPIndex() float64 {

	perf := p.Instance.GPU.FP32

	if perf == nil {
		panic("cannot determine price-to-gpu-performance index without GPU performance data")
	}

	return *perf / p.Avg

}

func (p *Prices) String() string {

	var b strings.Builder

	b.WriteString(p.Instance.String())

	fmt.Fprintf(&b, "💰 %.2f USD/h", p.Avg)
	fmt.Fprintf(&b, "\t▼ %.2f USD/h", p.Min)
	fmt.Fprintf(&b, "\t▲ %.2f USD/h", p.Max)
	fmt.Fprintf(&b, "\t🧩 %d AZs", p.AvailablityZones)

	return b.String()

}
