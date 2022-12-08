package algo

import (
	"math"
	"neural/market"
	"neural/utils"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func FindAlpacaRelativeOpenPosition(intervalBars TimeIntervalsBars) market.AlpacaPosition {
	intervalSimilarity := intervalBars.intervalSimilarity
	UpAverageMaxDiffSum, DownAverageMaxDiffSum := FindAverageDiffMaxPosition(intervalSimilarity)
	currentPrice := intervalSimilarity[0].bars[0].Close
	side := alpaca.Buy
	var takeProfit float64 = 0
	var stopLost float64 = 0
	if UpAverageMaxDiffSum > DownAverageMaxDiffSum {
		side = alpaca.Buy
		takeProfit = currentPrice + (currentPrice * (UpAverageMaxDiffSum / 1))
		stopLost = currentPrice - (currentPrice * (DownAverageMaxDiffSum / 1))
	} else {
		side = alpaca.Sell
		takeProfit = currentPrice - (currentPrice * (DownAverageMaxDiffSum / 1))
		stopLost = currentPrice + (currentPrice * (UpAverageMaxDiffSum / 1))
	}

	position := market.AlpacaPosition{
		Side:       side,
		TakeProfit: takeProfit,
		StopLost:   stopLost,
	}
	return position
}

type UpDownDiffType struct {
	Scale   float64
	DiffNum float64
	Side    alpaca.Side
}
type MaxUpDownDiffs struct {
	UpMaxDiff  *UpDownDiffType
	DowMaxDiff *UpDownDiffType
}

func FindAverageDiffMaxPosition(interval []*BarsInterval) (float64, float64) {
	var UpMaxDiffSum float64 = 0
	var DownMaxDiffSum float64 = 0
	for _, barInterval := range interval {
		upMaxDiffArray, dowMaxDiffArray := DiffToArray(barInterval)
		// upMaxDiff := FindMaxDiff(upMaxDiffArray)
		// downMaxDiff := FindMaxDiff(dowMaxDiffArray)
		// UpMaxDiffSum = UpMaxDiffSum + upMaxDiff.DiffNum
		// DownMaxDiffSum = DownMaxDiffSum + downMaxDiff.DiffNum
		//////////////////////////////////////////////////////
		UpMaxDiffSum = UpMaxDiffSum + FindAverageDiff(upMaxDiffArray)
		DownMaxDiffSum = DownMaxDiffSum + FindAverageDiff(dowMaxDiffArray)
	}
	UpAverageMaxDiffSum := UpMaxDiffSum / float64(len(interval))
	DownAverageMaxDiffSum := DownMaxDiffSum / float64(len(interval))

	return UpAverageMaxDiffSum, DownAverageMaxDiffSum
}
func FindAverageDiff(diffs []UpDownDiffType) float64 {
	var averageDiff float64 = 0
	if len(diffs) == 0 {
		return averageDiff
	}
	for _, diff := range diffs {
		averageDiff = averageDiff + diff.DiffNum
	}
	return averageDiff / float64(len(diffs))
}

//	func FindMaxDiff(diffs []UpDownDiffType) UpDownDiffType {
//		var maxDiff UpDownDiffType
//		for _, diff := range diffs {
//			if diff.DiffNum > maxDiff.DiffNum {
//				maxDiff = diff
//			}
//		}
//		return maxDiff
//	}
func DiffToArray(barsInterval *BarsInterval) ([]UpDownDiffType, []UpDownDiffType) {
	upMaxDiffArray := []UpDownDiffType{}
	dowMaxDiffArray := []UpDownDiffType{}
	EachDiffs(barsInterval, func(diff *UpDownDiffType) {
		if diff.Side == alpaca.Buy {
			upMaxDiffArray = append(upMaxDiffArray, *diff)
		} else {
			dowMaxDiffArray = append(dowMaxDiffArray, *diff)
		}
	})
	// fmt.Println("REF_1: ", len(upMaxDiffArray), "REF_2: ", len(dowMaxDiffArray))
	return upMaxDiffArray, dowMaxDiffArray
}

func FindMaxAndMinDifferentia(barsInterval *BarsInterval) *MaxUpDownDiffs {
	upMaxDiffNum := &UpDownDiffType{
		DiffNum: math.Inf(-1),
	}
	dowMaxDiffNum := &UpDownDiffType{
		DiffNum: math.Inf(-1),
	}
	EachDiffs(barsInterval, func(diff *UpDownDiffType) {
		if diff.Side == alpaca.Buy {
			if diff.DiffNum > upMaxDiffNum.DiffNum {
				upMaxDiffNum = diff
			}
		} else {
			if diff.DiffNum > dowMaxDiffNum.DiffNum {
				dowMaxDiffNum = diff
			}
		}
	})

	return &MaxUpDownDiffs{
		UpMaxDiff:  upMaxDiffNum,
		DowMaxDiff: dowMaxDiffNum,
	}
}
func EachDiffs(barsInterval *BarsInterval, callback func(arg *UpDownDiffType)) {
	startIndex := barsInterval.startIndex
	endIndex := barsInterval.endIndex
	price := barsInterval.bars[0]
	interval := endIndex - startIndex
	plusInterval := barsInterval.barsList[int(math.Max(0, float64(startIndex-interval))):endIndex]
	priceIndex := utils.IndexOf(plusInterval, price)
	barsToOneNumValue := BarsToOneNumber(plusInterval)
	priceNum := BarsToOneNumber(barsInterval.bars)[0]
	plusIntervalPriceNum := barsToOneNumValue[priceIndex]
	plusIntervalNum := barsToOneNumValue[0:priceIndex]
	scale := priceNum / plusIntervalPriceNum
	// price/plus =diffNum/x
	//////////////////
	// upMaxDiffNum := math.Inf(-1)
	//////////////////////////////////
	// dowMaxDiffNum := math.Inf(-1)

	for _, num := range plusIntervalNum {
		min, max := utils.FindMinAndMax([]float64{plusIntervalPriceNum, num})
		diffNum := max - min
		side := alpaca.Buy
		if max == plusIntervalPriceNum {
			// dowMaxDiffNum = diffNum
			side = alpaca.Sell
		} else {
			side = alpaca.Buy
		}
		callback(&UpDownDiffType{
			Scale:   scale,
			DiffNum: diffNum,
			Side:    side,
		})
	}
}
