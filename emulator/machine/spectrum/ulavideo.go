package spectrum

import (
	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/video"
)

// -----------------------------------------------------------------------------
// ZX Spectrum - ULA video tv output
// -----------------------------------------------------------------------------

// Size constants
const (
	tvPaperWidth   = 256
	tvPaperHeight  = 192
	tvBorderWidth  = 32
	tvBorderHeight = 24
	tvScreenWidth  = tvPaperWidth + 2*tvBorderWidth
	tvScreenHeight = tvPaperHeight + 2*tvBorderHeight
)

// Timing constants
//
// Vertical   : 312 lines = 16 sync, 48 border top, 192 paper, 56 border bottom
// Horizontal : 224 Ts = 24 border left, 128 paper, 24 border right, 48 retrace
// A line starts with its left border, 24 Ts before its first paper pixel, so
// the first paper pixel (0,0) is at Tstate 14336 = 64 lines * 224 Ts.
const (
	tvScanlines        = 312
	tvLineTstates      = 224
	tvFrameTstates     = tvScanlines * tvLineTstates // 69888
	tvLeftTstates      = 24                          // left border Ts before paper
	tvPaperTstates     = tvPaperWidth / 2            // 1 Ts == 2 px
	tvBorderTstates    = tvBorderWidth / 2           // visible border Ts
	tvPaperLine        = 64                          // first paper line
	tvScreenTstate     = tvPaperLine * tvLineTstates // 14336
	tvFirstVisibleLine = tvPaperLine - tvBorderHeight
	tvEndVisibleLine   = tvPaperLine + tvPaperHeight + tvBorderHeight
	tvColVisible       = tvLeftTstates - tvBorderTstates // line Ts : first visible
	tvColPaper         = tvLeftTstates                   // line Ts : first paper
	tvColBorderRight   = tvColPaper + tvPaperTstates     // line Ts : first right border
	tvColEndVisible    = tvColBorderRight + tvBorderTstates
)

// Video memory constants
const (
	tvDataSize  = 0x1800                   // 6 Kbytes (6144)
	tvAttrSize  = 0x0300                   // 768 bytes
	tvVideoSize = tvDataSize + tvAttrSize  // 6912 bytes
	tvVideoAddr = 0x0                      // bank at 0x4000
	tvAttrAddr  = tvVideoAddr + tvDataSize // 0x1800
)

var (
	// Paper line addresses table : bitmap & attribute
	tvPaperAddr [tvPaperHeight][2]int

	// ZX Spetrum 16/48k RGBA colour palette
	zxPaletteRGBA = []uint32{
		/* Bright 0 (black, blue, red, magenta, green, cyan, yellow, white) */
		0xff000000, 0xffc00000, 0xff0000c0, 0xffc000c0, 0xff00c000, 0xffc0c000, 0xff00c0c0, 0xffc0c0c0,
		/* BRIGHT 1 (black, blue, red, magenta, green, cyan, yellow, white) */
		0xff000000, 0xffff0000, 0xff0000ff, 0xffff00ff, 0xff00ff00, 0xffffff00, 0xff00ffff, 0xffffffff,
	}
)

// UlaTv is the spectrum RF video device.
//
// The TV beam is rendered on demand: Update draws from the last rendered
// Tstate up to the current clock Tstate. It must be called before any change
// that affects the output (border, video memory) and at the end of the frame,
// so the beam position is always derived from the machine clock.
type UlaTv struct {
	screen    *video.Screen // The video screen
	clock     device.Clock  // The machine clock
	scrdata   []byte        // The screen data
	border    byte          // The border current colour index
	borderRgb uint32        // The border rgb colour
	flash     bool          // Flash state
	frames    int           // Frame count
	tstate    int           // Frame Tstate rendered so far (beam position)
	scrByte   byte          // Current paper bitmap byte
	attInk    uint32        // Current paper ink colour
	attPaper  uint32        // Current paper paper colour
}

// NewUlaTv creates the ULA TV video device
func NewUlaTv(spectrum *Spectrum) *UlaTv {
	var tv = new(UlaTv)
	tv.clock = spectrum.clock
	tv.scrdata = spectrum.memory.Bank(zxVideoMemory).Data()
	tv.screen = video.NewScreen(tvScreenWidth, tvScreenHeight, zxPaletteRGBA)
	return tv
}

