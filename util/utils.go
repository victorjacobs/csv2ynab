package util

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

func SanitizeAmount(amount string) float64 {
	reg, err := regexp.Compile("[^0-9-,]+")
	if err != nil {
		panic(err)
	}

	cleanedString := reg.ReplaceAllString(amount, "")
	cleanedString = strings.ReplaceAll(cleanedString, ",", ".")

	if cleanedString == "" {
		return 0
	}

	amountFloat, err := strconv.ParseFloat(cleanedString, 64)
	if err != nil {
		log.Printf("Failed float conversion: %v", err)
	}

	return amountFloat
}
