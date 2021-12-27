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
	list := flag.Bool("list", false, "lists all budgets and accounts on the given API key")
	flag.Parse()

	if *list {
		ynabClient, err := ynab.NewClient(config)
		if err != nil {
			log.Fatal(err)
		}

		budgets, err := ynabClient.GetBudgets()
		if err != nil {
			log.Fatal(err)
		}

		for _, budget := range budgets {
			fmt.Printf("Budget %q (%v)\n", budget.Name, budget.Id)
			accounts, err := ynabClient.GetAccounts(budget.Id)
			if err != nil {
				log.Fatal(err)
			}

			for _, account := range accounts {
				if account.Deleted || account.Closed {
					continue
				}

				fmt.Printf("\tAccount %q (%v)\n", account.Name, account.Id)
			}

			fmt.Print("\n")
		}

		return
	}

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
