package csv

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/victorjacobs/csv2ynab/util"
	"github.com/victorjacobs/csv2ynab/ynab"
)

func Convert(filePath string) ([]ynab.Transaction, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fileContents := string(fileBytes)

	// Replace any \r with \n
	fileContents = strings.ReplaceAll(fileContents, "\r", "\n")

	// Create csv reader
	r := csv.NewReader(strings.NewReader(fileContents))
	r.Comma = ';'

	// Get header
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	// Find indices for columns
	dateIndex := -1
	payeeIndex := -1
	memoIndex := -1
	outflowIndex := -1
	inflowIndex := -1
	for i, v := range header {
		if v == "Datum" {
			dateIndex = i
		}

		if v == "Naam tegenpartij" {
			payeeIndex = i
		}

		if v == "Omschrijving" {
			memoIndex = i
		}

		if v == "debet" {
			outflowIndex = i
		}

		if v == "credit" {
			inflowIndex = i
		}
	}

	if dateIndex == -1 || payeeIndex == -1 || memoIndex == -1 || inflowIndex == -1 || outflowIndex == -1 {
		log.Fatal("File not valid")
	}

	var transactions []ynab.Transaction

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		memo := record[memoIndex]
		payee := strings.TrimSpace(record[payeeIndex])

		if p, present := payeeFromMemo(memo); payee == "" && present {
			payee = p
		}

		outflow := util.SanitizeAmount(strings.TrimSpace(record[outflowIndex]))
		inflow := util.SanitizeAmount(strings.TrimSpace(record[inflowIndex]))

		var amount float64
		if outflow != 0 {
			amount = outflow
		} else if inflow != 0 {
			amount = inflow
		}

		transactions = append(transactions, ynab.Transaction{
			Date:   record[dateIndex],
			Payee:  strings.Title(strings.ToLower(payee)),
			Amount: amount,
			Memo:   record[memoIndex],
		})
	}

	return transactions, nil
}

func payeeFromMemo(s string) (string, bool) {
	if strings.Contains(s, "KBC-PLUSREKENING") {
		return "KBC", true
	}

	if strings.HasPrefix(s, "BETALING VIA") {
		return stringBetween(s, "UUR", "MET")
	}

	if strings.HasPrefix(s, "EUROPESE DOMICILIERING") {
		return stringBetween(s, "SCHULDEISER     : ", "REF. ")
	}

	return "", false
}

func stringBetween(s string, start string, end string) (string, bool) {
	split := strings.Split(s, start)
	if len(split) < 2 {
		return "", false
	}
	split = strings.Split(split[1], end)
	if len(split) < 2 {
		return "", false
	}

	return strings.TrimSpace(split[0]), true
}
