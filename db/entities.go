package db

import (
	financeGo "github.com/piquette/finance-go"
	"gorm.io/gorm"
)

type Order struct {
	order    string
	interval uint
	symbol   string
	gorm.Model
}
type Bar struct {
	Symbol string
	financeGo.ChartBar
	gorm.Model
}
