package main

import (
	"fmt"
	"neural/db"
	"neural/market"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println("outputTT")
	db.Init()
	// timeNow := time.Now()
	// market.GetMarketDataDb("TD", 2, time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()-3, 1, 1, 1, 1, time.UTC))
	market.GetBars("TD", []float64{0.5}, time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC))

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
