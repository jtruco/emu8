//go:build sdltest

package sdl

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"testing"
	"unsafe"

	"github.com/jtruco/emu8/config"
	"github.com/jtruco/emu8/emulator/device/video"
	"github.com/veandco/go-sdl2/sdl"
)

// SDL video output tests and benchmarks : see README.md in this package

var sdlInitErr error // SDL initialization error : tests are skipped

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard) // SDL UI logs
	sdlInitErr = sdl.Init(sdl.INIT_VIDEO)
	code := m.Run()
	sdl.Quit()
	os.Exit(code)
}

// geometry is a machine screen geometry
type geometry struct {
	name          string
	width, height int
	scaleX        float32
}

var geometries = []geometry{
	{"zx", 320, 240, 1},    // Spectrum : border + 256x192 paper
	{"cpc", 768, 264, 0.5}, // CPC : border + 640x200 (mode 2 pixels)
}

// newTestVideo returns an initialized SDL video of a new screen. The test
// goroutine is locked to its thread while the video exists: SDL (OpenGL
// contexts, X11) must be used from one thread, and Go moves goroutines.
func newTestVideo(tb testing.TB, g geometry, scale float32) *Video {
	tb.Helper()
	if sdlInitErr != nil {
		tb.Skip("SDL not available:", sdlInitErr)
	}
	runtime.LockOSThread()
	tb.Cleanup(runtime.UnlockOSThread)
	screen := video.NewScreen(g.width, g.height, nil)
	screen.SetScaleX(g.scaleX)
	v := NewVideo()
	v.Title = "emu8 test " + g.name
	v.Scale = scale
	v.Screen = screen
	if err := v.Init(); err != nil {
		tb.Skip("SDL video not available:", err)
	}
	tb.Cleanup(v.Destroy)
	return v
}

// pumpEvents processes the pending window events, as the app event loop
func pumpEvents() {
	for sdl.PollEvent() != nil {
	}
}

// update draws the screen changes as the video controller does
func (v *Video) update(refresh bool) {
	pumpEvents()
	if refresh || v.Screen.IsDirty() {
		v.OnUpdate(refresh)
	}
	v.Screen.SetDirty(false)
}

// fill sets all the screen pixels with a pattern
func fill(screen *video.Screen, seed uint32) {
	for y := 0; y < screen.Height(); y++ {
		for x := 0; x < screen.Width(); x++ {
			screen.SetPixel(x, y, 0xff000000|(uint32(x)*0x0105+uint32(y)*0x010300+seed*0x4f1bbd)&0xffffff)
		}
	}
}

// block sets the pixels of a size x size block
func block(screen *video.Screen, x0, y0, size int, color uint32) {
	for y := y0; y < y0+size && y < screen.Height(); y++ {
		for x := x0; x < x0+size && x < screen.Width(); x++ {
			screen.SetPixel(x, y, color)
		}
	}
}

// renderToTarget draws into an offscreen target of the window size instead
// of the window: window pixels may be undefined when read back (OpenGL pixel
// ownership : not mapped or covered windows)
func (v *Video) renderToTarget(t *testing.T) {
	t.Helper()
	target, err := v.renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_TARGET, v.wRect.W, v.wRect.H)
	if err != nil {
		t.Skip("SDL render target not available:", err)
	}
	t.Cleanup(func() { target.Destroy() })
	if err := v.renderer.SetRenderTarget(target); err != nil {
		t.Skip("SDL render target not available:", err)
	}
	t.Cleanup(func() { v.renderer.SetRenderTarget(nil) })
}

// checkWindow reads the rendered window back and compares it with the
// screen pixels (1:1 scale)
func (v *Video) checkWindow(t *testing.T, step string) {
	t.Helper()
	w, h := int(v.wRect.W), int(v.wRect.H)
	pixels := make([]uint32, w*h)
	if err := v.renderer.ReadPixels(&v.wRect, uint32(sdl.PIXELFORMAT_RGBA32), unsafe.Pointer(&pixels[0]), 4*w); err != nil {
		t.Fatalf("%s: ReadPixels: %v", step, err)
	}
	diffs, first := 0, ""
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			got, want := pixels[y*w+x]&0xffffff, v.Screen.GetPixel(x, y)&0xffffff // RGB
			if got != want {
				if diffs == 0 {
					first = fmt.Sprintf("(%d,%d) = %06x, want %06x", x, y, got, want)
				}
				diffs++
			}
		}
	}
	if diffs > 0 {
		t.Errorf("%s: %d window pixels differ from the screen, first %s", step, diffs, first)
	}
}

// The window shows the screen pixels after partial, full and forced updates
func TestVideoOutputMatchesScreen(t *testing.T) {
	for _, g := range geometries {
		g.scaleX = 1 // 1:1 window pixels
		t.Run(g.name, func(t *testing.T) {
			v := newTestVideo(t, g, 1)
			v.renderToTarget(t)
			screen := v.Screen
			steps := []struct {
				name    string
				draw    func()
				refresh bool
			}{
				{"first frame", func() { fill(screen, 1) }, false},
				{"sprite 1", func() { block(screen, 3, 2, 16, 0xff0000ff) }, false},
				{"sprite 2", func() { block(screen, 11, 7, 16, 0xff00ff00) }, false},
				{"sprite 3", func() { block(screen, 120, 100, 16, 0xffff0000) }, false},
				{"top and bottom rows", func() {
					screen.SetPixel(1, 0, 0xffffffff)
					screen.SetPixel(screen.Width()-1, screen.Height()-1, 0xffffffff)
				}, false},
				{"single pixel", func() { screen.SetPixel(200, 150, 0xff123456) }, false},
				{"no changes", func() {}, false},
				{"full change", func() { fill(screen, 2) }, false},
				{"refresh", func() {}, true},
			}
			for _, s := range steps {
				s.draw()
				if s.refresh || screen.IsDirty() { // else nothing is drawn
					v.render(s.refresh)
				}
				v.checkWindow(t, s.name)
				screen.SetDirty(false)
			}
		})
	}
}

// One screen update (texture upload, copy and present)
func BenchmarkVideoOutput(b *testing.B) {
	for _, g := range geometries {
		b.Run(g.name+"/full", func(b *testing.B) {
			v := newTestVideo(b, g, config.DefaultVideoScale)
			frames := make([][]uint32, 2) // two different screens
			for i := range frames {
				fill(v.Screen, uint32(i))
				frames[i] = append([]uint32(nil), v.Screen.Data()...)
			}
			b.SetBytes(int64(4 * len(frames[0])))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				copy(v.Screen.Data(), frames[i&1])
				v.Screen.SetDirty(true)
				v.update(false)
			}
		})
		b.Run(g.name+"/partial", func(b *testing.B) {
			v := newTestVideo(b, g, config.DefaultVideoScale)
			fill(v.Screen, 0)
			v.update(true)
			const size = 32
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer() // moving sprite
				x := (i * 3) % (g.width - size)
				y := (i * 2) % (g.height - size)
				block(v.Screen, x, y, size, 0xff000000|uint32(i*0x10101))
				b.StartTimer()
				v.update(false)
			}
		})
	}
}
