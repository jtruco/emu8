package spectrum

import (
	"github.com/jtruco/emu8/cpu"
	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/memory"
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
	bank     *memory.Bank  // The spectrum video memory bank
	clock    cpu.Clock     // The CPU clock
	tstate   int           // Current videoframe tstate
	palette  []int32       // The video palette
	border   byte          // The border current colour index
	flash    bool          // Flash state
	frames   int           // Frame count
	accurate bool          // Accurate scanlines simulation
}

// NewTVVideo creates the video device
func NewTVVideo(spectrum *Spectrum) *TVVideo {
	tv := &TVVideo{}
	tv.palette = zxPalette
	tv.screen = video.NewScreen(tvTotalWidth, tvTotalHeight, tv.palette)
	tv.screen.SetDisplay(video.Rect{X: tvDisplayLeft, Y: tvDisplayTop, W: tvDisplayWidth, H: tvDisplayHeight})
	tv.bank = spectrum.memory.GetBankMap(1).Bank()
	tv.bank.AddBusListener(tv)
	tv.clock = spectrum.clock
	return tv
}

// ProcessBusEvent processes the bus event
func (tv *TVVideo) ProcessBusEvent(event *device.BusEvent) {
	if tv.accurate && event.Address < tvVideoSize {
		if event.Type == device.EventBusWrite && event.Order == device.OrderBefore {
			tv.DoScanlines()
		}
	}
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
	screendata := tv.bank.Data()
	flash := tv.flash

	// 3 banks, of 8 rows, of 8 lines, of 32 cols
	baddr := 0
	y := 0
	for b := 0; b < 3; b++ {
		for r := 0; r < 8; r++ {
			laddr := baddr + r*32
			for l := 0; l < 8; l++ {
				caddr := tvVideoAddr + laddr
				aaddr := tvAttrAddr + (y/8)*32
				for c := 0; c < 32; c++ {
					data := screendata[caddr]
					attr := screendata[aaddr]
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
	display := tv.screen.Display()
	for y := display.Y; y < tvBorderTop; y++ {
		for x := 0; x < tvTotalWidth; x++ {
			tv.screen.SetPixel(x, y, border)
		}
	}
	for y := tvBorderTop + tvScreenHeight; y < tvTotalHeight; y++ {
		for x := 0; x < tvTotalWidth; x++ {
			tv.screen.SetPixel(x, y, border)
		}
	}
	for y := tvBorderTop; y < display.Y+display.H; y++ {
		for x := 0; x < tvBorderLeft; x++ {
			tv.screen.SetPixel(x, y, border)
		}
		for x := tvBorderLeft + tvScreenWidth; x < tvTotalWidth; x++ {
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
	sx := tvBorderLeft + (c << 3)
	y += tvBorderTop
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
	// TODO
}
