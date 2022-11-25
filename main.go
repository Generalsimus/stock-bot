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
	fmt.Println("outputTT")

	// market.InitAndGetAlpacaClient()
	// marketObj := market.NewMarket()

	SymbolsSimilarity := algo.GetSymbolsSimilarity()
	// openPos :=
	algo.FindAlpacaRelativeOpenPosition(SymbolsSimilarity[0].BestIntervals[0])

	/////////////////////////////////////////
	// fmt.Printf("%+v\n", openPos)
	// symbol := SymbolsSimilarity[0].Symbol
	// marketObj.OrderMarket(symbol, openPos.StopLost, openPos.TakeProfit)
	//////////////////////////////////
	drawValue := algo.ConvertToDrawWindow(SymbolsSimilarity[0].BestIntervals[0], options.ViewCandles)
	draw.DrawOnNewWindow(drawValue, options.ViewCandles)
}
