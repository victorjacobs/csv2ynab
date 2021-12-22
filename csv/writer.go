package csv

import (
	"bufio"
	"fmt"
	"os"

	csvencoding "encoding/csv"

	"github.com/victorjacobs/csv2ynab/model"
)

func Write(outputPath string, transactions []model.Transaction) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	w := csvencoding.NewWriter(bufio.NewWriter(outputFile))
	defer w.Flush()

	// Write header
	w.Write([]string{"Date", "Payee", "Inflow", "Memo"})

	for _, transaction := range transactions {
		w.Write([]string{
			transaction.Date.Format("02/01/2006"),
			transaction.Payee,
			fmt.Sprintf("%f", transaction.Amount),
			transaction.Memo,
		})
	}

	return nil
}
