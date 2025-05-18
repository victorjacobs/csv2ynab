package model

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

type Transaction struct {
	Date           time.Time
	Payee          string
	Description    string
	Amount         float64
	IdempotencyKey string // In case the fields above don't provide enough uniqueness
}

func (t *Transaction) ImportId() string {
	h := sha1.New()

	hashKey := fmt.Sprintf("%v%v%v%v%v", t.Payee, t.Amount, t.Description, t.Date, t.IdempotencyKey)
	h.Write([]byte(hashKey))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
