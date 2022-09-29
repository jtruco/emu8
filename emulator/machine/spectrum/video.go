package spectrum

import (
	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/bus"
	"github.com/jtruco/emu8/emulator/device/video"
)

// -----------------------------------------------------------------------------
// Video constants & vars
// -----------------------------------------------------------------------------

// Video screen constants
const (
	tvScreenWidth  = 256
	tvScreenHeight = 192
	tvBorderLeft   = 48
	tvBorderRight  = 48
	tvBorderTop    = 64
	tvBorderBottom = 56
	tvTotalWidth   = tvScreenWidth + tvBorderLeft + tvBorderRight
	tvTotalHeight  = tvScreenHeight + tvBorderTop + tvBorderBottom
	tvViewLeft     = 16 // Border : 32
	tvViewTop      = 40 // Border : 24
	tvViewWidth    = tvScreenWidth + 2*(tvBorderLeft-tvViewLeft)
	tvViewHeight   = tvScreenHeight + 2*(tvBorderTop-tvViewTop)
)

// Video memory constants
const (
	tvDataSize  = 0x1800                   // 6 Kbytes (6144)
	tvAttrSize  = 0x0300                   // 768 bytes
	tvVideoSize = tvDataSize + tvAttrSize  // 6912 bytes
	tvVideoAddr = 0x0                      // bank at 0x4000
	tvAttrAddr  = tvVideoAddr + tvDataSize // 0x1800
)

// TV video timings constants
const (
	tvTstatePixels      = 2
	tvLineTstates       = 224
	tvHBorderTstates    = 24
	tvFirstScreenTstate = 14336
	tvFirstScreenLine   = tvBorderTop
	tvLastScreenLine    = tvBorderTop + tvScreenHeight - 1
)

// ZX Spetrum 16/48k RGBA colour palette
var zxPaletteRGBA = []uint32{
	/* Bright 0 (black, blue, red, magenta, green, cyan, yellow, white) */
	0xff000000, 0xffc00000, 0xff0000c0, 0xffc000c0, 0xff00c000, 0xffc0c000, 0xff00c0c0, 0xffc0c0c0,
	/* BRIGHT 1 (black, blue, red, magenta, green, cyan, yellow, white) */
	0xff000000, 0xffff0000, 0xff0000ff, 0xffff00ff, 0xff00ff00, 0xffffff00, 0xff00ffff, 0xffffffff,
}

// -----------------------------------------------------------------------------
// ZX Spectrum TV video output
// -----------------------------------------------------------------------------

// TvVideo is the spectrum RF video device
type TvVideo struct {
	screen   *video.Screen // The video screen
	clock    device.Clock  // The system clock
	srcdata  []byte        // The screen data
	tstate   int           // Current videoframe tstate
	border   byte          // The border current colour index
	flash    bool          // Flash state
	frames   int           // Frame count
	accurate bool          // Accurate scanlines simulation
}

// NewTVVideo creates the video device
func NewTVVideo(spectrum *Spectrum) *TvVideo {
	tv := new(TvVideo)
	tv.screen = video.NewScreen(tvTotalWidth, tvTotalHeight, zxPaletteRGBA)
	tv.screen.SetView(tvViewLeft, tvViewTop, tvViewWidth, tvViewHeight)
	tv.clock = spectrum.clock
	tv.srcdata = spectrum.memory.Bank(zxVideoMemory).Data()
	spectrum.memory.Map(zxVideoMemory).OnPostAccess = tv.onVideoPostAccess
	tv.accurate = true
	return tv
}

// onVideoPostAccess on write in video memory
func (tv *TvVideo) onVideoPostAccess(code int, address uint16) {
	if code == bus.EventAfterWrite {
		if tv.accurate && address < tvVideoSize {
			tv.DoScanlines()
		}
	}
}

// SetAccurate sets de video emulation algorithm
func (tv *TvVideo) SetAccurate(accurate bool) {
	tv.accurate = accurate
}

