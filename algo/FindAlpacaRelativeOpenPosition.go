package algo

import (
	"fmt"
	"math"
	"neural/utils"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

type AlpacaPosition struct {
	Side       alpaca.Side
	StopLost   float64
	TakeProfit float64
}

func FindAlpacaRelativeOpenPosition(intervalBars TimeIntervalsBars) {
	intervalSimilarity := intervalBars.intervalSimilarity
	maxDiffPosition := FindAverageDiffMaxPosition(intervalSimilarity)
	fmt.Println("maxDiffPosition", maxDiffPosition)
	// var finalMaxNum float64 = math.Inf(-1)
	// var finalMinNum float64 = math.Inf(1)
	// intervalSimilarity := intervalBars.intervalSimilarity
	// currentPrice, _ := intervalSimilarity[0].bars[0].Close.Float64()
	// diffDetails := FindSimilarityIntervalMaxAndMinDifferentials(intervalSimilarity)

	// diffPosition := FindRelativePosition(diffDetails)
	// position := AlpacaPosition{}
	// if diffPosition.Side == alpaca.Buy {
	// 	position = AlpacaPosition{
	// 		Side:       diffPosition.Side,
	// 		TakeProfit: (currentPrice + currentPrice*(diffPosition.TakeProfit/1)),
	// 		StopLost:   (currentPrice - currentPrice*(diffPosition.StopLost/1)),
	// 	}
	// } else {
	// 	position = AlpacaPosition{
	// 		Side:       diffPosition.Side,
	// 		TakeProfit: (currentPrice - currentPrice*(diffPosition.TakeProfit/1)),
	// 		StopLost:   (currentPrice + currentPrice*(diffPosition.StopLost/1)),
	// 	}
	// }
	// diff := diffDetails.Diffs
	// if diff.UpMaxDiffNum > diff.DowMaxDiffNum {
	// 	return AlpacaPosition{
	// 		TakeProfit: (currentPrice + currentPrice*(diffDetails.UpMaxAverageDiffNum/1)),
	// 		StopLost:   (currentPrice - currentPrice*(diff.DowMaxDiffNum/1)),
	// 	}
	// } else {
	// 	return AlpacaPosition{
	// 		TakeProfit: (currentPrice - currentPrice*(diffDetails.DowMaxAverageDiffNum/1)),
	// 		StopLost:   (currentPrice + currentPrice*(diff.UpMaxDiffNum/1)),
	// 	}
	// }
	// fmt.Println(diff)
	// fmt.Printf("POSITON: %+v\n", position)
	// return position
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

func FindAverageDiffMaxPosition(interval []*BarsInterval) *UpDownDiffType {
	var maxDiff *UpDownDiffType
	for _, barInterval := range interval {
		upMaxDiffArray, dowMaxDiffArray := DiffToArray(barInterval)
		upMaxDiff := FindMaxDiff(upMaxDiffArray)
		downMaxDiff := FindMaxDiff(dowMaxDiffArray)

		if upMaxDiff.DiffNum > downMaxDiff.DiffNum && upMaxDiff.DiffNum > maxDiff.DiffNum {
			maxDiff = upMaxDiff
		} else if downMaxDiff.DiffNum > upMaxDiff.DiffNum && downMaxDiff.DiffNum > maxDiff.DiffNum {
			maxDiff = downMaxDiff
		}
	}

	return maxDiff
}
func FindMaxDiff(diffs []*UpDownDiffType) *UpDownDiffType {
	maxDiff := diffs[0]
	for _, diff := range diffs {
		if diff.DiffNum > maxDiff.DiffNum {
			maxDiff = diff
		}
	}
	return maxDiff
}
func DiffToArray(barsInterval *BarsInterval) ([]*UpDownDiffType, []*UpDownDiffType) {
	upMaxDiffArray := []*UpDownDiffType{}
	dowMaxDiffArray := []*UpDownDiffType{}
	EachDiffs(barsInterval, func(diff *UpDownDiffType) {
		if diff.Side == alpaca.Buy {
			upMaxDiffArray = append(upMaxDiffArray, diff)
		} else {
			dowMaxDiffArray = append(dowMaxDiffArray, diff)
		}
	})
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
