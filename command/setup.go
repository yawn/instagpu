package command

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yawn/instagpu/provider"
	"github.com/yawn/instagpu/provider/aws"
	"golang.org/x/sync/errgroup"
)

var setupProviderAWS bool
var setupTimeout time.Duration

var setupCmd = &cobra.Command{

	Use:   "setup",
	Short: "Setup foundational infrastructure for the selected provider(s)",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, _ := context.WithTimeout(context.Background(), setupTimeout)

		var providers []provider.Provider

		if setupProviderAWS {

			provider, err := aws.New(ctx)

			if err != nil {
				return errors.Wrapf(err, "failed to configure aws")
			}

			providers = append(providers, provider)

		}

		if len(providers) == 0 {
			return fmt.Errorf("no providers selected")
		}

		wg, ctx := errgroup.WithContext(ctx)

		for _, provider := range providers {
			wg.Go(func() error {
				return provider.Setup(ctx)
			})
		}

		return wg.Wait()

	},
}

func init() {

	flags := setupCmd.Flags()

	flags.BoolVar(&setupProviderAWS, "provider-aws", true, "Enable AWS")
	flags.DurationVar(&setupTimeout, "timeout", 5*time.Minute, "Timeout for all API operations")

	rootCmd.AddCommand(setupCmd)

}
