package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/victorjacobs/csv2ynab/model"
	"github.com/victorjacobs/csv2ynab/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Convert(data []byte) ([]*model.Transaction, error) {
	fileContents := string(data)

	// Replace any \r with \n
	fileContents = strings.ReplaceAll(fileContents, "\r", "\n")

	// Create csv fileReader
	fileReader := strings.NewReader(fileContents)
	csvReader := csv.NewReader(fileReader)

	// Guess record separator by looking at header
	headerLine := strings.Split(fileContents, "\n")[0]
	if strings.Contains(headerLine, ",") {
		csvReader.Comma = ','
	} else if strings.Contains(headerLine, ";") {
		csvReader.Comma = ';'
	} else {
		return nil, errors.New("could not determine CSV separator")
	}

	// Get header
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// Find indices for columns
	columnIndices := columnIndices(header, dateColumns, payeeColums, descriptionColumns, outflowColumns, inflowColumns, amountColumns, feeColumns)
	dateIndex := columnIndices[0]
	payeeIndex := columnIndices[1]
	descriptionIndex := columnIndices[2]
	outflowIndex := columnIndices[3]
	inflowIndex := columnIndices[4]
	amountIndex := columnIndices[5]
	feeIndex := columnIndices[6]

	if dateIndex == -1 || payeeIndex == -1 {
		return nil, fmt.Errorf("input file not valid")
	}

	caser := cases.Title(language.English)
	var transactions []*model.Transaction

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		payee := strings.TrimSpace(record[payeeIndex])

		var description string
		if descriptionIndex != -1 {
			description = record[descriptionIndex]

			if p, present := payeeFromDescription(description); payee == "" && present {
				payee = p
			}
		}

		var amount float64
		if outflowIndex != -1 && inflowIndex != -1 {
			outflow := util.SanitizeAmount(strings.TrimSpace(record[outflowIndex]))
			inflow := util.SanitizeAmount(strings.TrimSpace(record[inflowIndex]))

			if outflow != 0 {
				amount = outflow
			} else if inflow != 0 {
				amount = inflow
			}
		} else if amountIndex != -1 {
			amount = util.SanitizeAmount(strings.TrimSpace(record[amountIndex]))
		}

		if feeIndex != -1 {
			fee := util.SanitizeAmount(strings.TrimSpace(record[feeIndex]))
			amount = amount - fee
		}

		date, err := parseDate(record[dateIndex])
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &model.Transaction{
			Date:        date,
			Payee:       caser.String(strings.ToLower(payee)),
			Amount:      amount,
			Description: description,
		})
	}

	return transactions, nil
}

func parseDate(s string) (time.Time, error) {
	if strings.Contains(s, "/") {
		return time.Parse("02/01/2006", s)
	} else if strings.Contains(s, "-") {
		// Remove time from the thing
		dateTimeParts := strings.Split(s, " ")

		return time.Parse("2006-01-02", dateTimeParts[0])
	} else {
		return time.Time{}, errors.New("failed to determine date format")
	}
}

func payeeFromDescription(s string) (string, bool) {
	if strings.Contains(s, "KBC-PLUSREKENING") {
		return "KBC", true
	}

	if strings.HasPrefix(s, "BETALING VIA") {
		return util.StringBetween(s, "UUR", "MET")
	}

	if strings.HasPrefix(s, "EUROPESE DOMICILIERING") {
		return util.StringBetween(s, "SCHULDEISER     : ", "REF. ")
	}

	return "", false
}

func columnIndices(header []string, columnNamesList ...[]string) []int {
	indices := make([]int, len(columnNamesList))
	for i := range indices {
		indices[i] = -1
	}

	for columnIndex, headerValue := range header {
		for columnNamesIndex, columnNames := range columnNamesList {
			for _, columnName := range columnNames {
				if headerValue == columnName {
					indices[columnNamesIndex] = columnIndex
				}
			}
		}
	}

	return indices
}
