package sdl

import (
	"log"
	"sync"
	"unsafe"

	"github.com/jtruco/emu8/device/video"
	"github.com/veandco/go-sdl2/sdl"
)

// Video is the SDL video UI
type Video struct {
	_sync      sync.Mutex   // Sync object
	app        *App         // The SDL app
	device     video.Video  // The machine video device
	window     *sdl.Window  // The main window
	winsurface *sdl.Surface // The window surface
	surface    *sdl.Surface // The emulator surface
	buffer     []sdl.Rect   // The buffer
	scale      int          // Video scale configuration
	fullscreen bool         // Full Screen window mode
	wscale     float32
	hscale     float32
}

// NewVideo creates a new video UI
func NewVideo(app *App) *Video {
	video := new(Video)
	video.app = app
	return video
}

// Init initialices video
func (video *Video) Init() bool {
	// configuration
	video.device = video.app.emulator.Controller().Video().Video()
	video.scale = video.app.config.VideoScale
	video.fullscreen = video.app.config.FullScreen
	// initialization
	return video.initSDLVideo()
}

// Render renders screen to video UI
func (video *Video) Render(screen *video.Screen) {
	var srcRect sdl.Rect
	video._sync.Lock()
	defer video._sync.Unlock()
	display := screen.Display()
	rects := screen.DirtyRegions()
	for idx, rect := range rects {
		srcRect = sdl.Rect{
			X: int32(rect.X),
			Y: int32(rect.Y),
			W: int32(rect.W),
			H: int32(rect.H)}
		video.buffer[idx] = sdl.Rect{
			X: int32(float32(rect.X-display.X) * video.wscale),
			Y: int32(float32(rect.Y-display.Y) * video.hscale),
			W: int32(float32(rect.W) * video.wscale),
			H: int32(float32(rect.H) * video.hscale)}
		video.surface.BlitScaled(&srcRect, video.winsurface, &video.buffer[idx])
	}
	video.window.UpdateSurfaceRects(video.buffer[:len(rects)])
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
	video.device.Screen().SetDirty(true)
	return true
}

func (video *Video) createSDLWindow() bool {
	screen := video.device.Screen()
	display := video.device.Screen().Display()
	video.wscale = float32(video.scale) * screen.WScale()
	video.hscale = float32(video.scale) * screen.HScale()
	window, err := sdl.CreateWindow(
		video.app.config.AppTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(float32(display.W)*video.wscale),
		int32(float32(display.H)*video.hscale),
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
		uint32(sdl.PIXELFORMAT_BGRA32))
	if err != nil {
		log.Println("Error creating emulator surface : " + err.Error())
		return false
	}
	surface.SetBlendMode(sdl.BLENDMODE_NONE)
	video.surface = surface
	video.buffer = make([]sdl.Rect, len(screen.Regions()))
	return true
}
