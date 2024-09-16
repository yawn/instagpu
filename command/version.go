package command

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionVersion string

var versionCmd = &cobra.Command{

	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s", app),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s (%s)\n", app, versionVersion)
		return nil
	},
}

func init() {

	buildInfo, ok := debug.ReadBuildInfo()

	if !ok {
		panic("missing build information")
	}

	var rev string

	for _, setting := range buildInfo.Settings {

		if setting.Key == "vcs.revision" {
			rev = setting.Value
			break
		}

	}

	if rev == "" {
		rev = "???????"
	}

	versionVersion = rev[:7]

	rootCmd.AddCommand(versionCmd)

}
