package market

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

func GetMarketData(symbol string, timeFrame marketdata.TimeFrame, startDate time.Time) []marketdata.Bar {
	client := GetMarketDataClient()
	// sss, ret := client.GetLatestBars([]string{"TWTR"})
	// fmt.Println(sss)
	// fmt.Println("LLL", ret)
	// // marketdata.GetLatestBars
	// // marketdata.StreamTradeUpdatesInBackground(context.TODO(), func(tu alpaca.TradeUpdate) {
	// // 	log.Printf("TRADE UPDATE: %+v\n", tu)
	// // })
	bars, err := client.GetBars(symbol, marketdata.GetBarsParams{
		TimeFrame: timeFrame,
		Start:     startDate,
		End:       time.Now(),
	})
	if err != nil {
		fmt.Printf("GET '%v' BARS:\n", err)
		panic(err)
	}
	fmt.Printf("GET '%v' BARS:\n", symbol)
	for _, bar := range bars {
		// bar.High
		fmt.Printf("%+v\n", bar)
	}
	return bars
}
