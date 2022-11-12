package calcAlgo

import (
	"fmt"
	"math"
	"neural/draw"
	"neural/utils"
	"sort"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	financeGo "github.com/piquette/finance-go"
)

func calcDiffractionAlgo(bars []*financeGo.ChartBar) float64 {
	var openClosePercentSum float64 = 0
	var maxMinPercentSum float64 = 0
	// var firstElement *financeGo.ChartBar = bars[0]
	// var lastElement *financeGo.ChartBar = bars[len(bars)]
	// min, max := utils.FindMinAndMax(output)
	for _, bar := range bars {
		open, _ := bar.Open.Float64()
		close, _ := bar.Close.Float64()
		openClosePercentSum += open / close * 100
		high, _ := bar.High.Float64()
		low, _ := bar.Low.Float64()
		maxMinPercentSum += low / high * 100
	}
	return openClosePercentSum / maxMinPercentSum
}

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
	return (checkTimeStart.After(startTime) && checkTimeStart.Before(endTime)) || (checkTimeEnd.After(startTime) && checkTimeEnd.Before(endTime))
}

func DrawMeOnWindow(BarsList []BarsInterval) {
	// red := color.NRGBA{R: 0x80, A: 0xFF}
	// barViewCount := float64(len(bestSortedIntervals))
	draw.AddListener(func(e system.FrameEvent, ops *op.Ops, w *app.Window) {
		// barsListCount := len(BarsList)
		// var path clip.Path
		// path.Begin(ops)
		// for rowIndex, intervalsBar := range BarsList {

		// 	Bars := PlusBarIntervals(5, intervalsBar).interval

		// 	viewWidth := float64(draw.Width)
		// 	viewHeight := float64(draw.Height) / float64(barsListCount)
		// 	barPriceValues := getBarsValues(Bars)
		// 	min, max := utils.FindMinAndMax(barPriceValues)
		// 	fmt.Println("TIME_START: ", time.Unix(int64(Bars[0].Timestamp), 0), "TIME_END: ", time.Unix(int64(Bars[len(Bars)-1].Timestamp), 0))
		// 	fmt.Println("START: ", barPriceValues[0], "END: ", barPriceValues[len(barPriceValues)-1])
		// 	barsCount := len(barPriceValues)
		// 	for colIndex, priceValue := range barPriceValues {
		// 		y := float32(((viewHeight * float64(rowIndex+1)) - (((priceValue - min) / (max - min)) * viewHeight)))
		// 		x := float32((float64(colIndex) / float64(barsCount-1)) * float64(viewWidth))

		// 		if colIndex != 0 {
		// 			path.LineTo(f32.Pt(x, y))

		// 		}
		// 		path.MoveTo(f32.Pt(x, y))

		// 	}

		// }
		// paint.FillShape(ops, red,
		// 	// Stroke
		// 	// Outline
		// 	clip.Stroke{
		// 		Path:  path.End(),
		// 		Width: 2,
		// 	}.Op())

	})
	// defer draw.Init()
}
func SortBarsWithTimestamp(bars []*financeGo.ChartBar) []*financeGo.ChartBar {
	sort.Slice(bars, func(index1, index2 int) bool {
		element1 := bars[index1]
		element2 := bars[index2]

		return element1.Timestamp > element2.Timestamp
	})
	return bars
}

type BarsInterval struct {
	startIndex int
	endIndex   int
	algoNumber float64
	bars       []*financeGo.ChartBar
	allBars    []*financeGo.ChartBar
}

func SliceBars(bars []*financeGo.ChartBar, interval int) []BarsInterval {
	barsList := []BarsInterval{}
	barsCount := len(bars)
	for i := 0; i < barsCount; i = i + interval {
		startIndex := i
		endIndex := i + interval
		if endIndex > barsCount {
			break
		}
		items := bars[startIndex:endIndex]

		algoNumber := calcDiffractionAlgo(items)
		barsList = append(barsList, BarsInterval{
			startIndex: startIndex,
			endIndex:   endIndex,
			algoNumber: algoNumber,
			bars:       items,
			allBars:    bars,
		})
	}
	return barsList
}
func FindClosestIntervals(bestCandles int, intervalBars []BarsInterval) []BarsInterval {
	firstBarInterval := intervalBars[0]
	firstBarAlgoNumber := firstBarInterval.algoNumber
	sort.Slice(intervalBars, func(index1, index2 int) bool {
		interval1 := intervalBars[index1]
		interval2 := intervalBars[index2]

		return math.Abs(float64(firstBarAlgoNumber-interval1.algoNumber)) < math.Abs(float64(firstBarAlgoNumber-interval2.algoNumber))
	})
	closestBarsIntervals := []BarsInterval{firstBarInterval}
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

type TimeIntervalsBars struct {
	interval        int
	topBarsInterval []BarsInterval
}

func CalcManyIntervals(bestCandles int, startIntervalCount int, endIntervalCount int, viewCandles int, barsList []*financeGo.ChartBar) []TimeIntervalsBars {
	barsList = SortBarsWithTimestamp(barsList)
	timeIntervalTopBarsList := []TimeIntervalsBars{}
	fmt.Println("TIME_AAAAA: ", time.Unix(int64(barsList[0].Timestamp), 0))
	for interval := startIntervalCount; interval <= endIntervalCount; interval++ {
		slicedBars := SliceBars(barsList, interval)
		bestIntervalCandles := FindClosestIntervals(bestCandles, slicedBars)
		fmt.Println("INTERVAL: ", interval, startIntervalCount, endIntervalCount)
		fmt.Println("LEN: ", len(slicedBars), len(barsList), len(bestIntervalCandles))
		timeIntervalTopBarsList = append(timeIntervalTopBarsList, TimeIntervalsBars{
			interval:        interval,
			topBarsInterval: bestIntervalCandles,
		})
		// fmt.Println("TIME_EEEEE: ", time.Unix(int64(slicedBars[0].bars[0].Timestamp), 0))
		// calcDiffractionAlgo(interval , bars []*financeGo.ChartBar)
	}
	// bestSortedIntervals := calcAlgo.FindClosestBars(candlesCount, bestCandles, output)
	return timeIntervalTopBarsList
}
