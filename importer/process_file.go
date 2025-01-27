package importer

import (
	"fmt"
	"log"
	"strings"

	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/csv"
	"github.com/victorjacobs/csv2ynab/excel"
	"github.com/victorjacobs/csv2ynab/model"
	"github.com/victorjacobs/csv2ynab/ynab"
)

func ProcessFile(config config.Ynab, inputPath string, outputPath string) error {
	// Parse
	var transactions []model.Transaction
	var err error
	if strings.Contains(inputPath, "csv") {
		transactions, err = csv.Convert(inputPath)
	} else if strings.Contains(inputPath, "xlsx") {
		transactions, err = excel.Convert(inputPath)
	} else {
		return fmt.Errorf("input path %v not recognized", inputPath)
	}

	if err != nil {
		return err
	}

	log.Printf("Parsed %v transactions", len(transactions))

	// Write
	if outputPath != "" {
		log.Print("Writing result to CSV file")
		csv.Write(outputPath, transactions)
	} else {
		log.Print("Importing file into YNAB directly")
		err = ynab.Write(config, transactions)
	}

	if err != nil {
		return err
	}

	return nil
}
