package market

import (
	"log"
	"math"
	"neural/utils"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

type AlpacaPosition struct {
	Side       alpaca.Side
	StopLost   float64
	TakeProfit float64
}
type Market struct {
	client     alpaca.Client
	options    alpaca.ClientOpts
	account    alpaca.Account
	marketData MarketData
}

func (m Market) OrderMarket(symbol string, position AlpacaPosition) {
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

	if order, err := m.client.PlaceOrder(alpaca.PlaceOrderRequest{
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
	}); err != nil {
		log.Fatalf("failed place order: %v", err)
		log.Fatalf("failed place order: %v", err)

		utils.LogStruct("ORDER SUCCESSFUL: ", order)
	}
	log.Println("order sent")
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
