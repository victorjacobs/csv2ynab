package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/importer"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watches directories for files to import",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := rootCmd.PersistentFlags().GetString("config")
		if err != nil {
			return err
		}

		config, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		importer.Watch(config)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
