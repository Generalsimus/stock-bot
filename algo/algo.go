package algo

import (
	"fmt"
	"math"
	"neural/utils"
	"sort"
	"time"

	financeGo "github.com/piquette/finance-go"
)

type BarsInterval struct {
	startIndex    int
	endIndex      int
	similarityNum float64
	bars          []*financeGo.ChartBar
	barsList      []*financeGo.ChartBar
}
type TimeIntervalsBars struct {
	interval           int
	intervalSimilarity []*BarsInterval
}

func CalcManyIntervals(bestCandles int, startIntervalCount int, endIntervalCount int, viewCandles int, barsList []*financeGo.ChartBar) []TimeIntervalsBars {
	barsList = SortBarsWithTimestamp(barsList)
	timeIntervalTopBarsList := []TimeIntervalsBars{}
	for interval := startIntervalCount; interval <= endIntervalCount; interval++ {
		slicedBars := SliceBars(barsList, interval)
		intervalSimilarity := intervalSimilarity(interval, slicedBars)

		timeIntervalTopBarsList = append(timeIntervalTopBarsList, TimeIntervalsBars{
			interval:           interval,
			intervalSimilarity: intervalSimilarity[0:int(math.Min(float64(bestCandles), float64(len(intervalSimilarity))))],
		})

	}

	return SortWithDistanceSum(timeIntervalTopBarsList)
}

