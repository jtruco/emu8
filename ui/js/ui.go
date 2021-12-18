// Package js is the Js/Wasm user interface
package js

import (
	"errors"
	"syscall/js"

	"github.com/jtruco/emu8/emulator/device/video"
)

// JS/HTML user interface

var (
	_window                = js.Global().Get("window")
	_document              = js.Global().Get("document")
	_requestAnimationFrame = js.Global().Get("requestAnimationFrame")
)

// JS function helpers

type JsFunc struct{ js js.Func }

func JsFuncOf(f func()) JsFunc {
	return JsFunc{
		js: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			f()
			return nil
		}),
	}
}
func CreateElement(tagName string) js.Value {
	return _document.Call("createElement", tagName)
}
func RequestAnimationFrame(f JsFunc) {
	_requestAnimationFrame.Invoke(f.js)
}

// JS UI

// JsUi is the JavaScript / HTML emulator UI
type JsUi struct {
	name       string
	container  js.Value
	canvas     Canvas
	onKeyEvent func(int, bool)
	onAction   func(string)
}

func NewJsUi(name string) *JsUi {
	ui := new(JsUi)
	ui.name = name
	ui.canvas.Id = name
	return ui
}

func (ui *JsUi) Init() error {
	// get canvas rendering context
	if !ui.canvas.Bind(ui.name) {
		return errors.New("JS: Could not get canvas rendering context : " + ui.name)
	}
	ui.canvas.Focus()

	// ui action buttons
	for _, name := range emuActions {
		e := _document.Call("getElementById", (ui.name + "_" + name))
		if !e.IsNull() {
			e.Call("addEventListener", "click", js.FuncOf(ui.eventClick))
		}
	}

	// keyboard events
	ui.canvas.Value.Call("addEventListener", "keydown", js.FuncOf(ui.eventKeyDown))
	ui.canvas.Value.Call("addEventListener", "keyup", js.FuncOf(ui.eventKeyUp))
	return nil
}

func (ui *JsUi) eventClick(this js.Value, args []js.Value) interface{} {
	actionId := args[0].Get("target").Get("id").String()
	ui.onAction(actionId[len(ui.name)+1:])
	return nil
}

func (ui *JsUi) eventKeyDown(this js.Value, args []js.Value) interface{} {
	event := args[0]
	event.Call("preventDefault")
	code := jsKeyCodes[event.Get("code").String()]
	ui.onKeyEvent(code, true)
	return nil
}

func (ui *JsUi) eventKeyUp(this js.Value, args []js.Value) interface{} {
	event := args[0]
	code := jsKeyCodes[event.Get("code").String()]
	ui.onKeyEvent(code, false)
	event.Call("preventDefault")
	return nil
}

// JS/HTML Canvas helper

// Canvas is an HTML canvas with a 2d context
type Canvas struct {
	Id      string
	Value   js.Value
	Context js.Value
	Image   js.Value
	Data    js.Value
	Width   int
	Height  int
}

func NewCanvas(name string) *Canvas {
	c := new(Canvas)
	c.Id = name
	return c
}

func (c *Canvas) Create() bool {
	c.Value = _document.Call("createElement", "canvas")
	if c.Value.IsNull() {
		return false
	}
	c.Value.Set("width", c.Width)
	c.Value.Set("height", c.Height)
	return c.getContext2d()
}

func (c *Canvas) Bind(id string) bool {
	c.Value = _document.Call("getElementById", id)
	if c.Value.IsNull() {
		return false
	}
	c.Id = id
	c.Width = c.Value.Get("width").Int()
	c.Height = c.Value.Get("height").Int()
	return c.getContext2d()
}

func (c *Canvas) FullScreen() {
	if _document.Get("fullscreenElement").IsNull() {
		c.Value.Call("requestFullscreen")
	} else {
		_document.Call("exitFullscreen")
		c.Focus()
	}
}

func (c *Canvas) Focus() {
	c.Value.Call("focus")
}

func (c *Canvas) PutImageData(buffer []byte) {
	js.CopyBytesToJS(c.Data, buffer)
	c.Context.Call("putImageData", c.Image, 0, 0)
}

func (c *Canvas) DrawImage(source *Canvas, srect, drect *video.Rect) {
	c.Context.Call("drawImage",
		source.Value,
		srect.X, srect.Y, srect.W, srect.H,
		drect.X, drect.Y, drect.W, drect.H)
}

func (c *Canvas) getContext2d() bool {
	c.Context = c.Value.Call("getContext", "2d")
	c.Image = c.Context.Call("getImageData", 0, 0, c.Width, c.Height)
	c.Data = c.Image.Get("data")
	return !c.Data.IsNull()
}
