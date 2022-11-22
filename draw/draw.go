package draw

import (
	"fmt"
	"image/color"
	"neural/utils"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

//"fyne.io/fyne/v2/theme"

// type FrameCallback = func(e system.FrameEvent, ops *op.Ops, w *app.Window)

// var callBack FrameCallback = func(e system.FrameEvent, ops *op.Ops, w *app.Window) {

// }

// var WidthDp float32 = 500
// var HeightDp float32 = 500

// var Width float32 = 1500
// var Height float32 = 700
// var drawXY []float64 = []float64{}

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
		case system.DestroyEvent:
			panic(e.Err)
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

func FyneDashBoard() {
	var Width float32 = 500
	var Height float32 = 500
	myApp := fyneApp.New()
	window := myApp.NewWindow("TabContainer Widget")
	window.Resize(fyne.NewSize(Width, Height))

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	// button := widget.NewButton("Save", func() {

	// })
	// tabs := container.NewAppTabs()
	// for _, symbol := range finance.DefaultSymbols {
	// 	tabs.Append(container.NewTabItem(symbol, widget.NewLabel("Hello")))
	// }

	// tabs.SetTabLocation(container.TabLocationLeading)

	// content := container.NewVBox(
	// 	container.NewGridWithColumns(
	// 		4,
	// 		widget.NewSelect(finance.DefaultSymbols, func(symbol string) {
	// 			log.Println("AAAA:", symbol)
	// 		}),
	// 		layout.NewSpacer(),
	// 		layout.NewSpacer(),
	// 		container.NewGridWithRows(
	// 			1,
	// 			input,
	// 			button,
	// 		),
	// 	),
	// 	container.NewHBox(tabs),
	// )

	// window.SetContent(content)
	window.ShowAndRun()
}

// func AddListener(newCallBack FrameCallback) {
// 	callbackCache := callBack
// 	callBack = func(e system.FrameEvent, ops *op.Ops, w *app.Window) {
// 		newCallBack(e, ops, w)
// 		callbackCache(e, ops, w)
// 	}
// 	Ops.Reset()
// }

// func CalcDrawer(cords []float64) {
// 	min, max := utils.FindMinAndMax(cords)
// 	// ops := &Ops
// 	// ops.Reset()
// 	// defer clip.Rect{Max: image.Pt(100, 100)}.Push(ops).Pop()
// 	// paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
// 	// paint.PaintOp{}.Add(ops)
// 	// ops := &Ops
// 	fmt.Println(min + max)

// }
