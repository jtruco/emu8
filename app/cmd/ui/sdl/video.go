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
	_sync    sync.Mutex          // Sync object
	config   *config.VideoConfig // Video configuration
	device   video.Video         // Machine video device
	window   *sdl.Window         // Main window
	renderer *sdl.Renderer       // Window renderer
	surface  *sdl.Surface        // Emulator screen surface
	wRect    sdl.Rect            // Window rect
	sRect    sdl.Rect            // Surface rect
	srcRects []sdl.Rect          // The source regions cache
	dstRects []sdl.Rect          // The dest regions cache
	hwAccel  bool                // Hardware accelerated renderer
}

// NewVideo creates a new video UI
func NewVideo(config *config.Config) *Video {
	video := new(Video)
	video.config = &config.Video
	return video
}

// Init initialices video
func (video *Video) Init(device video.Video) bool {
	video._sync.Lock()
	defer video._sync.Unlock()

	video.device = device
	if !video.sdlCreateWindow() {
		return false
	}
	if !video.sdlCreateSurface() {
		return false
	}
	video.createRegions()
	return true
}

// Destroy free video resources
func (video *Video) Destroy() {
	video._sync.Lock()
	defer video._sync.Unlock()

	video.surface.Free()
	video.renderer.Destroy()
	video.window.Destroy()
}

// ToggleFullscreen enable / disable fullscreen mode
func (video *Video) ToggleFullscreen() {
	video._sync.Lock()
	defer video._sync.Unlock()

	video.config.FullScreen = !video.config.FullScreen
	// FIXME implement fullscreen
}

// Update updates screen changes to video display
func (video *Video) Update(screen *video.Screen) {
	video._sync.Lock()
	defer video._sync.Unlock()

	// create texture from screen surface
	texture, err := video.renderer.CreateTextureFromSurface(video.surface)
	if err != nil {
		return
	}
	defer texture.Destroy()

	// copy texture region/s
	rects := screen.DirtyRects()
	if video.hwAccel || len(rects) == 0 {
		video.renderer.Copy(texture, &video.sRect, &video.wRect)
	} else {
		for _, r := range rects {
			video.renderer.Copy(texture, &video.srcRects[r], &video.dstRects[r])
		}
	}

	// update changes
	video.renderer.Present()
}

func (video *Video) sdlCreateWindow() bool {
	var err error
	screen := video.device.Screen()
	video.wRect = sdl.Rect{
		X: 0, Y: 0,
		W: int32(float32(screen.View().W) * float32(video.config.Scale) * screen.ScaleX()),
		H: int32(float32(screen.View().H) * float32(video.config.Scale) * screen.ScaleY())}
	video.window, err = sdl.CreateWindow(
		config.DefaultAppTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		video.wRect.W, video.wRect.H,
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println("Error initializing SDL window : " + err.Error())
	}
	// renderer
	video.renderer, err = sdl.CreateRenderer(video.window, -1, 0)
	if err != nil {
		log.Println("Errror initializing window Renderer : " + err.Error())
	}
	info, _ := video.renderer.GetInfo()
	video.hwAccel = info.Flags&sdl.RENDERER_ACCELERATED != 0
	return true
}

func (video *Video) sdlCreateSurface() bool {
	var err error
	screen := video.device.Screen()
	pixels := unsafe.Pointer(&screen.Data()[0])
	video.surface, err = sdl.CreateRGBSurfaceWithFormatFrom(
		pixels, int32(screen.Width()), int32(screen.Height()),
		32, 4*int32(screen.Width()), uint32(sdl.PIXELFORMAT_RGBA32))
	video.sRect = sdl.Rect{
		X: int32(screen.View().X), Y: int32(screen.View().Y),
		W: int32(screen.View().W), H: int32(screen.View().H)}
	if err != nil {
		log.Println("Error creating emulator surface : " + err.Error())
		return false
	}
	return true
}

// render regions
func (video *Video) createRegions() {
	screen := video.device.Screen()
	view := screen.View()
	scaleX := float32(video.config.Scale) * screen.ScaleX()
	scaleY := float32(video.config.Scale) * screen.ScaleY()
	regions := screen.Rects()
	video.srcRects = make([]sdl.Rect, len(regions))
	video.dstRects = make([]sdl.Rect, len(regions))
	for i, rect := range regions {
		rect = rect.Intersect(&view) // clip to view
		video.srcRects[i] = sdl.Rect{
			X: int32(rect.X), Y: int32(rect.Y),
			W: int32(rect.W), H: int32(rect.H)}
		video.dstRects[i] = sdl.Rect{
			X: int32(float32(rect.X-view.X) * scaleX),
			Y: int32(float32(rect.Y-view.Y) * scaleY),
			W: int32(float32(rect.W) * scaleX),
			H: int32(float32(rect.H) * scaleY)}
	}
}
