package algo

import (
	"fmt"
	"neural/finance"
	"neural/options"
	"sort"
)

type SymbolBestTimeIntervalsBars struct {
	Symbol        string
	BestIntervals []TimeIntervalsBars
}

func GetSymbolsSimilarity() []*SymbolBestTimeIntervalsBars {
	symbolBestTimeIntervals := []*SymbolBestTimeIntervalsBars{}
	for _, symbol := range options.CheckSymbols {
		output := finance.GetSymbolIntervalBars(symbol, options.FinanceInterval, options.FinanceStartDate)
		intervals := CalcManyIntervals(
			options.BestCandles,
			options.StartIntervalCount,
			options.EndIntervalCount,
			options.ViewCandles,
			output,
		)
		symbolBestTimeIntervals = append(symbolBestTimeIntervals, &SymbolBestTimeIntervalsBars{
			Symbol:        symbol,
			BestIntervals: intervals,
		})
	}

	sortedBestSymbols := FindBestSymbolInterval(symbolBestTimeIntervals)
	fmt.Println("BEST_INTERVAL_SYMBOL: ", sortedBestSymbols[0].Symbol)
	// sss := sortedBestSymbols[0].bestIntervals[0]
	return sortedBestSymbols
}
func SumBestSymbolIntervalSimilarityNum(bestIntervals []TimeIntervalsBars) float64 {
	var value float64 = 0
	for _, interval := range bestIntervals {

		value = value + (SumDistanceSimilarityNum(interval.intervalSimilarity) / float64(len(interval.intervalSimilarity)))
	}
	return value / float64(len(bestIntervals))
}
func FindBestSymbolInterval(symbolBestTimeIntervals []*SymbolBestTimeIntervalsBars) []*SymbolBestTimeIntervalsBars {

	sort.Slice(symbolBestTimeIntervals, func(index1, index2 int) bool {
		symbolInterval1 := symbolBestTimeIntervals[index1]
		symbolInterval2 := symbolBestTimeIntervals[index2]

		return SumBestSymbolIntervalSimilarityNum(symbolInterval1.BestIntervals) < SumBestSymbolIntervalSimilarityNum(symbolInterval2.BestIntervals)
	})

	return symbolBestTimeIntervals
}
