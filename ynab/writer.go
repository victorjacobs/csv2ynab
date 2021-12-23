package ynab

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/model"
)

func Write(config config.Config, transactions []model.Transaction) error {
	// Interactively ask the user which budget and account to import into
	client, err := NewClient(config)
	if err != nil {
		return err
	}

	budgets, err := client.GetBudgets()
	if err != nil {
		return err
	}

	var budgetNames []string
	for _, budget := range budgets {
		budgetNames = append(budgetNames, budget.Name)
	}

	prompt := promptui.Select{
		Label: "Select budget to import into",
		Items: budgetNames,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	selectedBudget := budgets[i]

	accounts, err := client.GetAccounts(selectedBudget.Id)
	if err != nil {
		return err
	}

	var accountNames []string
	for _, account := range accounts {
		accountNames = append(accountNames, account.Name)
	}

	prompt = promptui.Select{
		Label: "Select account to import into",
		Items: accountNames,
	}

	i, _, err = prompt.Run()
	if err != nil {
		return err
	}

	selectedAccount := accounts[i]

	// prompt = promptui.Select{
	// 	Label: "Have the transactions all cleared?",
	// 	Items: []string{"Yes", "No"},
	// }

	// i, _, err = prompt.Run()
	// if err != nil {
	// 	return err
	// }

	// cleared := i == 0

	confirm := promptui.Prompt{
		Label:     fmt.Sprintf("Import transactions into %q account in %q", selectedAccount.Name, selectedBudget.Name),
		IsConfirm: true,
	}

	result, err := confirm.Run()
	if err != nil {
		return err
	}

	if result != "y" {
		return nil
	}

	// Convert transactions from internal representation to YNAB models
	var transactionsForPost []transaction
	for _, transaction := range transactions {
		transactionsForPost = append(transactionsForPost, transactionFromModel(transaction, selectedAccount.Id, false))
	}

	err = client.PostTransactions(selectedBudget.Id, selectedAccount.Id, transactionsForPost)
	if err != nil {
		return err
	}

	return nil
}
