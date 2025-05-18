package importer

import (
	log "github.com/sirupsen/logrus"

	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/csv"
	"github.com/victorjacobs/csv2ynab/ynab"
)

func ProcessFile(config *config.YNAB, inputPath string, outputPath string) error {
	transactions, err := ParseFromPath(inputPath)
	if err != nil {
		return err
	}

	// Write
	if outputPath != "" {
		log.Print("Writing result to CSV file")
		csv.Write(outputPath, transactions)
	} else {
		log.Print("Importing file into YNAB directly")
		err = ynab.Write(config, transactions)
	}

	return err
}
