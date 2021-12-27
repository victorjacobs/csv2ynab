package command

import (
	"fmt"
	"log"

	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/ynab"
)

func ListBudgetsAndAccounts(config config.Ynab) {
	ynabClient, err := ynab.NewClient(config)
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
}
