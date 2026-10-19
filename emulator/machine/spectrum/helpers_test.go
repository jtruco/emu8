package spectrum

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
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

// romFileSystem returns a memory file system with a zeroed ROM file
func romFileSystem() *vfs.MemFileSystem {
	fs := vfs.NewMemFileSystem()
	fs.SaveFile(vfs.NewFileInfoData(zxRomName, make([]byte, 0x4000)))
	return fs
}

// newTestSpectrum builds and initializes a 48K Spectrum with a zeroed ROM,
// at frame Tstate 0.
func newTestSpectrum(t *testing.T) *Spectrum {
	t.Helper()
	m := New(ZXSpectrum48K, machine.Settings{}).(*Spectrum)
	controller.New(m, romFileSystem())
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
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
