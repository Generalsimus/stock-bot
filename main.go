package main

import (
	"fmt"
	"neural/algo"
	"neural/draw"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	SymbolsSimilarity := algo.GetSymbolsSimilarity()
	fmt.Println("BEST Len", len(SymbolsSimilarity))
	draw.DrawControllerDashboard(SymbolsSimilarity)

}
