package db

import "gorm.io/gorm"

type Order struct {
	order    string
	interval uint
	symbol   string
	gorm.Model
}