// SetBorder sets de current border color
func (tv *UlaTv) SetBorder(color byte) {
	tv.Update(tv.clock.Tstates())
	tv.border = color
	tv.borderRgb = tv.screen.GetColor(int(color))
}

// Device

// Init initializes video device
func (tv *UlaTv) Init() { tv.Reset() }

// Reset resets video device
func (tv *UlaTv) Reset() {
	tv.screen.Clear(0)
	tv.tstate = 0
	tv.border = 7
	tv.borderRgb = tv.screen.GetColor(7)
	tv.flash = false
}

// Video

// Screen the video screen
func (tv *UlaTv) Screen() *video.Screen { return tv.screen }

// EndFrame renders the rest of the frame and starts a new one
func (tv *UlaTv) EndFrame() {
	tv.Update(tvFrameTstates)
	tv.tstate = 0
	tv.frames++
	tv.flash = (tv.frames & 0x10) == 0
}

// Emulation

// Resync sets the beam position to the clock without rendering (i.e. after
// the clock has been set by a snapshot load)
func (tv *UlaTv) Resync() {
	tv.tstate = tv.clock.Tstates()
	if tv.tstate > tvFrameTstates {
		tv.tstate = tvFrameTstates
	}
}

// Update renders the TV beam up to frame Tstate (exclusive)
func (tv *UlaTv) Update(tstate int) {
	if tstate > tvFrameTstates {
		tstate = tvFrameTstates
	}
	for tv.tstate < tstate {
		pos := tv.tstate + tvLeftTstates
		line, col := pos/tvLineTstates, pos%tvLineTstates
		end := col + tstate - tv.tstate
		if end > tvLineTstates {
			end = tvLineTstates
		}
		if line >= tvFirstVisibleLine && line < tvEndVisibleLine {
			tv.renderLine(line, col, end)
		}
		tv.tstate += end - col
	}
}

// renderLine renders line Tstates [col, end)
func (tv *UlaTv) renderLine(line, col, end int) {
	if col < tvColVisible {
		col = tvColVisible
	}
	if end > tvColEndVisible {
		end = tvColEndVisible
	}
	y := line - tvFirstVisibleLine
	paperY := line - tvPaperLine
	isPaper := paperY >= 0 && paperY < tvPaperHeight
	for ; col < end; col++ {
		x := (col - tvColVisible) << 1
		if isPaper && col >= tvColPaper && col < tvColBorderRight {
			tv.paintPaper(x, y, paperY, col-tvColPaper)
		} else {
			tv.screen.SetPixel(x, y, tv.borderRgb)
			tv.screen.SetPixel(x+1, y, tv.borderRgb)
		}
	}
}

// paintPaper paints the 2 pixels of a paper Tstate (4 Ts per byte)
func (tv *UlaTv) paintPaper(x, y, paperY, paperTs int) {
	if paperTs&0x03 == 0 {
		tv.fetch(paperY, paperTs>>2)
	}
	mask := byte(0x80) >> ((paperTs & 0x03) << 1)
	tv.paintPixel(x, y, tv.scrByte&mask)
	tv.paintPixel(x+1, y, tv.scrByte&(mask>>1))
}

func (tv *UlaTv) paintPixel(x, y int, bit byte) {
	if bit != 0 {
		tv.screen.SetPixel(x, y, tv.attInk)
	} else {
		tv.screen.SetPixel(x, y, tv.attPaper)
	}
}

// fetch reads the bitmap & attribute bytes of a paper column
func (tv *UlaTv) fetch(paperY, column int) {
	addrs := tvPaperAddr[paperY]
	tv.scrByte = tv.scrdata[addrs[0]+column]
	attByte := tv.scrdata[addrs[1]+column]
	ink := attByte & 0x07
	paper := (attByte >> 3) & 0x07
	if (attByte & 0x40) != 0 {
		ink |= 0x08
		paper |= 0x08
	}
	if (attByte&0x80) != 0 && tv.flash {
		ink, paper = paper, ink
	}
	tv.attInk = tv.screen.GetColor(int(ink))
	tv.attPaper = tv.screen.GetColor(int(paper))
}

func init() {
	// Paper line addresses table
	for y := 0; y < tvPaperHeight; y++ {
		tvPaperAddr[y][0] = tvVideoAddr + ((y & 0xc0) << 5) + ((y & 0x38) << 2) + ((y & 0x07) << 8)
		tvPaperAddr[y][1] = tvAttrAddr + ((y & 0xf8) << 2)
	}
}
