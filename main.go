package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/csv"
	"github.com/victorjacobs/csv2ynab/excel"
	"github.com/victorjacobs/csv2ynab/model"
	"github.com/victorjacobs/csv2ynab/ynab"
)

func main() {
	// Load config
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading .ynabrc: %+v", err)
	}

	// Flags
	filePath := flag.String("in", "", "path to input file [required]")
	outputPath := flag.String("out", "", "path to output CSV to, if not set transactions will be sent through the YNAB API")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Error: Input path is required")
		flag.Usage()
		return
	}

	// Parsing
	var transactions []model.Transaction
	if strings.Contains(*filePath, "csv") {
		transactions, err = csv.Convert(*filePath)
	} else if strings.Contains(*filePath, "xlsx") {
		transactions, err = excel.Convert(*filePath)
	} else {
		log.Fatalf("Input path %v not recognized", *filePath)
	}

	log.Printf("Parsed %v transactions", len(transactions))

	if err != nil {
		log.Fatal(err)
	}

	// Writing
	if *outputPath != "" {
		log.Print("Writing result to CSV file")
		csv.Write(*outputPath, transactions)
	} else {
		log.Print("Importing file into YNAB directly")
		err = ynab.Write(config, transactions)
	}

	if err != nil {
		log.Fatal(err)
	}
}
