package market

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

type Market struct {
	client  alpaca.Client
	options alpaca.ClientOpts
}

func (m Market) OrderMarket(symbol string, stopLoss float64, takeProfit float64) {
	qty := decimal.NewFromInt(1)
	//////////////////////////////////
	takeProfitDecimal := decimal.NewFromFloat(math.Floor(takeProfit*100) / 100)
	stopLossDecimal := decimal.NewFromFloat(math.Floor(stopLoss*100) / 100)
	TakeProfit := alpaca.TakeProfit{LimitPrice: &takeProfitDecimal}
	StopLoss := alpaca.StopLoss{
		LimitPrice: nil,
		StopPrice:  &stopLossDecimal,
	}
	if _, err := m.client.PlaceOrder(alpaca.PlaceOrderRequest{
		AssetKey:    &symbol,
		Qty:         &qty,
		Side:        alpaca.Buy,
		Type:        alpaca.Limit,
		TimeInForce: alpaca.GTC,
		OrderClass:  alpaca.Bracket,
		TakeProfit:  &TakeProfit,
		StopLoss:    &StopLoss,
		// StopPrice:  &stopLossDecimal,
		// LimitPrice: &takeProfitDecimal,
	}); err != nil {
		log.Fatalf("failed place order: %v", err)
	}
	log.Println("order sent")
}

func New() Market {
	fmt.Println("ALPACA CONNECTING...")
	alpaca.StreamTradeUpdatesInBackground(context.TODO(), func(tu alpaca.TradeUpdate) {
		log.Printf("TRADE UPDATE: %+v\n", tu)
	})
	options := alpaca.ClientOpts{

		ApiKey:    os.Getenv("AlpacaApiKey"),
		ApiSecret: os.Getenv("AlpacaApiSecret"),
		BaseURL:   os.Getenv("AlpacaBaseURL"),
	}
	client := alpaca.NewClient(options)
	acct, err := client.GetAccount()

	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("%+v\n", *acct)
	return Market{
		client:  client,
		options: options,
	}
}
