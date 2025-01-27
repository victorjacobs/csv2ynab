package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/importer"
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Imports a file into YNAB",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		outputPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		configPath, err := rootCmd.PersistentFlags().GetString("config")
		if err != nil {
			return err
		}

		config, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return importer.ProcessFile(config.Ynab, filePath, outputPath)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("output", "o", "", "file to write CSV to, instead of importing directly to YNAB")
}
