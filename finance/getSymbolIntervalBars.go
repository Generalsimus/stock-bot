package finance

import (
	"fmt"
	"time"

	financeGo "github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func GetSymbolIntervalBars(symbol string, interval datetime.Interval, startDate datetime.Datetime) []*financeGo.ChartBar {
	iter := GetStockData(symbol, interval, startDate)
	// iter.

	return FinanceIterToArray(iter)
}

func GetStockDataWithSymbolIntervalBars(symbol string, interval datetime.Interval) []*financeGo.ChartBar {
	iter := GetStockDataWithSymbolInterval(symbol, interval)
	return FinanceIterToArray(iter)
}

var FrameMaxIntervals map[datetime.Interval]uint

func GetMaxIntervalBars(symbol string, interval datetime.Interval) []*financeGo.ChartBar {
	timeNow := time.Now()
	timeNowTimeStamp := timeNow.Unix()
	nowDate := datetime.New(&timeNow)
	if maxTimeStamp, ok := FrameMaxIntervals[interval]; ok {
		maxStartTime := time.Unix(int64(maxTimeStamp), 0)
		return FinanceIterToArray(chart.Get(&chart.Params{
			Symbol:   symbol,
			Interval: interval,
			Start:    datetime.New(&maxStartTime),
			End:      nowDate,
		}))
	}
	dayTimeStamp := int64(60 * 60 * 24)
	timeStamp := timeNowTimeStamp - dayTimeStamp
	valueBars := []*financeGo.ChartBar{}
	for {
		startTime := time.Unix(int64(timeStamp), 0)
		bars := FinanceIterToArray(chart.Get(&chart.Params{
			Symbol:   symbol,
			Interval: interval,
			Start:    datetime.New(&startTime),
			End:      nowDate,
		}))
		fmt.Println("LEN: ", len(bars), dayTimeStamp, "\n", startTime, timeNowTimeStamp, timeStamp)
		if len(bars) == 0 {
			return valueBars
		}
		valueBars = bars
		timeStamp -= dayTimeStamp
	}
	return valueBars
}

func FinanceIterToArray(iter *chart.Iter) []*financeGo.ChartBar {
	output := []*financeGo.ChartBar{}
	// output := []float64{8, 7, 1, 2, 5}
	for iter.Next() {
		Bar := iter.Bar()
		// utils.LogStruct(Bar)
		output = append(output, Bar)

	}
	return output
}
