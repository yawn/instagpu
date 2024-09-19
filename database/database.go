//go:build go1.22
// +build go1.22

package database

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"slices"
	"sync"

	"github.com/pkg/errors"
	"github.com/yawn/instagpu/database/filter"
	"github.com/yawn/instagpu/detect"
	"github.com/yawn/instagpu/provider"
	"golang.org/x/sync/errgroup"
)

type Database []*detect.Prices

func New(ctx context.Context, providers ...provider.Provider) (Database, error) {

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

					instance.GPU.MeasureTFLOPS()

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

func (d Database) Filter(max uint16, filters ...filter.Filter) []*Result {

	var (
		results []*Result
		top     float64
	)

	for _, prices := range d {

		results = append(results, &Result{
			Prices: prices,
			Score:  prices.PTGPIndex(),
		})

	}

	if len(results) == 0 {
		return nil
	}

	slices.SortFunc(results, func(a, b *Result) int {
		return int(b.Score - a.Score)
	})

	top = results[0].Score

	for idx, result := range results {
		result.Index = idx
		result.IndexMax = len(results) - 1
		result.Relative = result.Score / top
	}

	results = slices.DeleteFunc(results, func(result *Result) bool {

		for _, filter := range filters {

			if !filter(result.Prices) {
				return true
			}

		}

		return false

	})

	return results[:min(int(max), len(results))]

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
