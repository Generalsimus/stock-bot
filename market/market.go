package market

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

type Market struct {
	client     alpaca.Client
	options    alpaca.ClientOpts
	account    alpaca.Account
	marketData MarketData
}

func (m Market) OrderMarket(symbol string, stopLoss float64, takeProfit float64) {
	// m.marketData.client.
	// (status *string, until *time.Time, limit *int, nested *bool)
	status := "all"
	until := time.Now()
	limit := 1000
	nested := false
	order, _ := m.client.ListOrders(&status, &until, &limit, &nested)
	// m
	fmt.Println("POSSSSS: ", order)
	for _, pos := range order {

		fmt.Println("Position: ", pos.Status, pos.Symbol, pos.ID)
	}

	//////////////////////////////////
	qty := decimal.NewFromInt(1)
	takeProfitDecimal := decimal.NewFromFloat(math.Floor((takeProfit)*100) / 100)
	stopLossDecimal := decimal.NewFromFloat(math.Floor((stopLoss)*100) / 100)
	//////////////////////////////////
	fmt.Println("T", takeProfitDecimal, "S", stopLossDecimal)
	TakeProfit := alpaca.TakeProfit{LimitPrice: &takeProfitDecimal}
	StopLoss := alpaca.StopLoss{
		LimitPrice: nil,
		StopPrice:  &stopLossDecimal,
	}

	if _, err := m.client.PlaceOrder(alpaca.PlaceOrderRequest{
		AccountID:   m.account.ID,
		AssetKey:    &symbol,
		Qty:         &qty,
		Side:        alpaca.Sell,
		Type:        alpaca.Market,
		TimeInForce: alpaca.Day,
		OrderClass:  alpaca.Bracket,
		TakeProfit:  &TakeProfit,
		StopLoss:    &StopLoss,
		// StopPrice:   &stopLossDecimal,
		// LimitPrice:  &takeProfitDecimal,
	}); err != nil {
		log.Fatalf("failed place order: %v", err)
	}
	log.Println("order sent")
}
