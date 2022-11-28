package market

import (
	"fmt"
	"neural/db"
	"neural/finance"
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

func DbBarsToCharBars(bars []db.Bar) []financeGo.ChartBar {
	sort.Slice(bars, func(index1, index2 int) bool {
		return bars[index1].Timestamp < bars[index2].Timestamp
	})

	charBars := []financeGo.ChartBar{}
	for _, charBar := range bars {
		charBars = append(charBars, financeGo.ChartBar{
			Open:      charBar.Open,
			Low:       charBar.Low,
			High:      charBar.High,
			Close:     charBar.Close,
			AdjClose:  charBar.AdjClose,
			Volume:    charBar.Volume,
			Timestamp: charBar.Timestamp,
		})
	}
	return charBars
}

func GetMarketDataDb(symbol string, startDate time.Time) []financeGo.ChartBar {
	var Bars []db.Bar
	db.Database.Where("symbol = ?", symbol).Where("timestamp >= ?", startDate.Unix()).Find(&Bars)
	fmt.Println("BARS ,", Bars)

	// charBars :=
	return DbBarsToCharBars(Bars)
}
func SaveBarsOnDb(symbol string, financeBars []*financeGo.ChartBar) {
	for _, bar := range financeBars {
		dbBar := db.Bar{
			Symbol: symbol,
		}
		dbBar.Open = bar.Open
		dbBar.Low = bar.Low
		dbBar.High = bar.High
		dbBar.Close = bar.Close
		dbBar.AdjClose = bar.AdjClose
		dbBar.Volume = bar.Volume
		dbBar.Timestamp = bar.Timestamp

		db.Database.Create(&dbBar)
	}
}
func GetBarsAndSave(symbol string, startTimeStamp float64, endTimeStamp float64) []*financeGo.ChartBar {
	cutTimeStamp := float64(250 * 60)
	bars := []*financeGo.ChartBar{}
	for timeStamp := startTimeStamp; timeStamp < endTimeStamp; timeStamp += cutTimeStamp {
		startTimeDate := time.Unix(int64(timeStamp), 0)
		endTimeDate := time.Unix(int64(timeStamp+cutTimeStamp), 0)
		params := &chart.Params{
			Symbol:   symbol,
			Start:    datetime.New(&startTimeDate),
			End:      datetime.New(&endTimeDate),
			Interval: datetime.OneMin,
		}
		iter := chart.Get(params)
		bars = append(bars, finance.FinanceIterToArray(iter)...)

	}
	SaveBarsOnDb(symbol, bars)
	return bars
}
func GetAndFillEmptyIntervals(symbol string, frameHour float64, startDate time.Time) []financeGo.ChartBar {
	bars := GetMarketDataDb(symbol, startDate)
	resultBars := []financeGo.ChartBar{}
	startTimeStamp := float64(startDate.Unix())
	endTimeStamp := float64(time.Now().Unix())
	index := 0
	count := len(bars)
	frameInTimeStamp := frameHour * float64(60)
	frameIndex := 0
	for {
		bar := bars[index]
		plusInterval := float64(frameIndex) * frameInTimeStamp
		startIndexTimeStamp := startTimeStamp + plusInterval
		endIndexTimeStamp := startTimeStamp + plusInterval + frameInTimeStamp
		barTimestamp := float64(bar.Timestamp)
		if barTimestamp > startIndexTimeStamp && barTimestamp < endIndexTimeStamp {
			resultBars = append(resultBars, bar)
			frameIndex++
		}

		index++
		if index >= count {
			if endIndexTimeStamp < endTimeStamp {
				GetBarsAndSave(symbol, startIndexTimeStamp, endTimeStamp)
				return GetAndFillEmptyIntervals(symbol, frameHour, startDate)
			}
			break
		}
	}
	return resultBars
}
func GetBars(symbol string, frameHour []float64, startDate time.Time) {
	for _, frame := range frameHour {
		bars := GetAndFillEmptyIntervals(symbol, frame, startDate)
		for index, bar := range bars {
			if index == 0 {
				continue
			}
			bar2 := bars[index-1]
			fmt.Println("DIFF: ", bar.Timestamp-bar2.Timestamp)
		}
	}

}
