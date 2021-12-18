package js

import (
	"reflect"
	"unsafe"

	"github.com/jtruco/emu8/emulator/device/video"
)

const bufferSuffix = "_buffer"

// Video is the video UI controller
type Video struct {
	uiCanvas     *Canvas
	bufferCanvas *Canvas
	screen       *video.Screen
	screenBytes  []byte
	view, drect  video.Rect
	jsUpdate     JsFunc
}

func NewVideo(canvas *Canvas) *Video {
	v := new(Video)
	v.uiCanvas = canvas
	v.bufferCanvas = NewCanvas(canvas.Id + bufferSuffix)
	v.jsUpdate = JsFuncOf(v.onUpdate)
	return v
}

func (v *Video) Init() bool {
	// update canvas size
	v.onCanvasResize()
	// build screen byte array (for CopyBytesToJS)
	v.screenBytes = toByteSlice(v.screen.Data())
	// create a screen back buffer canvas
	v.bufferCanvas.Width = v.screen.Width()
	v.bufferCanvas.Height = v.screen.Height()
	return v.bufferCanvas.Create()
}

func (v *Video) Update(screen *video.Screen) {
	RequestAnimationFrame(v.jsUpdate) // ui sync
}

func (v *Video) onUpdate() {
	v.bufferCanvas.PutImageData(v.screenBytes)
	v.uiCanvas.DrawImage(v.bufferCanvas, &v.view, &v.drect)
}

func (v *Video) onCanvasResize() {
	// scale factor & destination rect
	v.view = v.screen.View()
	factorW := float32(v.uiCanvas.Width) / (float32(v.view.W) * v.screen.ScaleX())
	factorH := float32(v.uiCanvas.Height) / (float32(v.view.H) * v.screen.ScaleY())
	scale := factorW
	if factorW > factorH {
		scale = factorH
	}
	v.drect.W = int(scale * float32(v.view.W) * v.screen.ScaleX())
	v.drect.H = int(scale * float32(v.view.H) * v.screen.ScaleY())
	if v.drect.W < (v.uiCanvas.Width - 1) {
		v.drect.X = (v.uiCanvas.Width - v.drect.W) / 2
	} else {
		v.drect.Y = (v.uiCanvas.Height - v.drect.H) / 2
	}
}

// toByteSlice converts to []byte
func toByteSlice(data []uint32) []byte {
	ph := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	bh := reflect.SliceHeader{ph.Data, ph.Len << 2, ph.Cap << 2}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
