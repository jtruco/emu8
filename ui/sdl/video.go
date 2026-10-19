package sdl

import (
	"fmt"
	"log"
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
	window     *sdl.Window   // Main window
	renderer   *sdl.Renderer // Window renderer
	texture    *sdl.Texture  // Emulator screen texture (streaming)
	wRect      sdl.Rect      // Window rect
	sRect      sdl.Rect      // Surface rect
	srcRects   []sdl.Rect    // The source regions cache
	dstRects   []sdl.Rect    // The dest regions cache
	hwAccel    bool          // Hardware accelerated renderer
}

// NewVideo creates a new video UI
func NewVideo() *Video {
	video := new(Video)
	video.Title = "SDL UI"
	video.Scale = 1
	return video
}

// Init initializes video
func (video *Video) Init() error {
	if err := video.sdlCreateWindow(); err != nil {
		return err
	}
	if err := video.sdlCreateTexture(); err != nil {
		return err
	}
	video.createRegions()
	video.updateScreen()
	return nil
}

// Destroy free video resources
func (video *Video) Destroy() {
	log.Println("SDL : Closing video resources")
	if video.texture != nil {
		video.texture.Destroy()
	}
	if video.renderer != nil {
		video.renderer.Destroy()
	}
	if video.window != nil {
		video.window.Destroy()
	}
}

// ToggleFullscreen enable / disable fullscreen mode
func (video *Video) ToggleFullscreen() {
	video.FullScreen = !video.FullScreen
	video.updateScreen()
}

// Update display
func (video *Video) Update(screen *video.Screen) {
	video.OnUpdate(false)
}

// OnUpdate updates screen changes to video display
func (video *Video) OnUpdate(refresh bool) {
	rects := video.Screen.DirtyRects() // none : all the screen
	if refresh {
		rects = nil
	}
	video.updateTexture(rects)
	// copy texture region/s : the accelerated renderer redraws the whole window
	if video.hwAccel || len(rects) == 0 {
		video.renderer.Copy(video.texture, &video.sRect, &video.wRect)
	} else {
		for _, r := range rects {
			video.renderer.Copy(video.texture, &video.srcRects[r], &video.dstRects[r])
		}
	}
	video.renderer.Present()
}

// updateTexture uploads the band of screen rows with dirty regions to the
// screen texture, or the whole screen if there are no regions
func (video *Video) updateTexture(rects []int) {
	width := video.Screen.Width()
	y0, y1 := 0, video.Screen.Height()
	if len(rects) > 0 {
		y0, y1 = y1, 0
		for _, r := range rects {
			rect := &video.srcRects[r]
			if int(rect.Y) < y0 {
				y0 = int(rect.Y)
			}
			if int(rect.Y+rect.H) > y1 {
				y1 = int(rect.Y + rect.H)
			}
		}
	}
	if y1 <= y0 {
		return
	}
	// SDL copies the pixels : no Go memory is kept by C
	pixels := video.Screen.Data()[y0*width:]
	band := sdl.Rect{X: 0, Y: int32(y0), W: int32(width), H: int32(y1 - y0)}
	video.texture.Update(&band, unsafe.Pointer(&pixels[0]), 4*width)
}

func (video *Video) sdlCreateWindow() error {
	var err error
	video.wRect = sdl.Rect{
		X: 0, Y: 0,
		W: int32(float32(video.Screen.View().W) * video.Scale * video.Screen.ScaleX()),
		H: int32(float32(video.Screen.View().H) * video.Scale * video.Screen.ScaleY())}
	video.window, err = sdl.CreateWindow(
		video.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		video.wRect.W, video.wRect.H, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("creating window : %w", err)
	}
	// renderer
	video.renderer, err = sdl.CreateRenderer(video.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil { // no hardware acceleration
		video.renderer, err = sdl.CreateRenderer(video.window, -1, sdl.RENDERER_SOFTWARE)
	}
	if err != nil {
		return fmt.Errorf("creating renderer : %w", err)
	}
	video.renderer.SetLogicalSize(video.wRect.W, video.wRect.H)
	info, _ := video.renderer.GetInfo()
	video.hwAccel = (info.Flags & sdl.RENDERER_ACCELERATED) != 0
	log.Println("SDL : Renderer is:", info.Name)
	return nil
}

func (video *Video) sdlCreateTexture() error {
	var err error
	video.texture, err = video.renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING,
		int32(video.Screen.Width()), int32(video.Screen.Height()))
	if err != nil {
		return fmt.Errorf("creating screen texture : %w", err)
	}
	video.sRect = sdl.Rect{
		X: int32(video.Screen.View().X), Y: int32(video.Screen.View().Y),
		W: int32(video.Screen.View().W), H: int32(video.Screen.View().H)}
	return nil
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

// createRegions creates the screen render regions
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