func SortWithDistanceSum(timeBarsIntervals []TimeIntervalsBars) []TimeIntervalsBars {
	sort.Slice(timeBarsIntervals, func(index1, index2 int) bool {
		interval1 := timeBarsIntervals[index1]
		interval2 := timeBarsIntervals[index2]
		// fmt.Println("VAL 1:%v VAL 2: %v", SumDistanceSimilarityNum(interval1.intervalSimilarity), SumDistanceSimilarityNum(interval2.intervalSimilarity))
		return SumDistanceSimilarityNum(interval1.intervalSimilarity) < SumDistanceSimilarityNum(interval2.intervalSimilarity)
	})
	return timeBarsIntervals
}
func getBarValue(bar *financeGo.ChartBar) float64 {
	BarOpenPrice, _ := bar.Open.Float64()
	barClosePrice, _ := bar.Close.Float64()
	barHighPrice, _ := bar.High.Float64()
	barLowPrice, _ := bar.Low.Float64()
	return (BarOpenPrice + barClosePrice + barHighPrice + barLowPrice) / 4
}
func MapToBarsValues(bars []*financeGo.ChartBar) []float64 {
	values := []float64{}
	for _, bar := range bars {
		values = append(values, getBarValue(bar))
	}
	return values
}
func BarsToOneNumber(bars []*financeGo.ChartBar) []float64 {
	barsToValue := MapToBarsValues(bars)
	values := []float64{}
	_, max := utils.FindMinAndMax(barsToValue)
	for _, bar := range bars {
		barValue := getBarValue(bar)
		values = append(values, barValue/max)

	}
	return values
}
func CalculateSimilarityBars(primelyBars []*financeGo.ChartBar, secondaryBars []*financeGo.ChartBar) float64 {
	var sumInterval float64 = 0
	primelyOneNumBars := BarsToOneNumber(primelyBars)
	secondaryOneNumBars := BarsToOneNumber(secondaryBars)
	for index, primelyNum := range primelyOneNumBars {
		secondaryNum := secondaryOneNumBars[index]
		min, max := utils.FindMinAndMax([]float64{primelyNum, secondaryNum})
		sumInterval = sumInterval + (max - min)
	}
	return sumInterval
}
func intervalSimilarity(interval int, slicedBars []*BarsInterval) []*BarsInterval {
	checkInterval := slicedBars[0]
	sort.Slice(slicedBars, func(index1, index2 int) bool {
		interval1 := slicedBars[index1]
		interval2 := slicedBars[index2]
		interval1.similarityNum = CalculateSimilarityBars(checkInterval.bars, interval1.bars) / float64(len(interval1.bars))
		interval2.similarityNum = CalculateSimilarityBars(checkInterval.bars, interval2.bars) / float64(len(interval1.bars))
		return interval1.similarityNum < interval2.similarityNum
	})
	return slicedBars
}
func CutBestInterval(bestCandles int, intervalBars []BarsInterval) []BarsInterval {
	closestBarsIntervals := []BarsInterval{}

labelFor:
	for _, algoBar := range intervalBars {

		if len(closestBarsIntervals) == bestCandles {
			break
		}
		for _, addedBars := range closestBarsIntervals {

			if isBarInIntervalBarsTime(addedBars.bars, algoBar.bars) {
				continue labelFor
				break
			}
		}

		closestBarsIntervals = append(closestBarsIntervals, algoBar)
	}
	return closestBarsIntervals
}
func SliceBars(bars []*financeGo.ChartBar, interval int) []*BarsInterval {
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
func ConvertToDrawWindow(timeIntervalBar TimeIntervalsBars, viewCandles int) [][]float64 {
	drawNumberRow := [][]float64{}
	for _, intervalTopBars := range timeIntervalBar.intervalSimilarity {
		drawNumberCol := []float64{}
		bars := intervalTopBars.bars
		fmt.Println("TIME_EEEEE:1 ", time.Unix(int64(bars[0].Timestamp), 0))
		// start := intervalTopBars.startIndex
		// end := int(math.Min(float64(len(intervalTopBars.barsList)-1), float64(intervalTopBars.endIndex+viewCandles)))
		// bars = intervalTopBars.barsList[start:end]
		// bars = intervalTopBars.barsList[intervalTopBars.startIndex:intervalTopBars.endIndex]
		startIndex := int(math.Max(0, float64(intervalTopBars.startIndex-viewCandles)))
		endIndex := intervalTopBars.endIndex
		bars = intervalTopBars.barsList[startIndex:endIndex]
		fmt.Println("TIME_EEEEE:2 ", time.Unix(int64(bars[0].Timestamp), 0), len(bars))
		for _, bar := range bars {

			drawNumberCol = append(drawNumberCol, getBarValue(bar))
		}

		drawNumberRow = append(drawNumberRow, utils.Reverse(drawNumberCol))
		// break
	}

	return drawNumberRow
}

// func Find
// func calcDiffractionAlgo(bars []*financeGo.ChartBar) float64 {
// 	var openClosePercentSum float64 = 0
// 	var maxMinPercentSum float64 = 0
// 	firstBar := bars[0]
// 	LastBar := bars[len(bars)-1]
// 	firstBarOpenPrice, _ := firstBar.Open.Float64()
// 	firstBarClosePrice, _ := LastBar.Open.Float64()
// 	barsCount := float64(len(bars))
// 	// var firstElement *financeGo.ChartBar = bars[0]
// 	// var lastElement *financeGo.ChartBar = bars[len(bars)]
// 	// min, max := utils.FindMinAndMax(output)

// 	for barIndex, bar := range bars {
// 		locationNum := float64(barIndex+1) / float64(barsCount)
// 		open, _ := bar.Open.Float64()
// 		close, _ := bar.Close.Float64()
// 		openClosePercentSum += open / close * locationNum
// 		high, _ := bar.High.Float64()
// 		low, _ := bar.Low.Float64()
// 		maxMinPercentSum += low / high * locationNum
// 	}
// 	value := (openClosePercentSum + maxMinPercentSum) / (barsCount * 2)
// 	if (firstBarOpenPrice - firstBarClosePrice) < 0 {
// 		value = value * -1
// 	}
// 	return value
// }

func getBarTampsList(bars []*financeGo.ChartBar) []int {
	timestampsList := []int{}
	for _, bar := range bars {
		timestampsList = append(timestampsList, bar.Timestamp)
	}
	return timestampsList
}
func isBarInIntervalBarsTime(checkIntervalCandles []*financeGo.ChartBar, insideIntervalCandles []*financeGo.ChartBar) bool {

	timestampsList := getBarTampsList(insideIntervalCandles)
	startTimestamp, endTimestamp := utils.FindMinAndMax(timestampsList)
	startTime := time.Unix(int64(startTimestamp), 0)
	endTime := time.Unix(int64(endTimestamp), 0)

	//////////////////
	checkTampsList := getBarTampsList(checkIntervalCandles)
	checkStartTimestamp, checkEndTimestamp := utils.FindMinAndMax(checkTampsList)

	checkTimeStart := time.Unix(int64(checkStartTimestamp), 0)
	checkTimeEnd := time.Unix(int64(checkEndTimestamp), 0)
	// fmt.Println(checkTimeStart, startTime, endTime)

	return checkTimeStart.Equal(startTime) || (checkTimeStart.After(startTime) && checkTimeStart.Before(endTime)) || (checkTimeEnd.After(startTime) && checkTimeEnd.Before(endTime))
}

func SortBarsWithTimestamp(bars []*financeGo.ChartBar) []*financeGo.ChartBar {
	sort.Slice(bars, func(index1, index2 int) bool {
		element1 := bars[index1]
		element2 := bars[index2]

		return element1.Timestamp > element2.Timestamp
	})
	return bars
}

// type BarsInterval struct {
// 	startIndex int
// 	endIndex   int
// 	algoNumber float64
// 	bars       []*financeGo.ChartBar
// 	allBars    []*financeGo.ChartBar
// }

// func FindClosestIntervals(bestCandles int, intervalBars []BarsInterval) []BarsInterval {
// 	firstBarInterval := intervalBars[0]
// 	firstBarAlgoNumber := firstBarInterval.algoNumber
// 	sort.Slice(intervalBars, func(index1, index2 int) bool {
// 		interval1 := intervalBars[index1]
// 		interval2 := intervalBars[index2]

// 		return math.Abs(float64(firstBarAlgoNumber-interval1.algoNumber)) < math.Abs(float64(firstBarAlgoNumber-interval2.algoNumber))
// 	})
// 	closestBarsIntervals := []BarsInterval{firstBarInterval}
// labelFor:
// 	for _, algoBar := range intervalBars {

// 		if len(closestBarsIntervals) == bestCandles {
// 			break
// 		}
// 		for _, addedBars := range closestBarsIntervals {

// 			if isBarInIntervalBarsTime(addedBars.bars, algoBar.bars) {
// 				continue labelFor
// 				break
// 			}
// 		}

// 		closestBarsIntervals = append(closestBarsIntervals, algoBar)
// 	}

// 	return closestBarsIntervals
// }

// type TimeIntervalsBars struct {
// 	interval        int
// 	topBarsInterval []BarsInterval
// }

// func CalcManyIntervals(bestCandles int, startIntervalCount int, endIntervalCount int, viewCandles int, barsList []*financeGo.ChartBar) []TimeIntervalsBars {
// 	barsList = SortBarsWithTimestamp(barsList)
// 	timeIntervalTopBarsList := []TimeIntervalsBars{}
// 	fmt.Println("TIME_AAAAA: ", time.Unix(int64(barsList[0].Timestamp), 0))
// 	for interval := startIntervalCount; interval <= endIntervalCount; interval++ {
// 		slicedBars := SliceBars(barsList, interval)
// 		bestIntervalCandles := FindClosestIntervals(bestCandles, slicedBars)
// 		// fmt.Println("INTERVAL: ", interval, startIntervalCount, endIntervalCount)
// 		// fmt.Println("LEN: ", len(slicedBars), len(barsList), len(bestIntervalCandles))
// 		timeIntervalTopBarsList = append(timeIntervalTopBarsList, TimeIntervalsBars{
// 			interval:        interval,
// 			topBarsInterval: bestIntervalCandles,
// 		})
// 		// fmt.Println("TIME_EEEEE: ", time.Unix(int64(slicedBars[0].bars[0].Timestamp), 0))
// 		// calcDiffractionAlgo(interval , bars []*financeGo.ChartBar)
// 	}
// 	// bestSortedIntervals := calcAlgo.FindClosestBars(candlesCount, bestCandles, output)
// 	return timeIntervalTopBarsList
// }

// func sumAlgoNumberDiff(barsInterval []BarsInterval) float64 {
// 	var value float64 = 0
// 	for index, item := range barsInterval {
// 		if index == (len(barsInterval) - 1) {
// 			break
// 		}
// 		nextItem := barsInterval[index+1]
// 		// fmt.Println([]float64{item.algoNumber, nextItem.algoNumber})
// 		min, max := utils.FindMinAndMax([]float64{item.algoNumber, nextItem.algoNumber})

// 		value = value + (max - min)
// 	}
// 	return value
// }

// func FindBestTimeIntervalsBars(timeIntervalsBars []TimeIntervalsBars) []TimeIntervalsBars {
// 	sort.Slice(timeIntervalsBars, func(index1, index2 int) bool {
// 		element1 := timeIntervalsBars[index1]
// 		element2 := timeIntervalsBars[index2]

// 		return sumAlgoNumberDiff(element1.topBarsInterval) > sumAlgoNumberDiff(element2.topBarsInterval)
// 	})
// 	return timeIntervalsBars

// }
