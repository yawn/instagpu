package detect

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/pkg/errors"
	probing "github.com/prometheus-community/pro-bing"
)

type Region struct {
	Endpoint string
	Latency  struct {
		Avg uint64 `json:"avg"`
		Min uint64 `json:"min"`
		Max uint64 `json:"max"`
	} `json:"latency"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

func (r *Region) MeasureLatency(ctx context.Context) error {

	slog.Debug("measuring latency",
		slog.String("region", r.Name),
	)

	pinger, err := probing.NewPinger(r.Endpoint)

	if err != nil {
		return errors.Wrapf(err, "failed to setup pinging region %q", r.Name)
	}

	pinger.Count = 3

	if err := pinger.RunWithContext(ctx); err != nil {
		return errors.Wrapf(err, "failed to execute pinging region %q", r.Name)
	}

	stats := pinger.Statistics()

	r.Latency.Avg = uint64(stats.AvgRtt.Milliseconds())
	r.Latency.Max = uint64(stats.MaxRtt.Milliseconds())
	r.Latency.Min = uint64(stats.MinRtt.Milliseconds())

	return nil

}

func (r *Region) String() string {

	var b strings.Builder

	fmt.Fprintf(&b, "📍 %s-%s", r.Provider, r.Name)
	fmt.Fprintf(&b, "\t🐢 %dms", r.Latency.Avg)

	return b.String()

}
