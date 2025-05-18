package cmd

import (
	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/importer"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watches directories for files to import",
	RunE: func(cmd *cobra.Command, args []string) error {
		importer.Watch(cfg)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
