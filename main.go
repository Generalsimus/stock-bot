package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println("outputTT")

	// iter := finance.GetStockData("TD", datetime.OneHour, datetime.Datetime{Month: 1, Day: 1, Year: 2022})
	// iter.

	// output := []*financeGo.ChartBar{}
	// output := []float64{8, 7, 1, 2, 5}
	// for iter.Next() {
	// 	Bar := iter.Bar()

	// 	output = append(output, Bar)

	// }
	const (
		// რამდენი საუკეთესო დამთხვევა აირჩეს
		bestCandles int = 4
		// რამდენი ბარი შეამოწმოს მინიმალური
		startIntervalCount int = 4
		// რამდენი ბარი შეამოწმოს მაქსიმალური
		endIntervalCount int = 7
		// ვიზუალიზაცისას დამატებით რამდენი ბარი აჩვენოს ფანჯარაში
		viewCandles int = 5
	)
	// calcAlgo.CalcManyIntervals(bestCandles, startIntervalCount, endIntervalCount, viewCandles, output)

}
