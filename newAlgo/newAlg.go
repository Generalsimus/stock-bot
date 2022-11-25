package newAlgo

import (
	"fmt"
	"math"
	"neural/finance"
	"neural/options"
	"neural/utils"
	"sort"

	financeGo "github.com/piquette/finance-go"
)

type SymbolsBars struct {
	SymbolBars         map[string][]*financeGo.ChartBar
	StartIntervalCount int
	EndIntervalCount   int
	ViewCandles        int
}

func (d SymbolsBars) FindBestSymBolSimilarity() {

}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) GetSortedSymbolsWithSimilarity() []*SymbolBestTimeIntervalsBars {
	symbolBestTimeIntervals := []*SymbolBestTimeIntervalsBars{}
	for _, symbol := range options.CheckSymbols {
		output := finance.GetSymbolIntervalBars(symbol, options.FinanceInterval, options.FinanceStartDate)
		intervals := d.CalcManyIntervals(
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

	sortedBestSymbols := d.FindBestSymbolInterval(symbolBestTimeIntervals)
	fmt.Println("BEST_INTERVAL_SYMBOL: ", sortedBestSymbols[0].Symbol)
	// sss := sortedBestSymbols[0].bestIntervals[0]
	return sortedBestSymbols
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) FindBestSymbolInterval(symbolBestTimeIntervals []*SymbolBestTimeIntervalsBars) []*SymbolBestTimeIntervalsBars {

	sort.Slice(symbolBestTimeIntervals, func(index1, index2 int) bool {
		symbolInterval1 := symbolBestTimeIntervals[index1]
		symbolInterval2 := symbolBestTimeIntervals[index2]

		return d.SumBestSymbolIntervalSimilarityNum(symbolInterval1.BestIntervals) < d.SumBestSymbolIntervalSimilarityNum(symbolInterval2.BestIntervals)
	})

	return symbolBestTimeIntervals
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) SumDistanceSimilarityNum(barIntervals []*BarsInterval) float64 {
	var sumNum float64 = 0
	for _, interval := range barIntervals {
		sumNum = sumNum + (interval.similarityNum / float64(len(interval.bars)))
	}
	// fmt.Println("SUM: ", sumNum/float64(len(barIntervals)))
	return sumNum / float64(len(barIntervals))
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) SumBestSymbolIntervalSimilarityNum(bestIntervals []TimeIntervalsBars) float64 {
	var value float64 = 0
	for _, interval := range bestIntervals {

		value = value + (d.SumDistanceSimilarityNum(interval.intervalSimilarity) / float64(len(interval.intervalSimilarity)))
	}
	return value / float64(len(bestIntervals))
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) SliceBars(bars []*financeGo.ChartBar, interval int) []*BarsInterval {
	barsList := []*BarsInterval{}
	barsCount := len(bars)
	// checkInterval := slicedBars[0]
	for i := 0; i < barsCount; i = i + interval {
		startIndex := i
		endIndex := i + interval
		if endIndex >= barsCount {
			break
		}
		items := bars[startIndex:endIndex]

		// algoNumber := calcDiffractionAlgo(items)
		barsList = append(barsList, &BarsInterval{
			startIndex:    startIndex,
			endIndex:      endIndex,
			similarityNum: 0,
			bars:          items,
			barsList:      bars,
		})
	}
	return barsList
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) getBarValue(bar *financeGo.ChartBar) float64 {
	BarOpenPrice, _ := bar.Open.Float64()
	barClosePrice, _ := bar.Close.Float64()
	barHighPrice, _ := bar.High.Float64()
	barLowPrice, _ := bar.Low.Float64()
	return (BarOpenPrice + barClosePrice + barHighPrice + barLowPrice) / 4
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) MapToBarsValues(bars []*financeGo.ChartBar) []float64 {
	values := []float64{}
	for _, bar := range bars {
		values = append(values, d.getBarValue(bar))
	}
	return values
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) BarsToOneNumber(bars []*financeGo.ChartBar) []float64 {
	barsToValue := d.MapToBarsValues(bars)
	values := []float64{}
	_, max := utils.FindMinAndMax(barsToValue)
	for _, bar := range bars {
		barValue := d.getBarValue(bar)
		values = append(values, barValue/max)

	}
	return values
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) CalculateSimilarityBars(primelyBars []*financeGo.ChartBar, secondaryBars []*financeGo.ChartBar) float64 {
	var sumInterval float64 = 0
	primelyOneNumBars := d.BarsToOneNumber(primelyBars)
	secondaryOneNumBars := d.BarsToOneNumber(secondaryBars)
	for index, primelyNum := range primelyOneNumBars {
		secondaryNum := secondaryOneNumBars[index]
		min, max := utils.FindMinAndMax([]float64{primelyNum, secondaryNum})
		sumInterval = sumInterval + (max - min)
	}
	return sumInterval
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) intervalSimilarity(interval int, slicedBars []*BarsInterval) []*BarsInterval {
	checkInterval := slicedBars[0]
	sort.Slice(slicedBars, func(index1, index2 int) bool {
		interval1 := slicedBars[index1]
		interval2 := slicedBars[index2]
		interval1.similarityNum = d.CalculateSimilarityBars(checkInterval.bars, interval1.bars) / float64(len(interval1.bars))
		interval2.similarityNum = d.CalculateSimilarityBars(checkInterval.bars, interval2.bars) / float64(len(interval1.bars))
		return interval1.similarityNum < interval2.similarityNum
	})
	return slicedBars
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) CalcManyIntervals(bestCandles int, startIntervalCount int, endIntervalCount int, viewCandles int, barsList []*financeGo.ChartBar) []TimeIntervalsBars {
	barsList = d.SortBarsWithTimestamp(barsList)
	timeIntervalTopBarsList := []TimeIntervalsBars{}
	for interval := startIntervalCount; interval <= endIntervalCount; interval++ {
		slicedBars := d.SliceBars(barsList, interval)
		intervalSimilarity := d.intervalSimilarity(interval, slicedBars)

		timeIntervalTopBarsList = append(timeIntervalTopBarsList, TimeIntervalsBars{
			interval:           interval,
			intervalSimilarity: intervalSimilarity[0:int(math.Min(float64(bestCandles), float64(len(intervalSimilarity))))],
		})

	}

	return d.SortWithDistanceSum(timeIntervalTopBarsList)
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) SortWithDistanceSum(timeBarsIntervals []TimeIntervalsBars) []TimeIntervalsBars {
	sort.Slice(timeBarsIntervals, func(index1, index2 int) bool {
		interval1 := timeBarsIntervals[index1]
		interval2 := timeBarsIntervals[index2]
		// fmt.Println("VAL 1:%v VAL 2: %v", SumDistanceSimilarityNum(interval1.intervalSimilarity), SumDistanceSimilarityNum(interval2.intervalSimilarity))
		return d.SumDistanceSimilarityNum(interval1.intervalSimilarity) < d.SumDistanceSimilarityNum(interval2.intervalSimilarity)
	})
	return timeBarsIntervals
}

// ///////////////////////////////////////////////////////
func (d SymbolsBars) SortBarsWithTimestamp(bars []*financeGo.ChartBar) []*financeGo.ChartBar {
	sort.Slice(bars, func(index1, index2 int) bool {
		element1 := bars[index1]
		element2 := bars[index2]

		return element1.Timestamp > element2.Timestamp
	})
	return bars
}
