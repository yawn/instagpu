package detect

import "fmt"

type Prices struct {
	AvailablityZones uint
	Avg              float64
	Instance         *Instance
	Max              float64
	Min              float64
}

func (p *Prices) String() string {
	return fmt.Sprintf("%s %.2f USD/h (%.2f USD/h <-> %.2f USD/h) over %d AZs",
		p.Instance.String(),
		p.Avg,
		p.Min,
		p.Max,
		p.AvailablityZones,
	)
}
