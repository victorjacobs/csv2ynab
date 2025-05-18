package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/config"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "csv2ynab",
	Short: "Import CSV-like files into YNAB",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		c, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfg = c

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config file")
}
