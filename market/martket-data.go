package market

import (
	"encoding/json"
	"errors"
	"fmt"
	"neural/db"
	"neural/utils"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	financeGo "github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MarketData struct {
	client  marketdata.Client
	options marketdata.ClientOpts
	db      *gorm.DB
}

func (m MarketData) GetYahooFinanceData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("GetYahooFinanceData")
	params := &chart.Params{
		Symbol:   symbol,
		Interval: datetime.OneMin,
		Start:    datetime.New(&startTime),
		End:      datetime.New(&endTime),
	}
	iter := chart.Get(params)
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}
	var dbBars []db.Bar
	for iter.Next() {
		bar := iter.Bar()
		dbBars = append(dbBars, m.YahooBarToDbBar(symbol, *bar))
	}
	// dbBars := m.OptimizeBars(bars)
	m.SaveBarsOnDb(dbBars)
	return dbBars
}
func (m MarketData) GetAlpacaMarketData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("GetAlpacaMarketData")
	fmt.Println("REQUEST ALPACA BARS: \n", startTime, "\n", endTime, "\n", symbol)
	timeNow := time.Now()
	minute15 := int64(60 * 16)
	minEnd, _ := utils.FindMinAndMax([]int64{timeNow.Unix() - minute15, endTime.Unix()})
	quotes, err := m.client.GetBars(symbol, marketdata.GetBarsParams{
		TimeFrame:  marketdata.OneMin,
		Start:      startTime,
		End:        time.Unix(minEnd, 0),
		Adjustment: marketdata.Split,
		// TotalLimit: 5000,
		// AsOf:       "2022-06-10", // Leaving it empty yields the same results
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("RESPONSE ALPACA BARS: \n", quotes[0].Timestamp, "\n", quotes[len(quotes)-1].Timestamp)
	fmt.Println("GET ALPACA BARS: ", len(quotes))
	dbBars := m.AlpacaBarsToDbBars(symbol, quotes)
	m.SaveBarsOnDb(dbBars)
	return dbBars
}
func (m MarketData) SaveBarsOnDb(bars []db.Bar) []db.Bar {
	fmt.Println("SaveBarsOnDb")
	barsCount := len(bars)
	if barsCount != 0 {
		for _, bar := range bars {
			m.db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "timestamp"}},
				UpdateAll: true,
			}).Create(&bar)
		}
	}

	return bars
}
func (m MarketData) AlpacaBarToDbBar(symbol string, bar marketdata.Bar) db.Bar {
	// fmt.Println("AlpacaBarToDbBar")
	barToJson, _ := json.MarshalIndent(bar, "", "  ")
	barStructToJson := string(barToJson)
	dbBar := db.Bar{
		Symbol:          symbol,
		Timestamp:       bar.Timestamp.Unix(),
		Open:            bar.Open,
		Close:           bar.Close,
		High:            bar.High,
		Low:             bar.Low,
		BarStructToJson: barStructToJson,
	}
	return dbBar
}
func (m MarketData) AlpacaBarsToDbBars(symbol string, bars []marketdata.Bar) []db.Bar {
	fmt.Println("AlpacaBarsToDbBars")
	var dbBars []db.Bar

	for _, bar := range bars {
		dbBars = append(dbBars, m.AlpacaBarToDbBar(symbol, bar))
	}
	return dbBars
}
func (m MarketData) YahooBarToDbBar(symbol string, bar financeGo.ChartBar) db.Bar {
	// fmt.Println("YahooBarToDbBar")
	barToJson, _ := json.MarshalIndent(bar, "", "  ")
	barStructToJson := string(barToJson)
	open, _ := bar.Open.Float64()
	close, _ := bar.Close.Float64()
	high, _ := bar.High.Float64()
	low, _ := bar.Low.Float64()
	dbBar := db.Bar{
		Symbol:          symbol,
		Timestamp:       int64(bar.Timestamp),
		Open:            open,
		Close:           close,
		High:            high,
		Low:             low,
		BarStructToJson: barStructToJson,
	}
	return dbBar
}
func (m MarketData) GetMarketDataFromDb(symbol string, startTime time.Time) []db.Bar {
	fmt.Println("GetMarketDataFromDb")
	var bars []db.Bar
	m.db.Where("symbol = ?", symbol).Where("timestamp >= ?", startTime.Unix()).Order("timestamp asc").Find(&bars)
	// for index, _ := range bars {
	// 	if index == 0 {
	// 		continue
	// 	}
	// 	bar1 := bars[index-1]
	// 	bar2 := bars[index]
	// 	fmt.Println("DB BARS DIFF: \n", bar2.Timestamp-bar1.Timestamp)
	// }
	fmt.Println("DB BARS: ", len(bars), symbol)
	return bars
}
func (m MarketData) FindSymbolLasBar(symbol string) (db.Bar, error) {
	var bar db.Bar
	res := m.db.Where("symbol = ?", symbol).Where("timestamp = (SELECT MAX(timestamp) FROM bars WHERE symbol = ?) ", symbol).Find(&bar)
	if res.Error != nil || res.RowsAffected == 0 {
		fmt.Println("SYMBOL BAR NOT FOUND")
		return bar, errors.New("SYMBOL BAR NOT FOUND")
	}

	return bar, nil
}

func (m MarketData) FillMarketBars(symbol string, startTime time.Time, endTime time.Time) {
	fmt.Println("FillMarketBars")
	lastBar, err := m.FindSymbolLasBar(symbol)

	if err != nil {
		m.GetAlpacaMarketData(symbol, startTime, endTime)
	}

	m.GetYahooFinanceData(symbol,
		time.Unix(lastBar.Timestamp, 0).AddDate(0, 0, -3),
		endTime,
	)
}

var cacheBarsMap = map[string][]db.Bar{}

func (m MarketData) GetMarketCachedData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	key := fmt.Sprintf("%v:%v:%v", symbol, startTime, endTime)
	if val, ok := cacheBarsMap[key]; ok {
		return val
	}
	m.FillMarketBars(symbol, startTime, endTime)

	dbBars := m.GetMarketDataFromDb(symbol, startTime)
	cacheBarsMap[key] = dbBars
	return dbBars
}
func (m MarketData) CutBarsWithHourFrame(bars []db.Bar, hourFrame float64) []db.Bar {
	barsCount := len(bars)
	fmt.Println("CutBarsWithHourFrame", barsCount)
	if barsCount == 0 {
		panic("bARS FOR CUT FRAME NOT FOUND")
	}
	frameTimeStampInHour := int64(float64(60*60) * hourFrame)
	var slicedBars []db.Bar
	var lastBar db.Bar
	for index, bar := range bars {
		if index == 0 || (bar.Timestamp-lastBar.Timestamp) >= frameTimeStampInHour {
			lastBar = bar
			slicedBars = append(slicedBars, bar)
		}
	}
	return slicedBars
}

func (m MarketData) GetMarketCachedDataWithFrame(hourFrame float64, symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("GetMarketCachedDataWithFrame")

	bars := m.GetMarketCachedData(symbol, startTime, endTime)
	frameBars := m.CutBarsWithHourFrame(bars, hourFrame)

	return frameBars
}
