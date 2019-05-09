package spectrum

import (
	"github.com/jtruco/emu8/device/memory"
	"github.com/jtruco/emu8/device/video"
)

// -----------------------------------------------------------------------------
// Video constants & vars
// -----------------------------------------------------------------------------

// ZX Spectrum video constans
const (
	paperWidth    = 256
	paperHeight   = 192
	borderLeft    = 48
	borderRight   = 48
	borderTop     = 64
	borderBottom  = 56
	screenWidth   = paperWidth + borderLeft + borderRight
	screenHeight  = paperHeight + borderTop + borderBottom
	displayLeft   = 16
	displayTop    = 40
	displayWidth  = paperWidth + 2*(borderLeft-displayLeft)
	displayHeight = paperHeight + 2*(borderTop-displayTop)
	dataSize      = 0x1800              // 6 Kbytes
	attrSize      = 0x0300              // 768 bytes
	videoSize     = dataSize + attrSize // 6912 bytes
	videoAddr     = 0x0                 // bank at 0x4000
	attrAddr      = videoAddr + dataSize
)

// ZX Spetrum 16/48k RGB colour palette
var zxPalette = []int32{
	/* Bright 0 (black, blue, red, magenta, green, cyan, yellow, white) */
	0x000000, 0x0000c0, 0xc00000, 0xc000c0, 0x00c000, 0x00c0c0, 0xc0c000, 0xc0c0c0,
	/* BRIGHT 1 (black, blue, red, magenta, green, cyan, yellow, white) */
	0x000000, 0x0000ff, 0xff0000, 0xff00ff, 0x00ff00, 0x00ffff, 0xffff00, 0xffffff,
}

// -----------------------------------------------------------------------------
// ZX Spectrum TVVideo
// -----------------------------------------------------------------------------

// tvEvent
type tvEvent struct {
	TState int  // TState
	Border byte // Border colour index
}

// TVVideo is the spectrum RF video device
type TVVideo struct {
	screen  *video.Screen // The video screen
	bank    *memory.Bank  // The spectrum video memory bank
	palette []int32       // The video palette
	events  []tvEvent     // The video event queue
	border  byte          // The border current colour index
	flash   bool          // Flash state
	frames  int           // Frame count
}

// NewTVVideo creates the video device
func NewTVVideo(bank *memory.Bank) *TVVideo {
	tv := &TVVideo{}
	tv.palette = zxPalette
	tv.screen = video.NewScreen(screenWidth, screenHeight, tv.palette)
	tv.screen.SetDisplay(video.Rect{X: displayLeft, Y: displayTop, W: displayWidth, H: displayHeight})
	tv.bank = bank
	tv.events = make([]tvEvent, 0, 10)
	return tv
}

// SetBorder sets de current border color
func (tv *TVVideo) SetBorder(tstate int, colour byte) {
	tv.border = colour
	tv.events = append(tv.events, tvEvent{tstate, colour})
}

// Device

// Init initializes video device
func (tv *TVVideo) Init() { tv.Reset() }

// Reset resets video device
func (tv *TVVideo) Reset() {
	tv.screen.Clear(0)
	tv.border = 7
	tv.flash = false
}

// Video

// EndFrame updates screen video frame
func (tv *TVVideo) EndFrame() {
	tv.frames++
	tv.flash = (tv.frames & 0x10) == 0
	tv.paintScreen()
	tv.paintBorder()
	// tv.emulateBorder()
	tv.events = tv.events[:0]
}

// Screen the video screen
func (tv *TVVideo) Screen() *video.Screen { return tv.screen }

// Screen: accurate emulation

func (tv *TVVideo) emulateBorder() {
	// TODO
}

// Screen: simple and fast emulation

// paintScreen is a simple screen emulation
func (tv *TVVideo) paintScreen() {
	videodata := tv.bank.Data()
	flash := tv.flash

	// 3 banks, of 8 rows, of 8 lines, of 32 cols
	baddr := 0
	y := 0
	for b := 0; b < 3; b++ {
		for r := 0; r < 8; r++ {
			laddr := baddr + r*32
			for l := 0; l < 8; l++ {
				caddr := uint16(videoAddr + laddr)
				aaddr := uint16(attrAddr + (y/8)*32)
				for c := 0; c < 32; c++ {
					data := videodata[caddr]
					attr := videodata[aaddr]
					tv.paintByte(y, c, data, attr, flash)
					caddr++
					aaddr++
				}
				y++
				laddr += 0x100 // 8 * 32;
			}
		}
		baddr += 0x800 // 2Kbytes
	}

}

// paintBorder is a simple border emulation
func (tv *TVVideo) paintBorder() {
	// Border Top and Bottom and Paper
	border := tv.palette[tv.border]
	for y := 0; y < borderTop; y++ {
		for x := 0; x < screenWidth; x++ {
			tv.screen.SetPixel(x, y, border)
		}
	}
	for y := borderTop + paperHeight; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			tv.screen.SetPixel(x, y, border)
		}
	}
	for y := borderTop; y < borderTop+paperHeight; y++ {
		for x := 0; x < borderLeft; x++ {
			tv.screen.SetPixel(x, y, border)
		}
		for x := borderLeft + paperWidth; x < screenWidth; x++ {
			tv.screen.SetPixel(x, y, border)
		}
	}
}

// paintByte paints a byte
func (tv *TVVideo) paintByte(y, c int, data, attr byte, flash bool) {
	var ink, paper, mask byte
	ink = attr & 0x07
	if (attr & 0x40) != 0 {
		ink |= 0x08
	}
	paper = (attr >> 3) & 0x07
	if (attr & 0x40) != 0 {
		paper |= 0x08
	}
	mask = 0x80
	if attr&0x80 != 0 && flash {
		ink, paper = paper, ink
	}
	sx := borderLeft + (c << 3)
	y += borderTop
	for x := sx; x < sx+8; x++ {
		set := (data & mask) != 0
		if set {
			tv.screen.SetPixelIndex(x, y, int(ink))
		} else {
			tv.screen.SetPixelIndex(x, y, int(paper))
		}
		mask >>= 1
	}
}
