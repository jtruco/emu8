package spectrum

import (
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

// Sync constants
const (
	tvScanlines     = 312
	tvVSyncLines    = 16
	tvTopLine       = 48 + tvVSyncLines
	tvBottomLine    = tvScanlines - 56
	tvLineTstates   = 224
	tvHSyncTstates  = 200
	tvBorderLines   = tvBorderHeight
	tvBorderTstates = tvBorderWidth / 2
	tvPaperTstates  = tvPaperWidth / 2
	tvScreenTstate  = 14336
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
	// Scanline addresses table
	tvScanlineAddr [tvScanlines][2]int

	// ZX Spetrum 16/48k RGBA colour palette
	zxPaletteRGBA = []uint32{
		/* Bright 0 (black, blue, red, magenta, green, cyan, yellow, white) */
		0xff000000, 0xffc00000, 0xff0000c0, 0xffc000c0, 0xff00c000, 0xffc0c000, 0xff00c0c0, 0xffc0c0c0,
		/* BRIGHT 1 (black, blue, red, magenta, green, cyan, yellow, white) */
		0xff000000, 0xffff0000, 0xff0000ff, 0xffff00ff, 0xff00ff00, 0xffffff00, 0xff00ffff, 0xffffffff,
	}
)

// UlaTv is the spectrum RF video device
type UlaTv struct {
	screen     *video.Screen // The video screen
	scrdata    []byte        // The screen data
	border     byte          // The border current colour index
	borderRgb  uint32        // The border rgb colour
	flash      bool          // Flash state
	frames     int           // Frame count
	tstate     int           // Video TState
	scanline   int           // Video Scanline
	x, y       int           // Screen beam position
	scrAddr    int           // Screen memory address
	attAddr    int           // Attribute memory address
	scrByte    byte          // Screen byte
	bitCount   byte          // Bit position
	bitMask    byte          // Bit mask
	attInk     uint32        // Current ink color
	attPaper   uint32        // Current paper color
	isOnScreen bool
	isOnBorder bool
}

func NewUlaTv(spectrum *Spectrum) *UlaTv {
	var tv = new(UlaTv)
	tv.scrdata = spectrum.memory.Bank(zxVideoMemory).Data()
	tv.screen = video.NewScreen(tvScreenWidth, tvScreenHeight, zxPaletteRGBA)
	return tv
}

// SetBorder sets de current border color
func (tv *UlaTv) SetBorder(color byte) {
	tv.border = color
	tv.borderRgb = tv.screen.GetColor(int(color))
}

// Device

// Init initializes video device
func (tv *UlaTv) Init() { tv.Reset() }

// Reset resets video device
func (tv *UlaTv) Reset() {
	tv.screen.Clear(0)
	tv.border = 7
	tv.flash = false
	tv.onVSync()
}

// Video

// Screen the video screen
func (tv *UlaTv) Screen() *video.Screen { return tv.screen }

// EndFrame updates screen video frame
func (tv *UlaTv) EndFrame() {
	tv.frames++
	tv.flash = (tv.frames & 0x10) == 0
}

// Emulation

func (tv *UlaTv) Emulate(tstates int) {
	for i := 0; i < tstates; i++ {
		tv.OnClock()
	}
}

func (tv *UlaTv) OnClock() {
	// Spectrum 48k - Tv Scanlines timings
	// Vertical   : 312 Sl = 16 Sl sync, 48 Sl border top, 192 Sl Screen, 56 Sl border bottom
	// Horizontal : 224 Ts = 128 Ts screen, 24 Ts border right, 48 Ts retrace, 24 TS border left
	// First screen (0,0) pixel Tstate = 14336 TS = 64 Scanlines * 224 Tstates

	// tstate & scanline control
	if tv.tstate == tvHSyncTstates {
		tv.onHSync()
		if tv.scanline == tvScanlines {
			tv.onVSync()
		}
	} else if tv.tstate == tvLineTstates {
		tv.tstate = 0
	}

	// screen : 1 ts == 2px
	if tv.isOnScreen && tv.isTstateOnScreen() {
		if tv.isOnBorder || tv.isTstateOnBorder() {
			tv.paintBorderPixel()
			tv.paintBorderPixel()
		} else {
			tv.paintScreenPixel()
			tv.paintScreenPixel()
		}
	}
	tv.tstate++
}

func (tv *UlaTv) onHSync() {
	tv.x = 0
	tv.scanline++
	tv.isOnScreen = tv.isScanlineOnScreen()
	tv.isOnBorder = tv.isScanlineOnBorder()
	if tv.isOnScreen {
		tv.y = tv.scanline - tvTopLine + tvBorderLines
		if !tv.isOnBorder {
			var addrs = tvScanlineAddr[tv.scanline]
			tv.scrAddr = addrs[0]
			tv.attAddr = addrs[1]
		}
	}
}

func (tv *UlaTv) onVSync() {
	tv.x = 0
	tv.y = 0
	tv.scanline = 0
	tv.isOnBorder = tv.isScanlineOnBorder()
	tv.isOnScreen = tv.isScanlineOnScreen()
}

func (tv *UlaTv) isScanlineOnScreen() bool {
	return tv.scanline >= tvTopLine-tvBorderLines &&
		tv.scanline < tvBottomLine+tvBorderLines
}

func (tv *UlaTv) isScanlineOnBorder() bool {
	return tv.scanline < tvTopLine || tv.scanline >= tvBottomLine
}

func (tv *UlaTv) isTstateOnScreen() bool {
	return tv.tstate < tvPaperTstates+tvBorderTstates ||
		tv.tstate >= tvLineTstates-tvBorderTstates
}

func (tv *UlaTv) isTstateOnBorder() bool {
	return tv.tstate >= tvPaperTstates
}

func (tv *UlaTv) paintBorderPixel() {
	tv.screen.SetPixel(tv.x, tv.y, tv.borderRgb)
	tv.x++
}

func (tv *UlaTv) paintScreenPixel() {
	// read video mem
	if tv.bitCount == 0 {
		tv.bitMask = 0x80
		tv.scrByte = tv.scrdata[tv.scrAddr]
		attByte := tv.scrdata[tv.attAddr]
		attInk := attByte & 0x07
		if (attByte & 0x40) != 0 {
			attInk |= 0x08
		}
		attPaper := (attByte >> 3) & 0x07
		if (attByte & 0x40) != 0 {
			attPaper |= 0x08
		}
		if ((attByte & 0x80) != 0) && tv.flash {
			attInk = attInk ^ attPaper
			attPaper = attInk ^ attPaper
			attInk = attInk ^ attPaper
		}
		tv.attInk = tv.screen.GetColor(int(attInk))
		tv.attPaper = tv.screen.GetColor(int(attPaper))
	}
	// paint bit
	if (tv.scrByte & tv.bitMask) != 0 {
		tv.screen.SetPixel(tv.x, tv.y, tv.attInk)
	} else {
		tv.screen.SetPixel(tv.x, tv.y, tv.attPaper)
	}
	// next bit / byte
	tv.x++
	tv.bitCount++
	tv.bitMask >>= 1
	if tv.bitCount == 8 {
		tv.bitCount = 0
		tv.scrAddr++
		tv.attAddr++
	}
}

func init() {
	// Scanline addresses table
	for y := 0; y < tvPaperHeight; y++ {
		scrAddr := tvVideoAddr + ((y & 0xc0) << 5) + ((y & 0x38) << 2) + ((y & 0x07) << 8)
		attAddr := tvAttrAddr + ((y & 0xf8) << 2)
		tvScanlineAddr[tvTopLine+y][0] = scrAddr
		tvScanlineAddr[tvTopLine+y][1] = attAddr
	}
}
