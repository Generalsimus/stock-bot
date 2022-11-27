package main

import (
	"fmt"
	"neural/algo"
	"neural/draw"
	"neural/market"
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
	marketObj := market.NewMarket()

	SymbolsSimilarity := algo.GetSymbolsSimilarity()
	// openPos :=
	position := algo.FindAlpacaRelativeOpenPosition(SymbolsSimilarity[0].Interval)

	/////////////////////////////////////////
	symbol := SymbolsSimilarity[0].Symbol
	marketObj.CheckOrder(symbol)
	// fmt.Printf("%+v\n", openPos)
	// symbol := SymbolsSimilarity[0].Symbol
	marketObj.OrderMarket(symbol, position)
	//////////////////////////////////
	drawValue := algo.ConvertToDrawWindow(SymbolsSimilarity[0].Interval, options.ViewCandles)
	draw.DrawOnNewWindow(drawValue, options.ViewCandles)
}
