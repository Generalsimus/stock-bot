package market

import (
	"fmt"
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

func (m Market) OrderMarket(symbol string, position db.AlpacaOrder) {
	// m.marketData.client.
	utils.LogStruct("NEW ORDER: ", position)

	//////////////////////////////////
	qty := decimal.NewFromInt(1)
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
		AssetKey:    &symbol,
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
	if order, err := m.client.PlaceOrder(orderOptions); err != nil {
		fmt.Printf("failed place order: %v\n", err)

		utils.LogStruct("ORDER SUCCESSFUL: ", order)
	} else {

	}

}
func (m Market) CheckOrder(symbol string) {

	status := "all"
	until := time.Now()
	limit := 1000
	nested := false
	orders, _ := m.client.ListOrders(&status, &until, &limit, &nested)
	////////////////////////////////////////////////////////////////
	for _, order := range orders {
		utils.LogStruct("ORDER: ", order)
	}
	positions, _ := m.client.ListPositions()
	for _, pos := range positions {
		utils.LogStruct("POSITION: ", pos)
	}
}
