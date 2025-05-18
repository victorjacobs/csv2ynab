package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/ynab"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List accounts on YNAB account",
	RunE: func(cmd *cobra.Command, args []string) error {
		ynabClient, err := ynab.NewClient(cfg.YNAB)
		if err != nil {
			return fmt.Errorf("failed to create YNAB client: %w", err)
		}

		budgets, err := ynabClient.GetBudgets()
		if err != nil {
			return fmt.Errorf("failed to load budgets: %w", err)
		}

		for _, budget := range budgets {
			fmt.Printf("Budget %q (%v)\n", budget.Name, budget.Id)
			accounts, err := ynabClient.GetAccounts(budget.Id)
			if err != nil {
				return fmt.Errorf("failed to get accounts: %w", err)
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
