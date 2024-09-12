package detect

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type Regions []*Region

func (r Regions) MeasureLatency(ctx context.Context) error {

	wg, ctx := errgroup.WithContext(ctx)

	for _, region := range r {

		wg.Go(func() error {

			slog.Debug("measure latency of region",
				slog.String("region", region.Name),
			)

			return region.measureLatency(ctx)

		})

	}

	return wg.Wait()

}
