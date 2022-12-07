package db

import (
	"gorm.io/gorm"
)

type Order struct {
	order    string
	interval uint
	symbol   string
	gorm.Model
}
type Bar struct {
	Symbol          string
	Frame           string
	Timestamp       int64
	Open            float64
	Close           float64
	High            float64
	Low             float64
	BarStructToJson string
	gorm.Model
}
