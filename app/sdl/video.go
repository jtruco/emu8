package sdl

import (
	"log"
	"unsafe"

	"github.com/jtruco/emu8/device/video"
	"github.com/veandco/go-sdl2/sdl"
)

// Video is the SDL video UI
type Video struct {
	app        *App         // The SDL app
	window     *sdl.Window  // The main window
	winsurface *sdl.Surface // The window surface
	surface    *sdl.Surface // The emulator surface
	buffer     []sdl.Rect   // The buffer
	scale      int32
	fullscreen bool
}

// NewVideo creates a new video UI
func NewVideo(app *App) *Video {
	video := &Video{}
	video.app = app
	return video
}

// Init initialices video
func (video *Video) Init() bool {
	// Configuration
	video.scale = int32(video.app.config.VideoScale)
	video.fullscreen = video.app.config.FullScreen
	// SDL window & surface
	device := video.app.emulator.Controller().Video().Video()
	screen := device.Screen()
	display := screen.Display()
	window, err := sdl.CreateWindow(
		video.app.config.AppTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		(int32(display.W) * video.scale),
		(int32(display.H) * video.scale),
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println("Error initializing SDL window : " + err.Error())
		return false
	}
	surface, err := window.GetSurface()
	if err != nil {
		log.Println("Error creating window surface : " + err.Error())
		return false
	}
	video.window = window
	video.winsurface = surface
	// Emulator screen surface
	pixels := unsafe.Pointer(&screen.Data()[0])
	surface, err = sdl.CreateRGBSurfaceWithFormatFrom(
		pixels,
		int32(screen.Width()),
		int32(screen.Height()),
		32,
		4*int32(screen.Width()),
		uint32(sdl.PIXELFORMAT_BGRA32))
	if err != nil {
		log.Println("Error creating emulator surface : " + err.Error())
		return false
	}
	surface.SetBlendMode(sdl.BLENDMODE_NONE)
	video.surface = surface
	// Screen regions buffer
	video.buffer = make([]sdl.Rect, len(screen.Regions()))
	return true
}

// Render renders screen to video UI
func (video *Video) Render(screen *video.Screen) {
	var srcRect sdl.Rect
	display := screen.Display()
	rects := screen.DirtyRegions()
	count := len(rects)
	for idx, rect := range rects {
		srcRect = sdl.Rect{
			X: int32(rect.X),
			Y: int32(rect.Y),
			W: int32(rect.W),
			H: int32(rect.H)}
		video.buffer[idx] = sdl.Rect{
			X: int32(rect.X-display.X) * video.scale,
			Y: int32(rect.Y-display.Y) * video.scale,
			W: int32(rect.W) * video.scale,
			H: int32(rect.H) * video.scale}
		video.surface.BlitScaled(&srcRect, video.winsurface, &video.buffer[idx])
	}
	video.window.UpdateSurfaceRects(video.buffer[:count])
}
