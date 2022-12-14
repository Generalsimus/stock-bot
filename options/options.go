package options

import (
	"time"
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
	// CheckSymbols = []string{"LYTS"}
	CheckSymbols = []string{"LYTS", "SURG", "ERF", "PAYA", "DPSI", "ICVX", "ZVSA", "MOR", "NVAX", "OB", "TDUP", "PRE", "TIG"}
	// შესამოწმებელი დროის ინტერვალები საათობით
	CheckFrameHours = []float64{1}
	// CheckFrameHours = []float64{1, 4, 8, 16, 24}
	// CheckSymbols = []string{"ABT", "TM", "TXN", "TXN", "AXP", "TD"}
	// მაქსიმალური დროის ინტერვალი საიდანაც სანთლები შეგვიძლია წამოვიღოთ
	MaxGetBarsStartTime = time.Now().AddDate(-5, 0, 0).Round(time.Minute)
	// MaxGetBarsStartTime = time.Now().AddDate(0, 0, -2)
)
