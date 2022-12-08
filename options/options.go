package options

import (
	"time"

	"github.com/piquette/finance-go/datetime"
)

var (
	// რამდენი საუკეთესო დამთხვევა აირჩეს
	BestCandles int = 7
	// რამდენი ბარი შეამოწმოს მინიმალური
	StartIntervalCount int = 5
	// რამდენი ბარი შეამოწმოს მაქსიმალური
	EndIntervalCount int = 120
	// ვიზუალიზაცისას დამატებით რამდენი ბარი აჩვენოს ფანჯარაში
	ViewCandles int = 30
	// შესამოწმებელი სიმბოლოები
	CheckSymbols = []string{"GOOGL"}
	// შესამოწმებელი დროის ინტერვალები საათობით
	CheckFrameHours = []float64{2, 4, 8, 16, 24, 48}
	// CheckSymbols = []string{"ABT", "TM", "TXN", "TXN", "AXP", "TD"}
	// ფინანსური სანთლების დროის ინტერვალი
	FinanceInterval = datetime.OneDay
	// ფინანსური სანთლების დროის ინტერვალი
	FinanceIntervals = []datetime.Interval{datetime.OneDay}
	// ფინანსური სანთლების წამოღების საწყისი წერტილი
	FinanceStartDate = datetime.Datetime{Month: 1, Day: 1, Year: 2020}
	// მაქსიმალური დროის ინტერვალი საიდანაც სანთლები შეგვიძლია წამოვიღოთ
	MaxGetBarsStartTime = time.Now().AddDate(-5, 0, 0)
)
