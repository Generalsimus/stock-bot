package market

import (
	"context"
	"fmt"
	"log"
	"neural/db"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func NewMarket() Market {
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
	account, err := client.GetAccount()

	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("%+v\n", *account)
	return Market{
		client:     client,
		options:    options,
		account:    *account,
		marketData: NewMarketData(),
	}
}
func NewMarketData() MarketData {
	dbConnect := db.GetDb()
	options := marketdata.ClientOpts{
		// Alternatively you can set your key and secret using the
		// APCA_API_KEY_ID and APCA_API_SECRET_KEY environment variables
		ApiKey:    os.Getenv("AlpacaApiKey"),
		ApiSecret: os.Getenv("AlpacaApiSecret"),
	}
	client := marketdata.NewClient(options)

	return MarketData{
		client:  client,
		options: options,
		db:      dbConnect,
	}
}
