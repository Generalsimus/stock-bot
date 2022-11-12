package draw

import (
	"fmt"
	"neural/utils"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

var Ops op.Ops

type FrameCallback = func(e system.FrameEvent, ops *op.Ops, w *app.Window)

var callBack FrameCallback = func(e system.FrameEvent, ops *op.Ops, w *app.Window) {

}

// var WidthDp float32 = 500
// var HeightDp float32 = 500
var Width unit.Dp = 1500
var Height unit.Dp = 700

func Init() {
	// w := app.NewWindow(app.Size(Width, Height))
	window := app.NewWindow(app.Size(Width, Height))
	ops := &Ops
	for e := range window.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			panic(e.Err)
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)
			ops.Reset()
			callBack(e, ops, window)
			e.Frame(ops)
		}
	}
	// app.Main()
}

func AddListener(newCallBack FrameCallback) {
	callbackCache := callBack
	callBack = func(e system.FrameEvent, ops *op.Ops, w *app.Window) {
		newCallBack(e, ops, w)
		callbackCache(e, ops, w)
	}
	Ops.Reset()
}

func CalcDrawer(cords []float64) {
	min, max := utils.FindMinAndMax(cords)
	// ops := &Ops
	// ops.Reset()
	// defer clip.Rect{Max: image.Pt(100, 100)}.Push(ops).Pop()
	// paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
	// paint.PaintOp{}.Add(ops)
	// ops := &Ops
	fmt.Println(min + max)

}
