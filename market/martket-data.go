package market

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type MarketData struct {
	client  marketdata.Client
	options marketdata.ClientOpts
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
