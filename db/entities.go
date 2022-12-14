package db

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"gorm.io/gorm"
)

type AlpacaOrder struct {
	Symbol     string
	Side       alpaca.Side
	StopLost   float64
	TakeProfit float64
	ExpiredAt  time.Time
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
