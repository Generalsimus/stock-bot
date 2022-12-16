package draw

import (
	"fmt"
	"image/color"
	"neural/algo"
	"neural/market"
	"neural/options"
	"neural/utils"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func DrawOnNewWindow[T int | uint | float64 | float32](draws [][]T, markLine int) {
	const (
		Height unit.Dp = 900
		Width  unit.Dp = 1700
	)
	var Ops op.Ops
	window := app.NewWindow(app.Size(Width, Height))
	ops := &Ops
	red := color.NRGBA{R: 0x80, A: 0xFF}
	for e := range window.Events() {
		switch e := e.(type) {
		// case system.DestroyEvent:
		// 	panic(e.Err)
		case system.FrameEvent:
			// gtx := layout.NewContext(ops, e)
			ops.Reset()
			var path clip.Path
			path.Begin(ops)
			rowCount := len(draws)
			viewWidth := float64(Width)
			viewHeight := float64(Height) / float64(rowCount)
			for rowIndex, row := range draws {
				min, max := utils.FindMinAndMax(row)
				rowCount := len(row)
				if rowIndex == 0 {
					fmt.Println("RowCount:", rowCount)
				}
				for colIndex, colValue := range row {
					y := float32(((viewHeight * float64(rowIndex+1)) - (((float64(colValue) - float64(min)) / (float64(max) - float64(min))) * viewHeight)))
					x := float32((float64(colIndex) / float64(rowCount-1)) * float64(viewWidth))
					if colIndex == (rowCount - markLine) {
						maxY := float32(((viewHeight * float64(rowIndex+1)) - (((float64(max) - float64(min)) / (float64(max) - float64(min))) * viewHeight)))
						minY := float32(((viewHeight * float64(rowIndex+1)) - (((float64(min) - float64(min)) / (float64(max) - float64(min))) * viewHeight)))
						// maxX := float32((float64(colIndex) / float64(rowCount-1)) * float64(viewWidth))
						path.LineTo(f32.Pt(x, y))
						path.LineTo(f32.Pt(x, maxY))
						path.MoveTo(f32.Pt(x, maxY))
						path.LineTo(f32.Pt(x, minY))
						path.MoveTo(f32.Pt(x, minY))
					}
					if colIndex != 0 {
						// fmt.Println("WWW", x, y)
						path.LineTo(f32.Pt(x, y))

					}
					path.MoveTo(f32.Pt(x, y))
				}
			}
			paint.FillShape(ops, red,
				// Stroke
				// Outline
				clip.Stroke{
					Path:  path.End(),
					Width: 2,
				}.Op())

			e.Frame(ops)
		}
	}

}

func ShowMarketMoveData(item *algo.SymbolBestTimeIntervalsBars) {

	drawValue := algo.ConvertToDrawWindow(item.Interval, options.ViewCandles)
	DrawOnNewWindow(drawValue, options.ViewCandles)
}
func FineOpenOrder(item *algo.SymbolBestTimeIntervalsBars, window fyne.Window) {
	symbolObject, priceObject, hourFrameObject, sideObject, stopLostObject, takeProfitObject := GetItemLabels(item)
	marketOrder := market.NewMarket()
	dialog.ShowCustomConfirm(
		"Open Order?",
		"Yes",
		"No",
		container.NewGridWithRows(
			1,
			container.NewGridWithRows(
				3,
				symbolObject,
				hourFrameObject,
				priceObject,
			),
			container.NewGridWithRows(
				3,
				sideObject,
				stopLostObject,
				takeProfitObject,
			),
		),
		func(val bool) {
			if val {
				orderPosition := algo.FindAlpacaRelativeOpenPosition(item)
				if marketOrder.CheckOrderIsExpired(orderPosition) {
					order, err := marketOrder.OrderMarket(orderPosition)
					if err != nil {
						fmt.Printf("failed place order: %v\n", err)
					} else {
						marketOrder.SaveOnDb(orderPosition)
						utils.LogStruct("ORDER SUCCESSFUL: ", order)
					}
				} else {
					fmt.Println("ORDER ALREADY EXiST")
				}

			}
		},
		window,
	)

}
func GetItemLabels(item *algo.SymbolBestTimeIntervalsBars) (fyne.CanvasObject, fyne.CanvasObject, fyne.CanvasObject, fyne.CanvasObject, fyne.CanvasObject, fyne.CanvasObject) {
	symbolObject := widget.NewLabel(fmt.Sprintf("Symbol: %v", item.Symbol))
	priceObject := widget.NewLabel(fmt.Sprintf("Price: %v", item.LastBar.Close))
	hourFrameObject := widget.NewLabel(fmt.Sprintf("Hour Frame: %v", item.HourFrame))
	position := algo.FindAlpacaRelativeOpenPosition(item)
	sideObject := widget.NewLabel(fmt.Sprintf("Side: %v", position.Side))
	stopLostObject := widget.NewLabel(fmt.Sprintf("Stop Lost: %.2f", position.StopLost))
	takeProfitObject := widget.NewLabel(fmt.Sprintf("Take Profit: %.2f", position.TakeProfit))

	return symbolObject, priceObject, hourFrameObject, sideObject, stopLostObject, takeProfitObject
}

func CreateSymbolBox(item *algo.SymbolBestTimeIntervalsBars, window fyne.Window) *fyne.Container {
	symbolObject, priceObject, hourFrameObject, sideObject, stopLostObject, takeProfitObject := GetItemLabels(item)
	//  := GetPriceObject(item.LastBar.Close)
	// space :=
	onOpenOrder := widget.NewButton("Open Order", func() {
		FineOpenOrder(item, window)
	})
	onOpenMarketVisual := widget.NewButton("Show Market", func() {
		go func() {
			ShowMarketMoveData(item)
		}()
	})
	containerBox := container.NewGridWithRows(
		1,
		container.NewGridWithRows(
			3,
			symbolObject,
			hourFrameObject,
			priceObject,
		),
		container.NewGridWithRows(
			3,
			sideObject,
			stopLostObject,
			takeProfitObject,
		),
		layout.NewSpacer(),
		container.NewGridWithRows(
			2,
			onOpenOrder,
			onOpenMarketVisual,
		),
	)
	return containerBox
}
func DrawControllerDashboard(items []*algo.SymbolBestTimeIntervalsBars) {

	var Width float32 = 700
	var Height float32 = 900
	myApp := fyneApp.New()
	window := myApp.NewWindow("TabContainer Widget")
	window.Resize(fyne.NewSize(Width, Height))

	var drawItems []fyne.CanvasObject
	minCount, _ := utils.FindMinAndMax([]int{100, len(items) - 1})
	maxDrawableItems := items[0:minCount]
	for _, item := range maxDrawableItems {
		drawItems = append(drawItems, CreateSymbolBox(item, window), widget.NewSeparator())
	}

	content := container.NewVScroll(
		container.NewVBox(
			drawItems...,
		),
	)
	// content.MinSize
	window.SetContent(content)
	window.ShowAndRun()
}
