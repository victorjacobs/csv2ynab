package util

import (
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func SanitizeAmount(amount string) float64 {
	r, err := regexp.Compile("[^0-9-,.]+")
	if err != nil {
		panic(err)
	}

	cleanedString := r.ReplaceAllString(amount, "")
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

func StringBetween(s, start, end string) (string, bool) {
	split := strings.Split(s, start)
	if len(split) < 2 {
		return "", false
	}
	split = strings.Split(split[1], end)
	if len(split) < 2 {
		return "", false
	}

	return strings.TrimSpace(split[0]), true
}
