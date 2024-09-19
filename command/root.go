package command

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

const app = "instagpu"

var rootDebug bool

var rootCmd = &cobra.Command{
	Use: app,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		logOptions := &slog.HandlerOptions{}

		if rootDebug {
			logOptions.Level = slog.LevelDebug.Level()
		}

		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, logOptions)).With(
			slog.String("version", versionVersion),
		))

	},
}

func init() {

	flags := showCmd.Flags()

	flags.BoolVar(&rootDebug, "debug", false, "Enable debug logging")

}

func Run() error {
	return rootCmd.Execute()
}
