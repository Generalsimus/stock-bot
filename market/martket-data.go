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
	"gorm.io/gorm/clause"
)

type MarketData struct {
	client  marketdata.Client
	options marketdata.ClientOpts
	db      *gorm.DB
}

func (m MarketData) GetMarketData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	fmt.Println("REQUEST ALPACA BARS: \n", startTime, "\n", endTime)
	timeNow := time.Now()
	minute15 := int64(60 * 60 * 15)
	minEnd, _ := utils.FindMinAndMax([]int64{timeNow.Unix() - minute15, endTime.Unix()})
	quotes, err := m.client.GetBars(symbol, marketdata.GetBarsParams{
		TimeFrame:  marketdata.OneMin,
		Start:      startTime,
		End:        time.Unix(minEnd, 0),
		Adjustment: marketdata.Split,
		// PageLimit:  1000,
		// AsOf:      "2022-06-10", // Leaving it empty yields the same results
	})
	if err != nil {
		panic(err)
	}
	for index, _ := range quotes {
		if (index + 1) == len(quotes) {
			break
		}
		bar1 := quotes[index]
		bar2 := quotes[index+1]
		diff := bar2.Timestamp.Unix() - bar1.Timestamp.Unix()
		if (diff) != 60 {
			fmt.Println(
				"\nDIFF:",
				diff,
				"\nDIFF MINUT:",
				diff/60,
				"\nBAR 1 TIME:",
				bar1.Timestamp,
				"\nBAR 2 TIME:",
				bar2.Timestamp,
			)
			panic("NOT RELEVANT BAR")
		}

	}
	dbBars2 := m.AlpacaBarsToDbBars(symbol, quotes)
	dbBars := m.OptimizeBars(dbBars2)

	m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "timestamp"}},
		// DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
	}).Create(&dbBars)
	fmt.Println("DB BARS: ", len(dbBars2))
	// panic("SS")
	return dbBars
}
func (m MarketData) AlpacaBarToDbBar(symbol string, bar marketdata.Bar) db.Bar {
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
	var dbBars []db.Bar
	for _, bar := range bars {
		dbBars = append(dbBars, m.AlpacaBarToDbBar(symbol, bar))
	}
	return dbBars
}

func (m MarketData) GetMarketDataFromDb(symbol string, startTime time.Time) []db.Bar {
	var Bars []db.Bar
	m.db.Where("symbol = ?", symbol).Where("timestamp >= ?", startTime.Unix()-2000).Find(&Bars)
	fmt.Println("DB BARS: ", len(Bars))
	return Bars
}

func (m MarketData) OptimizeBars(bars []db.Bar) []db.Bar {
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
	fmt.Println("FILLABLE_BAR: ", len(bars))
	if len(bars) == 0 {
		m.GetMarketData(symbol, startTime, endTime)
		return m.GetMarketCachedData(symbol, startTime, endTime)
	}
	var newBars []db.Bar
	barsCunt := len(bars)
	for index, _ := range bars {
		if (index + 1) == barsCunt {
			break
		}
		bar1 := bars[index]
		bar2 := bars[index+1]
		timeStampDiff := float64(bar2.Timestamp - bar1.Timestamp)
		if timeStampDiff != 60 {
			startTimeBar := time.Unix(bar1.Timestamp-int64(timeStampDiff), 0)
			endTimeBar := time.Unix(bar2.Timestamp+int64(timeStampDiff), 0)

			m.GetMarketData(symbol, startTimeBar, endTimeBar)
			return m.GetMarketCachedData(symbol, startTime, endTime)
		}
	}
	return m.OptimizeBars(newBars)
}
func (m MarketData) GetMarketCachedData(symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	barsFromDb := m.OptimizeBars(m.GetMarketDataFromDb(symbol, startTime))

	filedBars := m.FillMarketBars(barsFromDb, symbol, startTime, endTime)

	return filedBars
}
func (m MarketData) GetMarketCachedDataWithFrame(hourFrame float64, symbol string, startTime time.Time, endTime time.Time) []db.Bar {
	var newBars []db.Bar
	bars := m.GetMarketCachedData(symbol, startTime, endTime)
	frameTimeStampInHour := int64(float64(60*60) * hourFrame)
	for timeStamp := startTime.Unix(); timeStamp < endTime.Unix(); timeStamp += frameTimeStampInHour {
		var closestBar db.Bar
		for index, bar := range bars {
			min1, max1 := utils.FindMinAndMax([]int64{closestBar.Timestamp, timeStamp})
			min2, max2 := utils.FindMinAndMax([]int64{bar.Timestamp, timeStamp})
			if ((max1 - min1) > (max2 - min2)) || index == 0 {
				closestBar = bar
			}
		}
		newBars = append(newBars, closestBar)
	}
	return newBars
}
