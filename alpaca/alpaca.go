package alpaca

import (
	"context"
	"fmt"
	"log"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

var AlpacaClient *alpaca.Client

func InitAlpaca() {
	fmt.Println("ALPACA CONNECTING...")
	alpaca.StreamTradeUpdatesInBackground(context.TODO(), func(tu alpaca.TradeUpdate) {
		log.Printf("TRADE UPDATE: %+v\n", tu)
	})
	AlpacaClient := alpaca.NewClient(alpaca.ClientOpts{
		// Alternatively you can set your key and secret using the
		// APCA_API_KEY_ID and APCA_API_SECRET_KEY environment variables
		ApiKey:    os.Getenv("AlpacaApiKey"),
		ApiSecret: os.Getenv("AlpacaApiSecret"),
		BaseURL:   os.Getenv("AlpacaBaseURL"),
	})
	acct, err := AlpacaClient.GetAccount()

	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *acct)
}
