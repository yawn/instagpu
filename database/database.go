package database

import (
	"container/heap"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/yawn/spottty/detect"
	"golang.org/x/sync/errgroup"
)

type Database []*detect.Prices

func New(ctx context.Context, providers ...Provider) (Database, error) {

	wg, ctx := errgroup.WithContext(ctx)

	var (
		logger  = slog.Default()
		mutex   sync.Mutex
		results []*detect.Prices
	)

	for _, provider := range providers {

		logger = logger.With(
			slog.String("provider", provider.Name()),
		)

		regions, err := provider.Regions(ctx)

		if err != nil {
			return nil, err
		}

		for _, region := range regions {

			logger := logger.With(
				slog.String("region", region.Name),
			)

			wg.Go(func() error {

				logger.Debug("measuring latency for region")

				return region.MeasureLatency(ctx)

			})

			wg.Go(func() error {

				logger.Debug("gathering instances")

				instances, err := provider.Instances(ctx, region)

				if err != nil {
					return err
				}

				for _, instance := range instances {

					logger := logger.With(
						slog.String("instance", instance.Name),
					)

					wg.Go(func() error {

						logger.Debug("gathering prices")

						prices, err := provider.Prices(ctx, region, instance)

						if err != nil {
							return err
						}

						if prices == nil {
							logger.Warn("instance in practice not available as spot")
							return nil
						}

						mutex.Lock()
						defer mutex.Unlock()

						results = append(results, prices)

						return nil

					})

				}

				return nil

			})

		}

	}

	logger.Debug("waiting for results")

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return results, nil

}

func (d Database) Filter(max int, filters ...Filter) []*Result {

	var (
		q       = make(queue, 0)
		results []*Result
	)

	for _, prices := range d {

		e := &Result{
			Prices: prices,
			Score:  prices.PTGPIndex(),
		}

		heap.Push(&q, e)

	}

	for _, prices := range q {
		prices.Relative = prices.Score / q[0].Score
	}

	for _, result := range q {

		var skip bool

		for _, filter := range filters {

			if ok := filter(result.Prices); !ok {
				skip = true
				break
			}

		}

		if !skip {

			results = append(results, result)

			if len(results) == max {
				break
			}

		}

	}

	return results

}

func (d Database) Save(path string) error {

	fh, err := os.Create(path)

	if err != nil {
		return errors.Wrapf(err, "failed to save to file %q", path)
	}

	defer fh.Close()

	enc := json.NewEncoder(fh)
	enc.SetIndent("", "\t")

	return enc.Encode(d)

}

func Load(path string) (Database, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %q", path)
	}

	defer fh.Close()

	dec := json.NewDecoder(fh)

	var db Database

	if err := dec.Decode(&db); err != nil {
		return nil, errors.Wrapf(err, "corrupt database in file %q", path)
	}

	return db, nil

}
