package algo

func SumDistanceSimilarityNum(barIntervals []*BarsInterval) float64 {
	var sumNum float64 = 0
	for _, interval := range barIntervals {
		sumNum = sumNum + (interval.similarityNum / float64(len(interval.bars)))
	}
	// fmt.Println("SUM: ", sumNum/float64(len(barIntervals)))
	return sumNum / float64(len(barIntervals))
}
