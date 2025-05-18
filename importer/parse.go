package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/victorjacobs/csv2ynab/csv"
	"github.com/victorjacobs/csv2ynab/excel"
	"github.com/victorjacobs/csv2ynab/model"
)

func ParseFromPath(path string) ([]*model.Transaction, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(filepath.Base(path), fileBytes)
}

func Parse(fileName string, data []byte) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	var err error
	if strings.Contains(fileName, "csv") {
		transactions, err = csv.Convert(data)
	} else if strings.Contains(fileName, "xlsx") {
		transactions, err = excel.Convert(data)
	} else {
		return nil, fmt.Errorf("input path %v not recognized", fileName)
	}

	if err != nil {
		return nil, err
	}

	log.Printf("Parsed %v transactions", len(transactions))

	return transactions, nil
}
