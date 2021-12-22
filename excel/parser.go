package excel

import (
	"strings"

	"github.com/victorjacobs/csv2ynab/util"
	"github.com/victorjacobs/csv2ynab/ynab"
	"github.com/xuri/excelize/v2"
)

func Convert(filePath string) ([]ynab.Transaction, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := f.Rows("Volgende uittreksel")
	if err != nil {
		return nil, err
	}

	var transactions []ynab.Transaction
	headerEncountered := false
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

		// Skip the payment of the bill
		if strings.HasPrefix(row[1], "BETALING VIA DOMICILIERIN") {
			continue
		}

		transaction := ynab.Transaction{
			Date:   row[0],
			Payee:  sanitizePayee(row[1]),
			Amount: util.SanitizeAmount(row[4]),
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
