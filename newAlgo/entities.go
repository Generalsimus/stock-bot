package newAlgo

import financeGo "github.com/piquette/finance-go"

type SymbolBestTimeIntervalsBars struct {
	Symbol        string
	BestIntervals []TimeIntervalsBars
}
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
