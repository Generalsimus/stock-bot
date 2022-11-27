package algo

import (
	"fmt"
	"neural/finance"
	"neural/options"
	"sort"
)

type SymbolBestTimeIntervalsBars struct {
	Symbol   string
	Interval TimeIntervalsBars
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
		for _, interval := range intervals {
			symbolBestTimeIntervals = append(symbolBestTimeIntervals, &SymbolBestTimeIntervalsBars{
				Symbol:   symbol,
				Interval: interval,
			})
		}
	}

	sortedBestSymbols := SortBestSymbolInterval(symbolBestTimeIntervals)
	fmt.Println("BEST_INTERVAL_SYMBOL: ", sortedBestSymbols[0].Symbol)
	// sss := sortedBestSymbols[0].bestIntervals[0]
	return sortedBestSymbols
}
func GetSymbolIntervalSimilarityNum(interval TimeIntervalsBars) float64 {
	return (SumDistanceSimilarityNum(interval.intervalSimilarity) / float64(len(interval.intervalSimilarity)))
}

func SortBestSymbolInterval(symbolBestTimeIntervals []*SymbolBestTimeIntervalsBars) []*SymbolBestTimeIntervalsBars {

	sort.Slice(symbolBestTimeIntervals, func(index1, index2 int) bool {
		symbolInterval1 := symbolBestTimeIntervals[index1]
		symbolInterval2 := symbolBestTimeIntervals[index2]

		return GetSymbolIntervalSimilarityNum(symbolInterval1.Interval) < GetSymbolIntervalSimilarityNum(symbolInterval2.Interval)
	})

	return symbolBestTimeIntervals
}
