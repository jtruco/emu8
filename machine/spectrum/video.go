package spectrum

import (
	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/video"
)

// -----------------------------------------------------------------------------
// Video constants & vars
// -----------------------------------------------------------------------------

// Video screen constants
const (
	tvScreenWidth   = 256
	tvScreenHeight  = 192
	tvBorderLeft    = 48
	tvBorderRight   = 48
	tvBorderTop     = 64
	tvBorderBottom  = 56
	tvTotalWidth    = tvScreenWidth + tvBorderLeft + tvBorderRight
	tvTotalHeight   = tvScreenHeight + tvBorderTop + tvBorderBottom
	tvDisplayLeft   = 16 // Border : 32
	tvDisplayTop    = 40 // Border : 24
	tvDisplayWidth  = tvScreenWidth + 2*(tvBorderLeft-tvDisplayLeft)
	tvDisplayHeight = tvScreenHeight + 2*(tvBorderTop-tvDisplayTop)
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
	tvHRetraceTstates   = 48
	tvHBorderTstates    = 24
	tvFirstScreenTstate = 14336
	tvFirstScreenLine   = tvBorderTop
	tvLastScreenLine    = tvBorderTop + tvScreenHeight - 1
	tvBytePixels        = 8
	tvLineBytes         = tvScreenWidth / tvBytePixels
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

// TVVideo is the spectrum RF video device
type TVVideo struct {
	screen   *video.Screen // The video screen
	spectrum *Spectrum     // The Spectrum machine
	srcdata  []byte        // The screen data
	tstate   int           // Current videoframe tstate
	border   byte          // The border current colour index
	flash    bool          // Flash state
	frames   int           // Frame count
	accurate bool          // Accurate scanlines simulation
}

// NewTVVideo creates the video device
func NewTVVideo(spectrum *Spectrum) *TVVideo {
	tv := new(TVVideo)
	tv.screen = video.NewScreen(tvTotalWidth, tvTotalHeight, zxPalette)
	tv.screen.SetDisplay(tvDisplayLeft, tvDisplayTop, tvDisplayWidth, tvDisplayHeight)
	tv.spectrum = spectrum
	tv.srcdata = spectrum.VideoMemory().Data()
	spectrum.VideoMemory().OnAccess.Bind(tv.onVideoAccess)
	tv.accurate = true
	return tv
}

// onVideoAccess processes the bus event
func (tv *TVVideo) onVideoAccess(event device.IEvent) {
	bevent := event.(*device.BusEvent)
	if bevent.Code() == device.EventBusWrite {
		if tv.accurate && bevent.Address < tvVideoSize {
			tv.DoScanlines()
		}
	}
}

// SetAccurate sets de video emulation algorithm
func (tv *TVVideo) SetAccurate(accurate bool) {
	tv.accurate = accurate
}

// SetBorder sets de current border color
func (tv *TVVideo) SetBorder(colour byte) {
	if tv.accurate {
		tv.DoScanlines()
	}
	tv.border = colour
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
func (tv *TVVideo) Screen() *video.Screen { return tv.screen }

// Screen: simple and fast emulation

// paintScreen is a simple screen emulation
func (tv *TVVideo) paintScreen() {
	// 3 banks, of 8 rows, of 8 lines, of 32 cols
	baddr := 0
	y, sy := 0, tvBorderTop
	for b := 0; b < 3; b++ {
		for r := 0; r < 8; r++ {
			laddr := baddr + r*32
			for l := 0; l < 8; l++ {
				caddr := tvVideoAddr + laddr
				aaddr := tvAttrAddr + (y>>3)<<5
				sx := tvBorderLeft
				for c := 0; c < 32; c++ {
					data := tv.srcdata[caddr]
					attr := tv.srcdata[aaddr]
					tv.paintByte(sy, sx, data, attr, tv.flash)
					caddr++
					aaddr++
					sx += tvBytePixels
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
func (tv *TVVideo) paintBorder() {
	// Border Top, Bottom and Paper
	border := tv.screen.GetColour(int(tv.border))
	display := tv.screen.Display()
	for y := display.Y; y < tvBorderTop; y++ {
		tv.scanlineBorder(y, 0, tvTotalWidth-1, border)
	}
	for y := tvBorderTop + tvScreenHeight; y < display.Y+display.H; y++ {
		tv.scanlineBorder(y, 0, tvTotalWidth-1, border)
	}
	for y := tvBorderTop; y < display.Y+display.H; y++ {
		tv.scanlineBorder(y, 0, tvBorderLeft-1, border)
		tv.scanlineBorder(y, tvBorderLeft+tvScreenWidth, tvTotalWidth-1, border)
	}
}

// paintByte paints a byte
func (tv *TVVideo) paintByte(y, sx int, data, attr byte, flash bool) {
	var ink, paper, mask byte
	ink = attr & 0x07
	if (attr & 0x40) != 0 {
		ink |= 0x08
	}
	paper = (attr >> 3) & 0x07
	if (attr & 0x40) != 0 {
		paper |= 0x08
	}
	if (attr&0x80) != 0 && flash {
		ink, paper = paper, ink
	}
	mask = 0x80
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

// Screen : accurate emulation

// DoScanlines refresh TV scanlines
func (tv *TVVideo) DoScanlines() {
	// Spectrum 48k - Tv Scanlines timings
	// Vertical   : 16 Sl sync, 48 Sl border top, 192 Sl Screen, 56 Sl boder bottom
	// Horizontal : 128 Ts screen, 24 Ts border right, 48 Ts retrace, 24 TS border left
	// First screen (0,0) pixel Tstate = 14336 TS = 64 Scanlines * 224 Tstates
	display := tv.screen.Display()
	border := tv.screen.GetColour(int(tv.border))
	tstate := tv.tstate
	endtstate := tv.spectrum.Clock().Tstates()
	limitBottom := display.Y*tvLineTstates - tvHBorderTstates
	limitTop := (display.Y+display.H)*tvLineTstates - tvHBorderTstates
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
		hBorder, vBorder := false, y < tvFirstScreenLine || y > tvLastScreenLine
		nextX, lastX := x, (tvTotalWidth - 1)
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

func (tv *TVVideo) scanlineScreen(y, x1, x2 int) {
	var attr, data, ink, paper, mask byte

	xx := x1 - tvBorderLeft
	yy := y - tvBorderTop
	scrAddr := 0x800*(yy>>6) | tvLineBytes*((yy&0x38)>>3) | ((yy & 0x07) << 8) | xx>>3
	attrAddr := (yy>>3)<<5 + xx>>3 + tvAttrAddr
	bit := byte(xx % tvBytePixels)
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
			ink = attr & 0x07
			if (attr & 0x40) != 0 {
				ink |= 0x08
			}
			paper = (attr >> 3) & 0x07
			if (attr & 0x40) != 0 {
				paper |= 0x08
			}
			if attr&0x80 != 0 && tv.flash {
				ink, paper = paper, ink
			}
		}
		// paint pixel
		set := (data & mask) != 0
		if set {
			tv.screen.SetPixelIndex(x, y, int(ink))
		} else {
			tv.screen.SetPixelIndex(x, y, int(paper))
		}
		// next pixel
		bit++
		mask >>= 1
		if bit == tvBytePixels {
			bit = 0
			mask = 1 << 7 // 0x80
			readmem = true
		}
	}
}

func (tv *TVVideo) scanlineBorder(y, x1, x2 int, colour int32) {
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
