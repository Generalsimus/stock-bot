package main

import (
	"fmt"
	"neural/algo"
	"neural/draw"
	"neural/options"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// marketData := market.NewMarketData()
	// fmt.Println(marketData)
	// hourFrame float64, symbol string, startTime time.Time, endTime time.Time

	// option.MaxGetBarsStartTime
	SymbolsSimilarity := algo.GetSymbolsSimilarity()
	fmt.Println("BEST Len", len(SymbolsSimilarity))
	fmt.Println("BEST Symbol", SymbolsSimilarity[0].Symbol)
	drawValue := algo.ConvertToDrawWindow(SymbolsSimilarity[0].Interval, options.ViewCandles)
	draw.DrawOnNewWindow(drawValue, options.ViewCandles)
	// for index, _ := range bars {
	// 	if index == 0 {
	// 		continue
	// 	}
	// 	bar1 := bars[index-1]
	// 	bar2 := bars[index]
	// 	fmt.Println("CUT BARS DIFF: \n", bar2.Timestamp-bar1.Timestamp)
	// }
	// utils.LogStruct("RESULT: ", len(bars))
	// lc, _ := time.LoadLocation("Asia/Tbilisi")

	// utils.LogStruct("LAST BAR TIME: ", time.Unix(bars[len(bars)-1].Timestamp, 0).In(lc))
	// for _, bar := range bars {

	// 	utils.LogStruct("RESULT: ", bar, time.Unix(bar.Timestamp, 0))
	// }
	// utils.LogStruct("RESULT: ", bars)
	// options := marketdata.ClientOpts{
	// 	// Alternatively you can set your key and secret using the
	// 	// APCA_API_KEY_ID and APCA_API_SECRET_KEY environment variables
	// 	ApiKey:    os.Getenv("AlpacaApiKey"),
	// 	ApiSecret: os.Getenv("AlpacaApiSecret"),
	// }
	// client := marketdata.NewClient(options)
	// account, err := client.GetAccount()
	// account.
	// account.GetQuotes()
	////////////////////////////////////////////////////////////////////////

	// quotes, err := client.GetBars("META", marketdata.GetBarsParams{
	// 	TimeFrame: marketdata.OneMin,
	// 	Start:     time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC),
	// 	End:       time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC),
	// 	AsOf:      "2022-06-10", // Leaving it empty yields the same results
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("TSLA quotes:")
	// for _, quote := range quotes {
	// 	fmt.Printf("%+v\n", quote)
	// }
	/////////////////////////////////////////
	// rsi2 := talib.Rsi(spy.Close, 2)
	// fmt.Println(rsi2)
	// fmt.Println("outputTT")
	// db.Init()
	// for _, symbol := range options.CheckSymbols {
	// 	for _, interval := range options.FinanceIntervals {
	// 		bars := finance.GetMaxIntervalBars(symbol, interval)
	// 		fmt.Println("BARS: ", len(bars))
	// 	}
	// }

	// timeNow := time.Now()
	// market.GetMarketDataDb("TD", 2, time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()-3, 1, 1, 1, 1, time.UTC))
	// market.GetBars("TD", []float64{0.5}, time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC))

	/*
		// market.InitAndGetAlpacaClient()
		marketObj := market.NewMarket()

		SymbolsSimilarity := algo.GetSymbolsSimilarity()
		// openPos :=
		position := algo.FindAlpacaRelativeOpenPosition(SymbolsSimilarity[0].Interval)
		/////////////////////////////////////////
		market.GetMarketDataDb("TD", 2, time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC))
		/////////////////////////////////////////
		symbol := SymbolsSimilarity[0].Symbol
		marketObj.CheckOrder(symbol)
		// fmt.Printf("%+v\n", openPos)
		// symbol := SymbolsSimilarity[0].Symbol
		marketObj.OrderMarket(symbol, position)
		//////////////////////////////////
		drawValue := algo.ConvertToDrawWindow(SymbolsSimilarity[0].Interval, options.ViewCandles)
		draw.DrawOnNewWindow(drawValue, options.ViewCandles)
	*/
}
