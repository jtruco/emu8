package spectrum

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
)

// Test palette colours (bright 0 / bright 1)
var (
	colBlack   = zxPaletteRGBA[0]
	colBlue    = zxPaletteRGBA[1]
	colRed     = zxPaletteRGBA[2]
	colCyan    = zxPaletteRGBA[5]
	colWhite   = zxPaletteRGBA[7]
	colBrRed   = zxPaletteRGBA[8+2]
	colBrGreen = zxPaletteRGBA[8+4]
)

// newTestSpectrum builds and initializes a 48K Spectrum without ROM files
// (the ROM bank stays zeroed), at frame Tstate 0.
func newTestSpectrum(t *testing.T) *Spectrum {
	t.Helper()
	vfs.SetFileSystem(vfs.NewMemFileSystem())
	m := New(ZXSpectrum48K).(*Spectrum)
	controller.New(m)
	m.Init()
	return m
}

// at sets the machine clock to a frame Tstate
func (m *Spectrum) at(tstate int) *Spectrum {
	m.clock.SetTstates(tstate)
	return m
}

// pixel returns the visible screen pixel at (x, y)
func (m *Spectrum) pixel(x, y int) uint32 {
	return m.tv.Screen().GetPixel(x, y)
}

// videoPoke writes video memory (offset from 0x4000) without bus events
func (m *Spectrum) videoPoke(offset int, value byte) {
	m.tv.scrdata[offset] = value
}

// firstPixel returns the first visible pixel, in raster order, with colour
func (m *Spectrum) firstPixel(colour uint32) (int, int, bool) {
	for y := 0; y < tvScreenHeight; y++ {
		for x := 0; x < tvScreenWidth; x++ {
			if m.pixel(x, y) == colour {
				return x, y, true
			}
		}
	}
	return 0, 0, false
}
