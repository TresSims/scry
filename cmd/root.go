package cmd

import (
	"os"

	"charm.land/log/v2"
	"github.com/TresSims/scry/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "scry",
	Short:             "A server health dashboard served over the ssh protocol.",
	PersistentPreRunE: configureLogging,
}

// Execute runs the rootCmd and handles logging and failing on errors
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("Failed to start scry", "error", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)

	config.ConfigureFlags(rootCmd)
}

func configureLogging(_ *cobra.Command, _ []string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	log.SetLevel(log.Level(cfg.Verbosity))

	log.Debug("config calculated", "values", cfg)

	return nil
}
