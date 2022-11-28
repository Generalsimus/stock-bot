package market

import (
	"fmt"
	"neural/db"
	"sort"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	financeGo "github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

type MarketData struct {
	client  marketdata.Client
	options marketdata.ClientOpts
}

func (m MarketData) GetMarketData(symbol string, timeFrame marketdata.TimeFrame, startDate time.Time) []marketdata.Bar {
	client := m.client
	// sss, ret := client.GetLatestBars([]string{"TWTR"})
	// fmt.Println(sss)
	// fmt.Println("LLL", ret)
	// // marketdata.GetLatestBars
	// // marketdata.StreamTradeUpdatesInBackground(context.TODO(), func(tu alpaca.TradeUpdate) {
	// // 	log.Printf("TRADE UPDATE: %+v\n", tu)
	// // })
	bars, err := client.GetBars(symbol, marketdata.GetBarsParams{
		TimeFrame: timeFrame,
		Start:     startDate,
		End:       time.Now(),
	})
	if err != nil {
		fmt.Printf("GET '%v' BARS:\n", err)
		panic(err)
	}
	fmt.Printf("GET '%v' BARS:\n", symbol)
	for _, bar := range bars {
		// bar.High
		fmt.Printf("%+v\n", bar)
	}
	return bars
}

// func DbBarsToCharBars(symbol string, hourFrame int, startDate time.Time) {

// }

func GetMarketDataDb(symbol string, hourFrame float64, startDate time.Time) []financeGo.ChartBar {
	var Bars []db.Bar
	db.Database.Where("symbol = ?", symbol).Find(&Bars)
	fmt.Println("BARS ,", Bars)

	charBars := DbBarsToCharBars(Bars, hourFrame, startDate)
	return charBars
}
func GetBarsAndSave(symbol string, startTimeStamp float64, endTimeStamp float64) {
	cutTimeStamp := float64(250 * 60)
	for timeStamp := startTimeStamp; timeStamp < endTimeStamp; timeStamp += cutTimeStamp {
		startTimeDate := time.Unix(startTime, 0)
		endTimeDate := time.Unix(endTime, 0)
		params := &chart.Params{
			Symbol:   symbol,
			Start:    datetime.New(&startTime),
			End:      datetime.New(&endTime),
			Interval: interval,
		}
		iter := chart.Get(params)
	}
	startTime := time.Unix(int64(startTimeStamp), 0)
	endTime := time.Unix(int64(endTimeStamp), 0)
	params := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&startTime),
		End:      datetime.New(&endTime),
		Interval: interval,
	}
	iter := chart.Get(params)
	// financeBars := finance.GetSymbolIntervalBars(symbol, datetime.OneMin, datetime.Datetime{Month: int(startDate.Month()), Day: startDate.Day(), Year: startDate.Year()})

}
func GetMinTimeInterval() {
}

func DbBarsToCharBars(bars []db.Bar, hourFrame float64, startDate time.Time) []financeGo.ChartBar {
	sort.Slice(bars, func(index1, index2 int) bool {
		return bars[index1].Timestamp < bars[index2].Timestamp
	})
	// for {
	// financeBars := finance.GetSymbolIntervalBars(symbol, datetime.OneMin, datetime.Datetime{Month: int(startDate.Month()), Day: startDate.Day(), Year: startDate.Year()})
	hourTimestamp := float64(3600)
	frameToTimestamp := hourFrame * hourTimestamp
	endTimestamp := float64(time.Now().Unix())
	timestampIntervals := float64(60)
	if frameToTimestamp < timestampIntervals {
		panic("Time Frame Is Low")
	}
	fmt.Println("STARTTTT")
	barIndex := 0
	eachTimestamp := float64(startDate.Unix())
	for eachTimestamp < endTimestamp {
		bar := bars[barIndex]
		baTimestamp := float64(bar.Timestamp)
		if eachTimestamp < baTimestamp {

		}
		barIndex++
		eachTimestamp += frameToTimestamp
	}
	// for _, bar := range financeBars {
	// 	dbBar := db.Bar{
	// 		Symbol: symbol,
	// 	}
	// 	dbBar.Open = bar.Open
	// 	dbBar.Low = bar.Low
	// 	dbBar.High = bar.High
	// 	dbBar.Close = bar.Close
	// 	dbBar.AdjClose = bar.AdjClose
	// 	dbBar.Volume = bar.Volume
	// 	dbBar.Timestamp = bar.Timestamp

	// 	db.Database.Create(&dbBar)
	// }

	// }
	charBars := []financeGo.ChartBar{}
	return charBars
}
