package detect

import "fmt"

type Prices struct {
	AvailablityZones uint      `json:"availability_zones"`
	Avg              float64   `json:"avg"`
	Instance         *Instance `json:"instance"`
	Max              float64   `json:"max"`
	Min              float64   `json:"min"`
}
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
