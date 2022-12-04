package market

import (
	"encoding/json"
	"fmt"
	"neural/db"
	"neural/utils"
	"sort"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type MarketData struct {
	client  marketdata.Client
	options marketdata.ClientOpts
	db      *gorm.DB
}

func (m MarketData) GetSyntheticAlpacaData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {

}

func (m MarketData) GetMarketData(symbol string, startTime time.Time, endTime time.Time) []marketdata.Bar {
	fmt.Println("REQUEST ALPACA BARS: \n", startTime, "\n", endTime)
	timeNow := time.Now()
	minute15 := int64(60 * 60 * 15)
	minEnd, _ := utils.FindMinAndMax([]int64{timeNow.Unix() - minute15, endTime.Unix()})
	quotes, err := m.client.GetBars(symbol, marketdata.GetBarsParams{
		TimeFrame:  marketdata.OneMin,
		Start:      startTime,
		End:        time.Unix(minEnd, 0),
		Adjustment: marketdata.Split,
		// AsOf:      "2022-06-10", // Leaving it empty yields the same results
	})
	if err != nil {
		panic(err)
	}
	return quotes
}
func (m MarketData) AlpacaBarToDbBar(symbol string, bar marketdata.Bar) db.Bar {
	barToJson, _ := json.MarshalIndent(bar, "", "  ")
	barStructToJson := string(barToJson)
	return db.Bar{
		Symbol:          symbol,
		Timestamp:       bar.Timestamp.Unix(),
		Open:            bar.Open,
		Close:           bar.Close,
		High:            bar.High,
		Low:             bar.Low,
		BarStructToJson: barStructToJson,
	}
}
func (m MarketData) AlpacaBarsToDbBars(symbol string, bars []marketdata.Bar) []db.Bar {
	var dbBars []db.Bar
	for _, bar := range bars {
		dbBars = append(dbBars, m.AlpacaBarToDbBar(symbol, bar))
	}
	return dbBars
}

// func (m MarketData) SaveBarsToDb(bars []db.Bar) []db.Bar {
// 	m.db.Create(&bars)
// 	return bars
// }

func (m MarketData) GetMarketDataFromDb(symbol string, startTime time.Time) []db.Bar {
	var Bars []db.Bar
	m.db.Where("symbol = ?", symbol).Where("timestamp >= ?", startTime.Unix()-2000).Find(&Bars)
	fmt.Println("DB BARS: ", len(Bars))
	return Bars
}

func (m MarketData) OptimizeBars(bars []db.Bar) []db.Bar {
	var newBars []db.Bar

	for _, bar := range bars {
		if slices.IndexFunc(newBars, func(el db.Bar) bool {
			return el.Timestamp == bar.Timestamp
		}) == -1 {
			newBars = append(newBars, bar)
		}
		// 2022-11-28 15:14:00 +0400 +04
		// 2022-11-28 20:29:00 +0400 +04
	}
	sort.Slice(newBars, func(index1, index2 int) bool {
		return newBars[index1].Timestamp < newBars[index2].Timestamp
	})
	fmt.Println("OPTIMIZE")
	barsCount := len(newBars)
	for index, _ := range newBars {
		if (index + 1) == barsCount {
			break
		}
		bar1 := bars[index]
		bar2 := bars[index+1]
		fmt.Println("OPTIMIZE BAR: ", time.Unix(bar1.Timestamp, 0))
		timeStampDiff := float64(bar2.Timestamp - bar1.Timestamp)
		if timeStampDiff != 60 {
			startTimeBar := time.Unix(bar1.Timestamp, 0)
			endTimeBar := time.Unix(bar2.Timestamp, 0)
			fmt.Println(
				"\nDIFF :",
				timeStampDiff,
				"\nBAR 1 TIME:",
				time.Unix(bar1.Timestamp, 0),
				// "\nSTART TIME 2:",
				// startTimeBar,
				"\nBAR 2 TIME:",
				time.Unix(bar2.Timestamp, 0),
				// "\nEND TIME 2:",
				// endTimeBar,
			)
			marketData := m.GetMarketData(bar1.Symbol, startTimeBar, endTimeBar)
			toDbBars := m.AlpacaBarsToDbBars(bar1.Symbol, marketData)
			fmt.Println("START BAR:")
			for _, bar := range toDbBars {
				fmt.Println("OPTIMIZE BAR: ", time.Unix(bar.Timestamp, 0))
			}
			fmt.Println("END BAR:")
			// panic("a problem")
		}
		// fmt.Println("TIMEST: ", time.Unix(bar.Timestamp, 0))
	}
	panic("a problem")
	fmt.Println("OPTIMIZE BARS: ", len(newBars))
	return newBars
}
func (m MarketData) SaveOnDb(bars []db.Bar) {
	for _, bar := range bars {
		m.db.Create(&bar)
	}
}
func (m MarketData) FillMarketBars(bars []db.Bar, symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("FILLABLE_BAR: ", len(bars))
	if len(bars) == 0 {
		marketData := m.GetMarketData(symbol, startTime, endTime)
		toDbBars := m.AlpacaBarsToDbBars(symbol, marketData)
		m.SaveOnDb(toDbBars)
		return toDbBars
	}
	var newBars []db.Bar
	// timeStampInHour := float64(60 * 60)
	// scaledTimeStampHour := timeStampInHour * 5
	barsCunt := len(bars)
	for _, bar := range bars {

		fmt.Println("SAVED TIMES: ", time.Unix(bar.Timestamp, 0))
	}
	for index, _ := range bars {
		if (index + 1) == barsCunt {
			break
		}
		bar1 := bars[index]
		bar2 := bars[index+1]
		timeStampDiff := float64(bar2.Timestamp - bar1.Timestamp)
		if timeStampDiff != 60 {
			// fmt.Println("PPPSSS: ", timeStampDiff)
			// plusSize := int64(timeStampDiff * 2)

			startTimeBar := time.Unix(bar1.Timestamp-int64(timeStampDiff), 0)
			weekday1 := time.Unix(bar1.Timestamp, 0).Weekday()
			weekday2 := time.Unix(bar2.Timestamp, 0).Weekday()
			endTimeBar := time.Unix(bar2.Timestamp+int64(timeStampDiff), 0)
			fmt.Println(
				"\nDIFF :",
				timeStampDiff,
				"\nBAR 1 TIME:",
				time.Unix(bar1.Timestamp, 0),
				"\nSTART TIME 2:",
				startTimeBar,
				"\nBAR 2 TIME:",
				time.Unix(bar2.Timestamp, 0),
				"\nEND TIME 2:",
				endTimeBar,
				"\nWEEK 1: ",
				weekday1,
				"\nWEEK 2: ",
				weekday2,
			)
			marketData := m.GetMarketData(symbol, startTimeBar, endTimeBar)
			toDbBars := m.AlpacaBarsToDbBars(symbol, marketData)
			m.SaveOnDb(toDbBars)
			fmt.Println(
				timeStampDiff,
				"\nEND BAR FIRST:",
				time.Unix(toDbBars[0].Timestamp, 0),
				"\nEND BAR LAST:",
				time.Unix(toDbBars[len(toDbBars)-1].Timestamp, 0),
			)
			for _, bar := range toDbBars {

				fmt.Println("TIMES: ", time.Unix(bar.Timestamp, 0))
			}
			// bars []db.Bar, symbol string, startTime time.Time, endTime time.Time
			for _, toDbBar := range toDbBars {
				index := slices.IndexFunc(bars, func(el db.Bar) bool {
					return toDbBar.Timestamp == el.Timestamp
				})
				utils.LogStruct("RESULT: ", time.Unix(toDbBar.Timestamp, 0), index)
			}
			// return m.FillMarketBars(
			// 	m.OptimizeBars(append(bars, toDbBars...)),
			// 	symbol,
			// 	startTimeBar,
			// 	endTimeBar,
			// )

			return m.GetMarketCachedData(symbol, startTime, endTime)
		}
	}
	return m.OptimizeBars(newBars)
}
func (m MarketData) GetMarketCachedData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {

	barsFromDb := m.OptimizeBars(m.GetMarketDataFromDb(symbol, startTime))
	// fmt.Println("🚀 --> file: martket-data.go:111 --> func --> barsFromDb", barsFromDb)

	filedBars := m.FillMarketBars(barsFromDb, symbol, startTime, endTime)

	return filedBars
}
func (m MarketData) GetMarketCachedDataWithFrame(hourFrame float64, symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	var newBars []db.Bar
	bars := m.GetMarketCachedData(symbol, startTime, endTime)
	frameTimeStampInHour := int64(float64(60*60) * hourFrame)
	for timeStamp := startTime.Unix(); timeStamp < endTime.Unix(); timeStamp += frameTimeStampInHour {
		closestBar := bars[0]
		for _, bar := range bars {
			min1, max1 := utils.FindMinAndMax([]int64{closestBar.Timestamp, timeStamp})
			min2, max2 := utils.FindMinAndMax([]int64{bar.Timestamp, timeStamp})
			if (max1 - min1) > (max2 - min2) {
				closestBar = bar
			}
		}
		newBars = append(newBars, closestBar)
	}
	return newBars
}

// func (m MarketData) GetMarketData(symbol string, timeFrame marketdata.TimeFrame, startDate time.Time) []marketdata.Bar {
// 	client := m.client
// 	// sss, ret := client.GetLatestBars([]string{"TWTR"})
// 	// fmt.Println(sss)
// 	// fmt.Println("LLL", ret)
// 	// // marketdata.GetLatestBars
// 	// // marketdata.StreamTradeUpdatesInBackground(context.TODO(), func(tu alpaca.TradeUpdate) {
// 	// // 	log.Printf("TRADE UPDATE: %+v\n", tu)
// 	// // })
// 	bars, err := client.GetBars(symbol, marketdata.GetBarsParams{
// 		TimeFrame: timeFrame,
// 		Start:     startDate,
// 		End:       time.Now(),
// 	})
// 	if err != nil {
// 		fmt.Printf("GET '%v' BARS:\n", err)
// 		panic(err)
// 	}
// 	fmt.Printf("GET '%v' BARS:\n", symbol)
// 	for _, bar := range bars {
// 		// bar.High
// 		fmt.Printf("%+v\n", bar)
// 	}
// 	return bars
// }

// func DbBarsToCharBars(bars []db.Bar) []financeGo.ChartBar {
// 	fmt.Println("DbBarsToCharBars")

// 	charBars := []financeGo.ChartBar{}
// 	for _, charBar := range bars {
// 		charBars = append(charBars, financeGo.ChartBar{
// 			Open:      charBar.Open,
// 			Low:       charBar.Low,
// 			High:      charBar.High,
// 			Close:     charBar.Close,
// 			AdjClose:  charBar.AdjClose,
// 			Volume:    charBar.Volume,
// 			Timestamp: charBar.Timestamp,
// 		})
// 	}

// 	uniqueCharBars := []financeGo.ChartBar{}

// 	for _, bat := range charBars {
// 		index := slices.IndexFunc(uniqueCharBars, func(el financeGo.ChartBar) bool {
// 			return el.Timestamp == bat.Timestamp
// 		})
// 		if index == -1 {
// 			uniqueCharBars = append(uniqueCharBars, bat)
// 		}

// 	}
// 	if len(charBars) > 0 {
// 		fmt.Println("CharBars1:", time.Unix(int64(charBars[0].Timestamp), 0))
// 		fmt.Println("CharBars2:", time.Unix(int64(charBars[len(charBars)-1].Timestamp), 0))
// 	}
// 	sort.Slice(uniqueCharBars, func(index1, index2 int) bool {
// 		return uniqueCharBars[index1].Timestamp < uniqueCharBars[index2].Timestamp
// 	})
// 	fmt.Println("DB UNIQUE BARS COUNT:   ", len(uniqueCharBars))
// 	return uniqueCharBars
// }

// func GetMarketDataDb(symbol string, startDate time.Time) []financeGo.ChartBar {
// 	fmt.Println("GetMarketDataDb")
// 	var Bars []db.Bar
// 	db.Database.Where("symbol = ?", symbol).Where("timestamp >= ?", startDate.Unix()).Find(&Bars)
// 	fmt.Println("DB BARS COUNT:   ", len(Bars))

// 	charBars := DbBarsToCharBars(Bars)
// 	return charBars
// }
// func SaveBarsOnDb(symbol string, financeBars []*financeGo.ChartBar) {
// 	fmt.Println("SaveBarsOnDb")
// 	for _, bar := range financeBars {
// 		dbBar := db.Bar{
// 			Symbol: symbol,
// 		}
// 		dbBar.Open = bar.Open
// 		dbBar.Low = bar.Low
// 		dbBar.High = bar.High
// 		dbBar.Close = bar.Close
// 		dbBar.AdjClose = bar.AdjClose
// 		dbBar.Volume = bar.Volume
// 		dbBar.Timestamp = bar.Timestamp
// 		db.Database.Create(&dbBar)
// 	}
// }
// func GetBarsAndSave(symbol string, startTimeStamp float64, endTimeStamp float64) []*financeGo.ChartBar {
// 	fmt.Println("GetBarsAndSave")
// 	cutTimeStamp := float64(1669220040 - 1669047240)
// 	bars := []*financeGo.ChartBar{}
// 	// fmt.Println("QWWW4", time.Unix(int64(startTimeStamp+cutTimeStamp), 0), time.Unix(int64(endTimeStamp), 0))
// 	// fmt.Println("QWWW3", time.Unix(int64(startTimeStamp), 0), time.Unix(int64(endTimeStamp), 0))
// 	for timeStamp := startTimeStamp; timeStamp < endTimeStamp; timeStamp += cutTimeStamp {

// 		startTimeDate := time.Unix(int64(timeStamp), 0)
// 		endTimeDate := time.Unix(int64(timeStamp+cutTimeStamp), 0)
// 		// fmt.Println("REQUEST INTERVAL", startTimeDate, endTimeDate)
// 		// fmt.Println("QWWW2", startTimeDate, endTimeDate)

// 		params := &chart.Params{
// 			Symbol:   symbol,
// 			Start:    datetime.New(&startTimeDate),
// 			End:      datetime.New(&endTimeDate),
// 			Interval: datetime.OneMin,
// 		}
// 		// time.Sleep(4000 * time.Second)
// 		iter := chart.Get(params)
// 		fmt.Println("REQUEST INTERVAL", startTimeDate, endTimeDate)
// 		utils.LogStruct("FinanceIterToArray")
// 		requestBars := finance.FinanceIterToArray(iter)
// 		utils.LogStruct("SSSSSSSSSSSSSSSSSSSSSSSSSSS", len(requestBars), startTimeDate, endTimeDate)
// 		SaveBarsOnDb(symbol, requestBars)
// 		time.Sleep(time.Second * 5)
// 		bars = append(bars, requestBars...)

// 	}
// 	// fmt.Println("SaveBarsOnDb")
// 	fmt.Println("bars")
// 	return bars
// }
// func GetValidBars(symbol string, startDate time.Time) []financeGo.ChartBar {
// 	bars := GetMarketDataDb(symbol, startDate)
// 	diff := 0
// 	count := len(bars)
// 	if count == 0 {
// 		GetBarsAndSave(symbol, float64(startDate.Unix()-3000), float64(time.Now().Unix()))
// 		return GetValidBars(symbol, startDate)
// 	}
// 	for index, bar1 := range bars {
// 		if diff == 0 {
// 			if startDate.Unix() < int64(bar1.Timestamp) {
// 				fmt.Println("QWWW1", startDate, time.Unix(int64(bar1.Timestamp), 0))
// 				// return []financeGo.ChartBar{}
// 				// time.Sleep(10 * time.Second)
// 				GetBarsAndSave(symbol, float64(startDate.Unix()-3000), float64(time.Now().Unix()))
// 				return GetValidBars(symbol, startDate)
// 			}
// 			continue
// 		}
// 		if (index + 1) < count {
// 			continue
// 		}
// 		bar2 := bars[index+1]
// 		newDiff := bar2.Timestamp - bar1.Timestamp
// 		fmt.Println("newDiff", newDiff)
// 		if newDiff != diff {
// 			// GetMarketDataDb(symbol, time.Unix(int64(bar1.Timestamp), 0))
// 			GetBarsAndSave(symbol, float64(bar1.Timestamp), float64(time.Now().Unix()))
// 			return GetValidBars(symbol, startDate)
// 		}
// 	}
// 	return bars
// }
// func GetAndFillEmptyIntervals(symbol string, frameHour float64, startDate time.Time) []financeGo.ChartBar {
// fmt.Println("GetAndFillEmptyIntervals")
// bars := GetValidBars(symbol, startDate)
// fmt.Println("RESULT LEN: ", len(bars))
// resultBars := []financeGo.ChartBar{}
// startTimeStamp := float64(startDate.Unix())
// endTimeStamp := float64(time.Now().Unix())
// fmt.Println("EEEEE", endTimeStamp, time.Now().Unix())
// index := 0
// count := len(bars)
// frameInTimeStamp := frameHour * float64(60*60)
// frameIndex := 0
// if count == 0 {
// 	GetBarsAndSave(symbol, startTimeStamp, endTimeStamp)
// 	return GetAndFillEmptyIntervals(symbol, frameHour, startDate)
// }

// fmt.Println(
// 	"\n FIRST Timestamp:",
// 	time.Unix(int64(bars[0].Timestamp), 0),
// )
// for {
// 	bar := bars[index]
// 	plusInterval := float64(frameIndex) * frameInTimeStamp
// 	startIndexTimeStamp := startTimeStamp + plusInterval
// 	endIndexTimeStamp := startTimeStamp + plusInterval + frameInTimeStamp
// 	barTimestamp := float64(bar.Timestamp)
// 	ifP := barTimestamp >= startIndexTimeStamp && barTimestamp <= endIndexTimeStamp

// 	if ifP {
// 		resultBars = append(resultBars, bar)
// 		frameIndex++
// 	}
// 	// if endIndexTimeStamp < endTimeStamp {
// 	// 	return GetAndFillEmptyIntervals(symbol, frameHour, startDate)
// 	// }
// 	index++
// 	if index >= count {
// 		if endIndexTimeStamp < endTimeStamp {

// 			fmt.Println(
// 				"TIME",
// 				"\n START Plus:",
// 				time.Unix(int64(startIndexTimeStamp), 0),
// 				"\n END Plus:",
// 				time.Unix(int64(endIndexTimeStamp), 0),
// 				"\n BAR: ",
// 				time.Unix(int64(bar.Timestamp), 0),
// 				"\n END OF END: ",
// 				time.Unix(int64(endTimeStamp), 0),
// 			)
// 			// GetBarsAndSave(symbol, startIndexTimeStamp, endTimeStamp)
// 			// return GetAndFillEmptyIntervals(symbol, frameHour, startDate)
// 		}
// 		break
// 	}
// }
// 	return bars
// }
// func GetBars(symbol string, frameHour []float64, startDate time.Time) {
// 	fmt.Println("GetBars")
// 	for _, frame := range frameHour {
// 		bars := GetAndFillEmptyIntervals(symbol, frame, startDate)
// 		for index, bar := range bars {
// 			if index == 0 {
// 				continue
// 			}
// 			bar2 := bars[index-1]
// 			fmt.Println("DIFF: ", bar.Timestamp-bar2.Timestamp)
// 		}
// 	}

// }
