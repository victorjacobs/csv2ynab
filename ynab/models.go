package ynab

import "github.com/victorjacobs/csv2ynab/model"

// YNAB-specific version of a transaction
type transaction struct {
	Date      string `json:"date"`
	PayeeName string `json:"payee_name"`
	Memo      string `json:"memo"`
	Amount    int    `json:"amount"`
	AccountId string `json:"account_id"`
	Cleared   string `json:"cleared,omitempty"`
	ImportId  string `json:"import_id,omitempty"`
}

func transactionFromModel(t model.Transaction, accountId string, cleared bool) transaction {
	// TODO add cleared
	amount := int(t.Amount * 1000)
	date := t.Date.Format("2006-01-02")
	var memo string
	if len(t.Memo) > 200 {
		memo = t.Memo[0:200]
	} else {
		memo = t.Memo
	}

	return transaction{
		Date:      date,
		PayeeName: t.Payee,
		Memo:      memo,
		Amount:    amount,
		AccountId: accountId,
	}
}

type Budget struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Account struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Private models to parse API responses
type postTransactionsRequest struct {
	Transactions []transaction `json:"transactions"`
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
