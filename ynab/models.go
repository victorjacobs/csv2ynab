package ynab

import (
	"github.com/victorjacobs/csv2ynab/model"
)

// YNAB-specific version of a transaction
type transaction struct {
	Date      string `json:"date"`
	PayeeName string `json:"payee_name"`
	Memo      string `json:"memo,omitempty"`
	Amount    int    `json:"amount"`
	AccountId string `json:"account_id"`
	Cleared   string `json:"cleared,omitempty"`
	ImportId  string `json:"import_id,omitempty"`
}

func transactionFromModel(t model.Transaction, accountId string, cleared bool) transaction {
	amount := int(t.Amount * 1000)
	date := t.Date.Format("2006-01-02")
	var memo string
	if len(t.Description) > 200 {
		memo = t.Description[0:200]
	} else {
		memo = t.Description
	}

	return transaction{
		Date:      date,
		PayeeName: t.Payee,
		Memo:      memo,
		Amount:    amount,
		AccountId: accountId,
		ImportId:  t.ImportId(),
	}
}

type Budget struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Account struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Closed  bool   `json:"closed"`
	Deleted bool   `json:"deleted"`
}

// Private models to parse API responses
type postTransactionsRequest struct {
	Transactions []transaction `json:"transactions"`
}

type postTransactionsResponse struct {
	Data postTransactionsResponseData `json:"data"`
}

type postTransactionsResponseData struct {
	DuplicateImportIds []string `json:"duplicate_import_ids"`
}

type getBudgetsResponse struct {
	Data getBudgetsResponseData `json:"data"`
}

type getBudgetsResponseData struct {
	Budgets []Budget `json:"budgets"`
}

type getAccountsResponse struct {
	Data getAccountsResponseData `json:"data"`
}

type getAccountsResponseData struct {
	Accounts []Account `json:"accounts"`
}
