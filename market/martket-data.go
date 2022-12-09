package market

import (
	"encoding/json"
	"fmt"
	"neural/db"
	"neural/utils"
	"sort"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	financeGo "github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"golang.org/x/exp/slices"
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
	var bars []db.Bar
	for iter.Next() {
		bar := iter.Bar()
		bars = append(bars, m.YahooBarToDbBar(symbol, *bar))
	}
	dbBars := m.OptimizeBars(bars)
	m.SaveBarsOnDb(dbBars)
	return dbBars
}
func (m MarketData) GetAlpacaMarketData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("GetAlpacaMarketData")
	fmt.Println("REQUEST ALPACA BARS: \n", startTime, "\n", endTime)
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
	fmt.Println("GET ALPACA BARS: ", len(quotes))
	dbBars := m.AlpacaBarsToDbBars(symbol, quotes)
	optiBars := m.OptimizeBars(dbBars)
	m.SaveBarsOnDb(optiBars)
	return optiBars
}
func (m MarketData) SaveBarsOnDb(bars []db.Bar) []db.Bar {
	fmt.Println("SaveBarsOnDb", len(bars))
	if len(bars) != 0 {
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
	// barToJson, _ := json.MarshalIndent(bar, "", "  ")
	barStructToJson := ""
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
	var Bars []db.Bar
	m.db.Where("symbol = ?", symbol).Where("timestamp >= ?", startTime.Unix()-2000).Find(&Bars)
	fmt.Println("DB BARS: ", len(Bars))
	return Bars
}

func (m MarketData) OptimizeBars(bars []db.Bar) []db.Bar {
	fmt.Println("OptimizeBars")
	var newBars []db.Bar

	for _, bar := range bars {
		index := slices.IndexFunc(newBars, func(el db.Bar) bool {
			return el.Timestamp == bar.Timestamp
		})
		if index == -1 {
			newBars = append(newBars, bar)
		}

	}
	sort.Slice(newBars, func(index1, index2 int) bool {
		return newBars[index1].Timestamp < newBars[index2].Timestamp
	})

	return newBars
}

func (m MarketData) FillMarketBars(bars []db.Bar, symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("FillMarketBars", len(bars))
	if len(bars) == 0 {
		bars = append(
			bars,
			m.GetAlpacaMarketData(symbol, startTime, endTime)...,
		)
	}

	lastBar := bars[0]

	bars = append(
		bars,
		m.GetYahooFinanceData(symbol,
			time.Unix(lastBar.Timestamp-1000, 0),
			endTime,
		)...,
	)

	return m.OptimizeBars(bars)
}
func (m MarketData) GetMarketCachedData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("GetMarketCachedData")
	barsFromDb := m.OptimizeBars(m.GetMarketDataFromDb(symbol, startTime))

	filedBars := m.FillMarketBars(barsFromDb, symbol, startTime, endTime)

	return filedBars
}
func (m MarketData) CutBarsWithHourFrame(bars []db.Bar, hourFrame float64) []db.Bar {
	fmt.Println("CutBarsWithHourFrame")
	var newBars []db.Bar
	if len(bars) == 0 {
		fmt.Println("bARS FOR CUT FRAME NOT FOUND")
	}
	frameTimeStampInHour := int64(float64(60*60) * hourFrame)
	startTime := time.Unix(bars[0].Timestamp, 0)
	endTime := time.Unix(bars[len(bars)-1].Timestamp, 0)

	for timeStamp := startTime.Unix(); timeStamp < endTime.Unix(); timeStamp += frameTimeStampInHour {
		var closestBar db.Bar
		for index, bar := range bars {
			min1, max1 := utils.FindMinAndMax([]int64{closestBar.Timestamp, timeStamp})
			min2, max2 := utils.FindMinAndMax([]int64{bar.Timestamp, timeStamp})
			if ((max1 - min1) > (max2 - min2)) || index == 0 {
				closestBar = bar
			}
		}

		if len(newBars) == 0 || newBars[len(newBars)-1] != closestBar {
			newBars = append(newBars, closestBar)
		}
	}
	return newBars
}
func (m MarketData) GetMarketCachedDataWithFrame(hourFrame float64, symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("GetMarketCachedDataWithFrame")

	bars := m.GetMarketCachedData(symbol, startTime, endTime)
	frameBars := m.CutBarsWithHourFrame(bars, hourFrame)
	return frameBars
}
