package algo

import (
	"fmt"
	"neural/db"
	"neural/market"
	"neural/options"
	"sort"
	"time"
)

type SymbolBestTimeIntervalsBars struct {
	Symbol    string
	LastBar   db.Bar
	HourFrame float64
	Interval  TimeIntervalsBars
}

func GetSymbolsSimilarity() []*SymbolBestTimeIntervalsBars {
	symbolBestTimeIntervals := []*SymbolBestTimeIntervalsBars{}
	marketData := market.NewMarketData()
	/////////////////////
	endTime := time.Now()
	// dayInTs := int64(60 * 60 * 50)
	// startTime := time.Unix(endTime.Unix()-(dayInTs*60), 0).Round(time.Minute)
	startTime := options.MaxGetBarsStartTime
	fmt.Println("GET BARS: \n", startTime, "\n", endTime)

	/////////
	for _, symbol := range options.CheckSymbols {
		symbolBars := marketData.GetMarketCachedData(symbol, startTime, endTime)
		lastBar := symbolBars[len(symbolBars)-1]
		for _, hourFrame := range options.CheckFrameHours {
			frameBars := marketData.CutBarsWithHourFrame(symbolBars, hourFrame)
			for index := range frameBars {
				if index == 0 {
					continue
				}
				bar1 := frameBars[index-1]
				bar2 := frameBars[index]
				fmt.Println("CUT BARS DIFF: \n", bar2.Timestamp-bar1.Timestamp)
			}
			fmt.Println("BARS LEN: ", len(frameBars))
			intervals := CalcManyIntervals(
				options.BestCandles,
				options.StartIntervalCount,
				options.EndIntervalCount,
				options.ViewCandles,
				frameBars,
			)
			for _, interval := range intervals {
				symbolBestTimeIntervals = append(symbolBestTimeIntervals, &SymbolBestTimeIntervalsBars{
					Symbol:    symbol,
					LastBar:   lastBar,
					HourFrame: hourFrame,
					Interval:  interval,
				})
			}
		}
	}

	sortedBestSymbols := SortBestSymbolInterval(symbolBestTimeIntervals)
	// fmt.Println("BEST_INTERVAL_SYMBOL: ", sortedBestSymbols[0].Symbol)
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
