package ynab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/victorjacobs/csv2ynab/config"
)

const apiEndpoint = "https://api.youneedabudget.com/v1"

type Client struct {
	config config.Config
}

func NewClient(config config.Config) Client {
	return Client{
		config: config,
	}
}

func (c *Client) PostTransactions(budgetId string, accountId string, transactions []transaction) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/budgets/%s/transactions", apiEndpoint, budgetId)

	requestBody := postTransactionsRequest{
		Transactions: transactions,
	}
	requestBodyMarshalled, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := c.createRequest("POST", url, requestBodyMarshalled)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 201 {
		return fmt.Errorf("posting transactions returned status %v", resp.StatusCode)
	}

	response := postTransactionsResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if len(response.Data.DuplicateImportIds) != 0 {
		log.Printf("File contained %v duplicates", len(response.Data.DuplicateImportIds))
	}

	return nil
}

func (c *Client) GetBudgets() ([]Budget, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/budgets", apiEndpoint)

	req, err := c.createRequest("GET", url, []byte{})
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting budgets returned status %v", resp.StatusCode)
	}

	response := getBudgetsResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Data.Budgets, nil
}

func (c *Client) GetAccounts(budgetId string) ([]Account, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/budgets/%s/accounts", apiEndpoint, budgetId)

	req, err := c.createRequest("GET", url, []byte{})
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting accounts returned status %v", resp.StatusCode)
	}

	response := getAccountsResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Data.Accounts, nil
}

func (c *Client) createRequest(method string, url string, json []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.Ynab.ApiKey))

	return req, nil
}
