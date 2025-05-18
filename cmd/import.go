package cmd

import (
	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/importer"
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Imports a file into YNAB",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		return importer.ProcessFile(cfg.YNAB, filePath, outputPath)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("output", "o", "", "file to write CSV to, instead of importing directly to YNAB")
}
