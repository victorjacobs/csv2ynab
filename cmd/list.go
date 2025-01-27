package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/ynab"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List accounts on YNAB account",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := rootCmd.PersistentFlags().GetString("config")
		if err != nil {
			return err
		}

		config, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		ynabClient, err := ynab.NewClient(config.Ynab)
		if err != nil {
			log.Fatal(err)
		}

		budgets, err := ynabClient.GetBudgets()
		if err != nil {
			log.Fatal(err)
		}

		for _, budget := range budgets {
			fmt.Printf("Budget %q (%v)\n", budget.Name, budget.Id)
			accounts, err := ynabClient.GetAccounts(budget.Id)
			if err != nil {
				log.Fatal(err)
			}

			for _, account := range accounts {
				if account.Deleted || account.Closed {
					continue
				}

				fmt.Printf("\tAccount %q (%v)\n", account.Name, account.Id)
			}

			fmt.Print("\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