// SetBorder sets de current border color
func (tv *TvVideo) SetBorder(colour byte) {
	if tv.accurate {
		tv.DoScanlines()
	}
	tv.border = colour
}

// Device

// Init initializes video device
func (tv *TvVideo) Init() { tv.Reset() }

// Reset resets video device
func (tv *TvVideo) Reset() {
	tv.screen.Clear(0)
	tv.border = 7
	tv.flash = false
}

// Video

// EndFrame updates screen video frame
func (tv *TvVideo) EndFrame() {
	if tv.accurate {
		tv.DoScanlines()
		tv.tstate = 0
	} else {
		tv.paintScreen()
		tv.paintBorder()
	}
	tv.frames++
	tv.flash = (tv.frames & 0x10) == 0
}

// Screen the video screen
func (tv *TvVideo) Screen() *video.Screen { return tv.screen }

// Screen: simple and fast emulation

// paintScreen is a simple screen emulation
func (tv *TvVideo) paintScreen() {
	// 3 banks, of 8 rows, of 8 lines, of 32 cols
	baddr := 0
	y, sy := 0, tvBorderTop
	for b := 0; b < 3; b++ { // banks
		for r := 0; r < 8; r++ { // rows
			laddr := baddr + r<<5
			for l := 0; l < 8; l++ { // lines
				caddr := tvVideoAddr + laddr
				aaddr := tvAttrAddr + (y>>3)<<5
				sx := tvBorderLeft
				for c := 0; c < 32; c++ { // columns
					data := tv.srcdata[caddr]
					attr := tv.srcdata[aaddr]
					tv.paintByte(sy, sx, data, attr)
					caddr++
					aaddr++
					sx += 8
				}
				y++
				sy++
				laddr += 0x100 // 8 * 32;
			}
		}
		baddr += 0x800 // 2Kbytes
	}
}

// paintBorder is a simple border emulation
func (tv *TvVideo) paintBorder() {
	// Border Top, Bottom and Paper
	border := tv.screen.GetColor(int(tv.border))
	view := tv.screen.View()
	for y := view.Y; y < tvBorderTop; y++ {
		tv.scanlineBorder(y, 0, tvTotalWidth-1, border)
	}
	for y := tvBorderTop + tvScreenHeight; y < view.Y+view.H; y++ {
		tv.scanlineBorder(y, 0, tvTotalWidth-1, border)
	}
	for y := tvBorderTop; y < view.Y+view.H; y++ {
		tv.scanlineBorder(y, 0, tvBorderLeft-1, border)
		tv.scanlineBorder(y, tvBorderLeft+tvScreenWidth, tvTotalWidth-1, border)
	}
}

// paintByte paints a byte
func (tv *TvVideo) paintByte(y, sx int, data, attr byte) {
	var ink, paper, mask byte
	ink = attr & 0x07
	if (attr & 0x40) != 0 {
		ink |= 0x08
	}
	paper = (attr >> 3) & 0x07
	if (attr & 0x40) != 0 {
		paper |= 0x08
	}
	if (attr&0x80) != 0 && tv.flash {
		ink, paper = paper, ink
	}
	inkRgb := tv.screen.GetColor(int(ink))
	paperRgb := tv.screen.GetColor(int(paper))
	mask = 0x80
	for x := sx; x < sx+8; x++ {
		set := (data & mask) != 0
		if set {
			tv.screen.SetPixel(x, y, inkRgb)
		} else {
			tv.screen.SetPixel(x, y, paperRgb)
		}
		mask >>= 1
	}
}

// Screen : accurate emulation

