package finance

import (
	financeGo "github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func GetSymbolIntervalBars(symbol string, interval datetime.Interval, startDate datetime.Datetime) []*financeGo.ChartBar {
	iter := GetStockData(symbol, interval, startDate)
	// iter.

	return FinanceIterToArray(iter)
}
func FinanceIterToArray(iter *chart.Iter) []*financeGo.ChartBar {
	output := []*financeGo.ChartBar{}
	// output := []float64{8, 7, 1, 2, 5}
	for iter.Next() {
		Bar := iter.Bar()

		output = append(output, Bar)

	}
	return output
}
