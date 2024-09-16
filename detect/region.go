package detect

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	probing "github.com/prometheus-community/pro-bing"
)

type Region struct {
	endpoint string
	Latency  struct {
		Avg int64
		Min int64
		Max int64
	}
	Name     string
	Provider string
}

func (r *Region) MeasureLatency(ctx context.Context) error {

	slog.Debug("measuring latency",
		slog.String("region", r.Name),
	)

	pinger, err := probing.NewPinger(r.endpoint)

	if err != nil {
		return errors.Wrapf(err, "failed to setup pinging region %q", r.Name)
	}

	pinger.Count = 3

	if err := pinger.RunWithContext(ctx); err != nil {
		return errors.Wrapf(err, "failed to execute pinging region %q", r.Name)
	}

	stats := pinger.Statistics()

	r.Latency.Avg = stats.AvgRtt.Milliseconds()
	r.Latency.Max = stats.MaxRtt.Milliseconds()
	r.Latency.Min = stats.MinRtt.Milliseconds()

	return nil

}

func (r *Region) String() string {
	return fmt.Sprintf("%s-%s (%d)", r.Provider, r.Name, r.Latency.Avg)
}
