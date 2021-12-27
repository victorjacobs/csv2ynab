package excel

import (
	"fmt"
	"strings"
	"time"

	"github.com/victorjacobs/csv2ynab/model"
	"github.com/victorjacobs/csv2ynab/util"
	"github.com/xuri/excelize/v2"
)

func Convert(filePath string) ([]model.Transaction, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := f.Rows("Volgende uittreksel")
	if err != nil {
		return nil, err
	}

	var transactions []model.Transaction
	headerEncountered := false
	rowNumber := 0
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		if len(row) == 0 {
			continue
		}

		if !headerEncountered {
			if row[0] == "Datum" {
				headerEncountered = true
			}

			continue
		}

		rowNumber += 1

		// Skip the payment of the bill
		if strings.HasPrefix(row[1], "BETALING VIA DOMICILIERIN") {
			continue
		}

		date, err := time.Parse("02/01/2006", row[0])
		if err != nil {
			return nil, err
		}

		transaction := model.Transaction{
			Date:           date,
			Payee:          sanitizePayee(row[1]),
			Amount:         util.SanitizeAmount(row[4]),
			IdempotencyKey: fmt.Sprintf("%v", rowNumber),
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func sanitizePayee(payee string) string {
	// Remove some Curve cruft
	cleaned := strings.ReplaceAll(payee, "CRV*", "")
	cleaned = strings.ReplaceAll(cleaned, "Vilnius", "")
	cleaned = strings.TrimSpace(cleaned)
	return strings.Title(strings.ToLower(cleaned))
}
