package finance

import (
	"fmt"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func GetStockData(symbol string, interval datetime.Interval, startDate datetime.Datetime) *chart.Iter {
	timeNow := time.Now()
	fmt.Println(datetime.New(&timeNow))
	params := &chart.Params{
		Symbol: symbol,
		Start:  &startDate,
		// End:    &datetime.Datetime{Month: 12, Day: 22, Year: 2022},
		End:      datetime.New(&timeNow),
		Interval: interval,
	}
	iter := chart.Get(params)
	// fmt.Println(iter.Meta())
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}
	return iter
}