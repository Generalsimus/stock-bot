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
	Symbol          string
	LastBar         db.Bar
	HourFrame       float64
	TimeFrameInHour float64
	Interval        TimeIntervalsBars
}

func GetSymbolsSimilarity() []*SymbolBestTimeIntervalsBars {
	marketData := market.NewMarketData()
	var symbolBestTimeIntervals []*SymbolBestTimeIntervalsBars
	/////////////////////
	endTime := time.Now()
	// dayInTs := int64(60 * 60 * 50)
	// startTime := time.Unix(endTime.Unix()-(dayInTs*60), 0).Round(time.Minute)
	startTime := options.MaxGetBarsStartTime
	fmt.Println("GET BARS: \n", startTime, "\n", endTime)

	/////////
	for _, hourFrame := range options.CheckFrameHours {

		for _, symbol := range options.CheckSymbols {
			symbolFrameBestTimeIntervals := []*SymbolBestTimeIntervalsBars{}
			symbolBars := marketData.GetMarketCachedData(symbol, startTime, endTime)
			frameBars := marketData.CutBarsWithHourFrame(symbolBars, hourFrame)
			lastBar := symbolBars[len(symbolBars)-1]
			////////////////////////////////////
			// fmt.Println("BARS LEN: ", len(frameBars))
			intervals := CalcManyIntervals(
				options.BestCandles,
				options.StartIntervalCount,
				options.EndIntervalCount,
				options.ViewCandles,
				frameBars,
			)
			for _, interval := range intervals {
				// interval.interval
				symbolFrameBestTimeIntervals = append(symbolFrameBestTimeIntervals, &SymbolBestTimeIntervalsBars{
					Symbol:          symbol,
					LastBar:         lastBar,
					HourFrame:       hourFrame,
					TimeFrameInHour: float64(interval.interval) * hourFrame,
					Interval:        interval,
				})
			}
			sortedBestSymbols := SortBestSymbolInterval(symbolFrameBestTimeIntervals)
			fmt.Println("sortedBestSymbols len", len(sortedBestSymbols))
			if len(sortedBestSymbols) != 0 {
				symbolBestTimeIntervals = append(symbolBestTimeIntervals, sortedBestSymbols[0])
			}
		}

	}

	// fmt.Println("BEST_INTERVAL_SYMBOL: ", sortedBestSymbols[0].Symbol)
	// sss := sortedBestSymbols[0].bestIntervals[0]
	fmt.Println("symbolBestTimeIntervals len", len(symbolBestTimeIntervals))
	return symbolBestTimeIntervals
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
