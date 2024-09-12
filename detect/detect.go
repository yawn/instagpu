package detect

type Price struct {
	Avg float64
	Min float64
	Max float64
}

type Target struct {
	Availablity  bool
	AZs          int
	InstanceType string
	Latency      struct {
		Avg int64
		Min int64
		Max int64
	}
	Price          *Price
	Region         string
	regionEndpoint string
}
