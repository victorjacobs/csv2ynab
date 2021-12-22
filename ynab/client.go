package ynab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	fmt.Printf("%+v", string(requestBodyMarshalled))

	req, err := c.createRequest("POST", url, requestBodyMarshalled)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

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
