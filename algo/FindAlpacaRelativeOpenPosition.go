package algo

import (
	"math"
	"neural/utils"
)

type AlpacaPosition struct {
	StopLost   float64
	TakeProfit float64
}

func FindAlpacaRelativeOpenPosition(intervalBars TimeIntervalsBars) AlpacaPosition {
	// var finalMaxNum float64 = math.Inf(-1)
	// var finalMinNum float64 = math.Inf(1)
	intervalSimilarity := intervalBars.intervalSimilarity
	currentPrice, _ := intervalSimilarity[0].bars[0].Close.Float64()
	diffDetails := FindMaxAndMinSimilarityIntervalDifferentia(intervalSimilarity)
	diff := diffDetails.Diffs
	if diff.UpMaxDiffNum > diff.DowMaxDiffNum {
		return AlpacaPosition{
			TakeProfit: (currentPrice + currentPrice*(diffDetails.UpMaxAverageDiffNum/1)),
			StopLost:   (currentPrice - currentPrice*(diff.DowMaxDiffNum/1)),
		}
	} else {
		return AlpacaPosition{
			TakeProfit: (currentPrice - currentPrice*(diffDetails.DowMaxAverageDiffNum/1)),
			StopLost:   (currentPrice + currentPrice*(diff.UpMaxDiffNum/1)),
		}
	}
	// fmt.Println(diff)
}

// ბარების საშუალო გამოთვალე ინტერვალით
type UpDownMaxDiffType struct {
	UpMaxDiffNum        float64
	UpMaxDiffTimestamp  int
	DowMaxDiffNum       float64
	DowMaxDiffTimestamp int
}

// ბარების საშუალო გამოთვალე ინტერვალით
type UpDownMaxDiffWithAverageType struct {
	Diffs                UpDownMaxDiffType
	UpMaxAverageDiffNum  float64
	DowMaxAverageDiffNum float64
}

func FindMaxAndMinSimilarityIntervalDifferentia(intervalSimilarity []*BarsInterval) *UpDownMaxDiffWithAverageType {
	upMaxDiffNum := math.Inf(-1)
	upMaxDiffTimestamp := 0
	upMaxDiffNumSum := float64(0)
	dowMaxDiffNum := math.Inf(-1)
	dowMaxDiffTimestamp := 0
	dowMaxDiffNumSum := float64(0)
	for index, barsInterval := range intervalSimilarity {
		if index != 0 {
			maxDiffs := FindMaxAndMinDifferentia(barsInterval)
			if maxDiffs.UpMaxDiffNum > upMaxDiffNum {
				upMaxDiffNum = maxDiffs.UpMaxDiffNum
				upMaxDiffNumSum = upMaxDiffNumSum + upMaxDiffNum
				upMaxDiffTimestamp = maxDiffs.UpMaxDiffTimestamp
			}
			if maxDiffs.DowMaxDiffNum > dowMaxDiffNum {
				dowMaxDiffNum = maxDiffs.DowMaxDiffNum
				dowMaxDiffNumSum = dowMaxDiffNumSum + dowMaxDiffNum
				dowMaxDiffTimestamp = maxDiffs.DowMaxDiffTimestamp
			}
		}
	}
	return &UpDownMaxDiffWithAverageType{
		Diffs: UpDownMaxDiffType{
			UpMaxDiffNum:        upMaxDiffNum,
			UpMaxDiffTimestamp:  upMaxDiffTimestamp,
			DowMaxDiffNum:       dowMaxDiffNum,
			DowMaxDiffTimestamp: dowMaxDiffTimestamp,
		},
		UpMaxAverageDiffNum:  upMaxDiffNumSum / float64(len(intervalSimilarity)-1),
		DowMaxAverageDiffNum: dowMaxDiffNumSum / float64(len(intervalSimilarity)-1),
	}
}
func FindMaxAndMinDifferentia(barsInterval *BarsInterval) *UpDownMaxDiffType {
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
	//console.log(4/2 = 8/4 == 2)
	//console.log(4/2 = 8/4 == 2)
	// p/pl=x/d =
	upMaxDiffNum := math.Inf(-1)
	upMaxDiffTimestamp := 0
	dowMaxDiffNum := math.Inf(-1)
	dowMaxDiffTimestamp := 0

	for index, num := range plusIntervalNum {
		min, max := utils.FindMinAndMax([]float64{plusIntervalPriceNum, num})
		timestamp := plusInterval[index].Timestamp
		diffNum := max - min
		if max == plusIntervalPriceNum {
			dowMaxDiffNum = diffNum * scale
			dowMaxDiffTimestamp = timestamp
		} else {
			upMaxDiffNum = diffNum * scale
			upMaxDiffTimestamp = timestamp
		}

	}
	return &UpDownMaxDiffType{
		UpMaxDiffNum:        upMaxDiffNum,
		UpMaxDiffTimestamp:  upMaxDiffTimestamp,
		DowMaxDiffNum:       dowMaxDiffNum,
		DowMaxDiffTimestamp: dowMaxDiffTimestamp,
	}
}
