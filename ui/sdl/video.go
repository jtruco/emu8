package sdl

import (
	"log"
	"sync"
	"unsafe"

	"github.com/jtruco/emu8/emulator/device/video"
	"github.com/veandco/go-sdl2/sdl"
)

// Video is the SDL video UI
type Video struct {
	Screen     *video.Screen // Machine video screen
	Title      string        // Window title
	FullScreen bool          // Fullscreen state
	Scale      float32       // Video scale
	Async      bool          // Asynchronous mode
	UpdateUi   chan bool     // Update UI channel
	window     *sdl.Window   // Main window
	renderer   *sdl.Renderer // Window renderer
	surface    *sdl.Surface  // Emulator screen surface
	wRect      sdl.Rect      // Window rect
	sRect      sdl.Rect      // Surface rect
	srcRects   []sdl.Rect    // The source regions cache
	dstRects   []sdl.Rect    // The dest regions cache
	hwAccel    bool          // Hardware accelerated renderer
	_sync      sync.Mutex    // Sync object
}

// NewVideo creates a new video UI
func NewVideo() *Video {
	video := new(Video)
	video.Title = "SDL UI"
	video.Scale = 1
	video.UpdateUi = make(chan bool, 2) // 2-slot channel
	return video
}

// Init initialices video
func (video *Video) Init() bool {
	video._sync.Lock()
	defer video._sync.Unlock()

	if !video.sdlCreateWindow() {
		return false
	}
	if !video.sdlCreateSurface() {
		return false
	}
	video.createRegions()
	video.updateScreen()
	return true
}

// Destroy free video resources
func (video *Video) Destroy() {
	video._sync.Lock()
	defer video._sync.Unlock()

	log.Println("SDL : Closing video resources")
	video.surface.Free()
	video.renderer.Destroy()
	video.window.Destroy()
}

// ToggleFullscreen enable / disable fullscreen mode
func (video *Video) ToggleFullscreen() {
	video._sync.Lock()
	defer video._sync.Unlock()

	video.FullScreen = !video.FullScreen
	video.updateScreen()
}

// Update display
func (video *Video) Update(screen *video.Screen) {
	if video.Async {
		video.UpdateUi <- true // Notify SDL main thread
		return
	}
	video.OnUpdate(false)
}

// OnUpdate updates screen changes to video display
func (video *Video) OnUpdate(refresh bool) {
	video._sync.Lock()
	defer video._sync.Unlock()

	// create texture from screen surface
	texture, err := video.renderer.CreateTextureFromSurface(video.surface)
	if err != nil {
		return
	}
	defer texture.Destroy()

	// copy texture region/s
	rects := video.Screen.DirtyRects()
	if refresh || video.hwAccel || len(rects) == 0 {
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
	video.wRect = sdl.Rect{
		X: 0, Y: 0,
		W: int32(float32(video.Screen.View().W) * video.Scale * video.Screen.ScaleX()),
		H: int32(float32(video.Screen.View().H) * video.Scale * video.Screen.ScaleY())}
	video.window, err = sdl.CreateWindow(
		video.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		video.wRect.W, video.wRect.H, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println("SDL : Error creating Window:", err.Error())
	}
	// renderer
	video.renderer, err = sdl.CreateRenderer(video.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println("SDL : Errror creating Renderer:", err.Error())
	}
	video.renderer.SetLogicalSize(video.wRect.W, video.wRect.H)
	info, _ := video.renderer.GetInfo()
	video.hwAccel = (info.Flags & sdl.RENDERER_ACCELERATED) != 0
	log.Println("SDL : Renderer is:", info.Name)
	return true
}

func (video *Video) sdlCreateSurface() bool {
	var err error
	pixels := unsafe.Pointer(&video.Screen.Data()[0])
	video.surface, err = sdl.CreateRGBSurfaceWithFormatFrom(
		pixels, int32(video.Screen.Width()), int32(video.Screen.Height()),
		32, 4*int32(video.Screen.Width()), uint32(sdl.PIXELFORMAT_RGBA32))
	video.sRect = sdl.Rect{
		X: int32(video.Screen.View().X), Y: int32(video.Screen.View().Y),
		W: int32(video.Screen.View().W), H: int32(video.Screen.View().H)}
	if err != nil {
		log.Println("Error creating emulator surface:", err.Error())
		return false
	}
	return true
}

// updateScreen update screen state
func (video *Video) updateScreen() {
	// check fullscreen mode
	if video.FullScreen {
		video.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
		sdl.ShowCursor(sdl.DISABLE)
		log.Println("SDL : Screen state is: fullscreen")
	} else {
		video.window.SetFullscreen(0)
		sdl.ShowCursor(sdl.ENABLE)
		log.Println("SDL : Screen state is: windowed")
	}
	video.renderer.Clear()
	// force screen refresh
	video.Screen.SetDirty(true)
}

// render regions
func (video *Video) createRegions() {
	view := video.Screen.View()
	scaleX := video.Scale * video.Screen.ScaleX()
	scaleY := video.Scale * video.Screen.ScaleY()
	regions := video.Screen.Rects()
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
