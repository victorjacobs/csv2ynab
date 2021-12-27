package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/victorjacobs/csv2ynab/command"
	"github.com/victorjacobs/csv2ynab/config"
)

func main() {
	// Flags
	filePath := flag.String("in", "", "path to input file [required]")
	outputPath := flag.String("out", "", "path to output CSV to, if not set transactions will be sent through the YNAB API")
	list := flag.Bool("list", false, "lists all budgets and accounts on the given API key")
	watch := flag.Bool("watch", false, "watches directories for new files to process")
	configPath := flag.String("config", "", "path to configuration")
	flag.Parse()

	// Load config
	config, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Error loading .ynabrc: %+v", err)
	}

	if *list {
		command.ListBudgetsAndAccounts(config.Ynab)
		return
	}

	if *watch {
		command.WatchDirectories(config)
		return
	}

	if *filePath == "" {
		fmt.Println("Error: Input path is required")
		flag.Usage()
		return
	}

	err = command.ProcessFile(config.Ynab, *filePath, *outputPath)
	if err != nil {
		log.Fatal(err)
	}
}
