package sdl

import (
	"log"
	"sync"
	"unsafe"

	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/device/video"
	"github.com/veandco/go-sdl2/sdl"
)

// Video is the SDL video UI
type Video struct {
	_sync      sync.Mutex     // Sync object
	device     video.Video    // The machine video device
	config     *config.Config // Configuration
	window     *sdl.Window    // The main window
	winsurface *sdl.Surface   // The window surface
	surface    *sdl.Surface   // The emulator surface
	buffer     []sdl.Rect     // The update surface buffer
	srcRects   []sdl.Rect     // The source regions cache
	dstRects   []sdl.Rect     // The dest regions cache
	fullscreen bool           // Full Screen window mode
	scale      float32        // Video scale configuration
	scaleX     float32        // Render horizontal scale factor
	scaleY     float32        // Render vertical scale factor
}

// NewVideo creates a new video UI
func NewVideo(config *config.Config) *Video {
	video := new(Video)
	video.config = config
	return video
}

// Init initialices video
func (video *Video) Init(device video.Video) bool {
	// configuration
	video.device = device
	video.scale = float32(video.config.Video.Scale)
	video.fullscreen = video.config.Video.FullScreen
	// initialization
	return video.initSDLVideo()
}

// Render renders screen to video UI
func (video *Video) Render(screen *video.Screen) {
	video._sync.Lock()
	defer video._sync.Unlock()

	rects := screen.DirtyRects()
	if len(rects) == 0 {
		idx := len(video.srcRects) - 1
		video.surface.BlitScaled(&video.srcRects[idx], video.winsurface, &video.dstRects[idx])
		video.window.UpdateSurface()
	} else {
		for idx, rect := range rects {
			video.buffer[idx] = video.dstRects[rect]
			video.surface.BlitScaled(&video.srcRects[rect], video.winsurface, &video.dstRects[rect])
		}
		video.window.UpdateSurfaceRects(video.buffer[:len(rects)])
	}
}

func (video *Video) initScreenRects() {
	screen := video.device.Screen()
	viewport := screen.View()
	screen.SetDirty(true)
	// create regions cache
	rects := screen.Rects()
	l := len(rects)
	video.srcRects = make([]sdl.Rect, l+1)
	video.dstRects = make([]sdl.Rect, l+1)
	for i, r := range rects {
		rect := r.Intersect(&viewport) // only in view
		video.srcRects[i] = sdl.Rect{
			X: int32(rect.X), Y: int32(rect.Y),
			W: int32(rect.W), H: int32(rect.H)}
		video.dstRects[i] = sdl.Rect{
			X: int32(float32(rect.X-viewport.X) * video.scaleX),
			Y: int32(float32(rect.Y-viewport.Y) * video.scaleY),
			W: int32(float32(rect.W) * video.scaleX),
			H: int32(float32(rect.H) * video.scaleY)}
	}
	// viewport
	video.srcRects[l] = sdl.Rect{
		X: int32(viewport.X), Y: int32(viewport.Y),
		W: int32(viewport.W), H: int32(viewport.H)}
	video.dstRects[l] = sdl.Rect{
		X: 0, Y: 0,
		W: int32(float32(viewport.W) * video.scaleX),
		H: int32(float32(viewport.H) * video.scaleY)}
}

// ToggleFullscreen enable / disable fullscreen mode
func (video *Video) ToggleFullscreen() {
	video.fullscreen = !video.fullscreen
	video.initSDLVideo()
}

func (video *Video) initSDLVideo() bool {
	video._sync.Lock()
	defer video._sync.Unlock()

	video.destroySDLWindow()
	if !video.createSDLWindow() {
		return false
	}
	if !video.createEmulatorSurface() {
		return false
	}
	video.initScreenRects()
	return true
}

func (video *Video) createSDLWindow() bool {
	screen := video.device.Screen()
	viewport := video.device.Screen().View()
	video.scaleX = video.scale * screen.ScaleX()
	video.scaleY = video.scale * screen.ScaleY()
	window, err := sdl.CreateWindow(
		video.config.AppTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(float32(viewport.W)*video.scaleX),
		int32(float32(viewport.H)*video.scaleY),
		0)
	if err != nil {
		log.Println("Error initializing SDL window : " + err.Error())
	}
	if video.fullscreen {
		window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	}
	window.Show()
	surface, err := window.GetSurface()
	if err != nil {
		log.Println("Error creating window surface : " + err.Error())
		return false
	}
	video.window = window
	video.winsurface = surface
	return true
}

func (video *Video) destroySDLWindow() {
	if video.window != nil {
		video.window.Destroy()
	}
}

func (video *Video) createEmulatorSurface() bool {
	screen := video.device.Screen()
	pixels := unsafe.Pointer(&screen.Data()[0])
	surface, err := sdl.CreateRGBSurfaceWithFormatFrom(
		pixels,
		int32(screen.Width()),
		int32(screen.Height()),
		32,
		4*int32(screen.Width()),
		uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		log.Println("Error creating emulator surface : " + err.Error())
		return false
	}
	surface.SetBlendMode(sdl.BLENDMODE_NONE)
	video.surface = surface
	video.buffer = make([]sdl.Rect, len(screen.Rects()))
	return true
}
