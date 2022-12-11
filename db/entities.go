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
	Timestamp       int64 `gorm:"unique"`
	Open            float64
	Close           float64
	High            float64
	Low             float64
	BarStructToJson string
	gorm.Model
}
