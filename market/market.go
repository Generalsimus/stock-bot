package market

import (
	"math"
	"neural/db"
	"neural/utils"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

//	type AlpacaPosition struct {
//		Side       alpaca.Side
//		StopLost   float64
//		TakeProfit float64
//	}
type Market struct {
	client     alpaca.Client
	options    alpaca.ClientOpts
	account    alpaca.Account
	marketData MarketData
	db         *gorm.DB
}

func (m Market) OrderMarket(position db.AlpacaOrder) (*alpaca.Order, error) {
	// m.marketData.client.
	utils.LogStruct("NEW ORDER: ", position)

	//////////////////////////////////
	// qty := decimal.NewFromInt(1)
	qty := decimal.NewFromFloat(0.5)
	takeProfitDecimal := decimal.NewFromFloat(math.Floor((position.TakeProfit)*100) / 100)
	stopLossDecimal := decimal.NewFromFloat(math.Floor((position.StopLost)*100) / 100)
	side := position.Side
	//////////////////////////////////
	TakeProfit := alpaca.TakeProfit{LimitPrice: &takeProfitDecimal}
	StopLoss := alpaca.StopLoss{
		LimitPrice: nil,
		StopPrice:  &stopLossDecimal,
	}
	orderOptions := alpaca.PlaceOrderRequest{
		AccountID:   m.account.ID,
		AssetKey:    &position.Symbol,
		Qty:         &qty,
		Side:        side,
		Type:        alpaca.Market,
		TimeInForce: alpaca.Day,
		OrderClass:  alpaca.Bracket,
		TakeProfit:  &TakeProfit,
		StopLoss:    &StopLoss,
		// StopPrice:   &stopLossDecimal,
		// LimitPrice:  &takeProfitDecimal,
	}
	utils.LogStruct("Order Options: ", orderOptions)
	return m.client.PlaceOrder(orderOptions)

}
func (m Market) SaveOnDb(alpacaOrderDetails db.AlpacaOrder) {
	m.db.Create(&alpacaOrderDetails)
}
func (m Market) CheckOrderIsExpired(orderPosition db.AlpacaOrder) bool {
	var dbOrder db.AlpacaOrder
	res := m.db.Where("symbol = ?", orderPosition.Symbol).Where("side = ?", orderPosition.Side).Where("hour_frame = ?", orderPosition.HourFrame).Where("expired_at <= ?", time.Now()).Find(&dbOrder)

	return res.RowsAffected == 0
}