// DoScanlines refresh TV scanlines
func (tv *TvVideo) DoScanlines() {
	// Spectrum 48k - Tv Scanlines timings
	// Vertical   : 16 Sl sync, 48 Sl border top, 192 Sl Screen, 56 Sl boder bottom
	// Horizontal : 128 Ts screen, 24 Ts border right, 48 Ts retrace, 24 TS border left
	// First screen (0,0) pixel Tstate = 14336 TS = 64 Scanlines * 224 Tstates
	view := tv.screen.View()
	border := tv.screen.GetColor(int(tv.border))
	tstate := tv.tstate
	endtstate := tv.clock.Tstates()
	limitBottom := view.Y*tvLineTstates - tvHBorderTstates
	limitTop := (view.Y+view.H)*tvLineTstates - tvHBorderTstates
	if endtstate < limitBottom || tstate > limitTop {
		return
	}
	if tstate < limitBottom {
		tstate = limitBottom
	}
	if endtstate > limitTop {
		endtstate = limitTop
	}
	tv.tstate = endtstate
	x, y := tvTstateToXY(tstate)
	endX, endY := tvTstateToXY(endtstate)
	for y <= endY {
		// horizontal 448 px : 48 border left + 256 screen/border  + 48 border right + 96 sync
		var hBorder, vBorder bool
		var nextX, lastX int
		vBorder = y < tvFirstScreenLine || y > tvLastScreenLine
		lastX = (tvTotalWidth - 1)
		if y == endY && endX < lastX {
			lastX = endX
		}
		for x <= lastX {
			// detect trace type
			if x < tvBorderLeft {
				hBorder = true
				nextX = tvBorderLeft - 1
			} else if x >= (tvBorderLeft + tvScreenWidth) {
				hBorder = true
				nextX = tvTotalWidth - 1
			} else {
				hBorder = vBorder
				nextX = tvBorderLeft + tvScreenWidth - 1
			}
			// paint scanline trace
			if lastX < nextX {
				nextX = lastX
			}
			if hBorder {
				tv.scanlineBorder(y, x, nextX, border)
			} else {
				tv.scanlineScreen(y, x, nextX)
			}
			// next trace
			x = nextX + 1
		}
		// next scanline
		x, y = 0, y+1
	}
}

func (tv *TvVideo) scanlineScreen(y, x1, x2 int) {
	var attr, data, mask byte
	var inkRgb, paperRgb uint32

	xx := x1 - tvBorderLeft
	yy := y - tvBorderTop
	scrAddr := (yy & 0xc0 << 5) | ((yy & 0x38) << 2) | ((yy & 0x07) << 8) | (xx >> 3)
	attrAddr := tvAttrAddr + ((yy & 0xf8) << 2) + (xx >> 3)
	bit := byte(xx) & 0x07
	mask = 1 << (7 - bit)
	readmem := true
	for x := x1; x <= x2; x++ {
		if readmem {
			readmem = false
			// read memory data & attr
			data = tv.srcdata[scrAddr]
			attr = tv.srcdata[attrAddr]
			scrAddr++
			attrAddr++
			// calculate ink + paper
			ink := attr & 0x07
			if (attr & 0x40) != 0 {
				ink |= 0x08
			}
			paper := (attr >> 3) & 0x07
			if (attr & 0x40) != 0 {
				paper |= 0x08
			}
			if attr&0x80 != 0 && tv.flash {
				ink, paper = paper, ink
			}
			inkRgb = tv.screen.GetColor(int(ink))
			paperRgb = tv.screen.GetColor(int(paper))
		}
		// paint pixel
		set := (data & mask) != 0
		if set {
			tv.screen.SetPixel(x, y, inkRgb)
		} else {
			tv.screen.SetPixel(x, y, paperRgb)
		}
		// next pixel
		bit++
		mask >>= 1
		if bit == 8 {
			bit = 0
			mask = 1 << 7 // 0x80
			readmem = true
		}
	}
}

func (tv *TvVideo) scanlineBorder(y, x1, x2 int, colour uint32) {
	for x := x1; x <= x2; x++ {
		tv.screen.SetPixel(x, y, colour)
	}
}

func tvTstateToXY(tstate int) (int, int) {
	tstate = tstate + tvHBorderTstates
	y := tstate / tvLineTstates
	x := tstate % tvLineTstates * tvTstatePixels
	return x, y
}
