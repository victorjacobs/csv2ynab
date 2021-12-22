package model

import "time"

type Transaction struct {
	Date   time.Time
	Payee  string
	Memo   string
	Amount float64
}
