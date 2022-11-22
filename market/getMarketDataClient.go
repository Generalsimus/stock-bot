package market

import (
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func GetMarketDataClient() marketdata.Client {

	return marketdata.NewClient(marketdata.ClientOpts{
		// Alternatively you can set your key and secret using the
		// APCA_API_KEY_ID and APCA_API_SECRET_KEY environment variables
		ApiKey:    os.Getenv("AlpacaApiKey"),
		ApiSecret: os.Getenv("AlpacaApiSecret"),
	})
}
