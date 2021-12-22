package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	csvencoding "encoding/csv"

	"github.com/victorjacobs/csv2ynab/csv"
	"github.com/victorjacobs/csv2ynab/excel"
	"github.com/victorjacobs/csv2ynab/ynab"
)

func main() {
	flag.Parse()

	args := flag.Args()

	filePath := args[0]
	if filePath == "" {
		log.Fatal("Panik!")
	}

	var err error
	var transactions []ynab.Transaction
	if strings.Contains(filePath, "csv") {
		transactions, err = csv.Convert(filePath)
	} else if strings.Contains(filePath, "xlsx") {
		transactions, err = excel.Convert(filePath)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Write file
	outputFile, err := os.Create("/tmp/test.csv")
	if err != nil {
		log.Fatal(err)
	}
	w := csvencoding.NewWriter(bufio.NewWriter(outputFile))
	defer w.Flush()

	// Write header
	w.Write([]string{"Date", "Payee", "Inflow", "Memo"})

	fmt.Printf("%+v\n", transactions)

	for _, transaction := range transactions {
		w.Write([]string{
			transaction.Date,
			transaction.Payee,
			fmt.Sprintf("%f", transaction.Amount),
			transaction.Memo,
		})
	}
}
