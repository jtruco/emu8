package cpc

import "github.com/jtruco/emu8/device/video"

// -----------------------------------------------------------------------------
// Amstrad CPC - VDU Video
// -----------------------------------------------------------------------------

// CPC Video constants
const (
	videoHMode0       = 160
	videoHMode1       = 320
	videoHMode2       = 640
	videoHRows        = videoHMode2 >> 3 // bytes * row
	videoVLines       = 200
	videoVCols        = videoVLines >> 3
	videoScreenWidth  = videoHMode2
	videoScreenHeight = videoVLines
	videoHBorder      = 32
	videoVBorder      = 32
	videoVSpare       = 48
	videoTotalWidth   = videoHMode2 + videoHBorder*2
	videoTotalHeight  = videoScreenHeight + videoVBorder*2
	videoWidthScale   = 0.5
	videoTotalBytes   = videoHRows*videoVLines + videoVSpare*8
)

// CPC 464 RGB colour palette (27 colors)
var cpcPalette = []int32{
	0x808080, 0x808080, 0x00ff80, 0xffff80, 0x000080, 0xff0080, 0x008080, 0xff8080,
	0xff0080, 0xffff80, 0xffff00, 0xffffff, 0xff0000, 0xff00ff, 0xff8000, 0xff80ff,
	0x000080, 0x00ff80, 0x00ff00, 0x00ffff, 0x000000, 0x0000ff, 0x008000, 0x0080ff,
	0x800080, 0x80ff80, 0x80ff00, 0x80ffff, 0x800000, 0x8000ff, 0x808000, 0x8080ff,
}

// VduVideo device
type VduVideo struct {
	cpc     *AmstradCPC
	screen  *video.Screen
	srcdata []byte
	border  byte
}

// NewVduVideo creates a new vdu
func NewVduVideo(cpc *AmstradCPC) *VduVideo {
	vdu := &VduVideo{}
	vdu.cpc = cpc
	vdu.screen = video.NewScreen(videoTotalWidth, videoTotalHeight, cpcPalette)
	vdu.screen.SetWScale(videoWidthScale)
	vdu.srcdata = cpc.memory.Map(5).Bank().Data()
	return vdu
}

// Screen the video screen
func (vdu *VduVideo) Screen() *video.Screen { return vdu.screen }

// Init initializes video device
func (vdu *VduVideo) Init() { vdu.Reset() }

// Reset resets video device
func (vdu *VduVideo) Reset() {
	vdu.screen.Clear(0)
}

// EndFrame updates screen video frame
func (vdu *VduVideo) EndFrame() {
	vdu.paintBorder()
	vdu.paintScreen()
}

// paintScreen paints screen
func (vdu *VduVideo) paintScreen() {
	mode := vdu.cpc.gatearray.mode
	palette := vdu.cpc.gatearray.palette
	x, y := videoHBorder, videoVBorder
	row, col := 0, 0
	for addr := 0; addr < videoTotalBytes; addr++ {
		data := vdu.srcdata[addr]
		// paintbyte
		switch mode {
		case 0, 3: // 4 bpp
			idx := ((data & 0x80) >> 7) | ((data & 0x20) >> 3) | ((data & 0x08) >> 2) | ((data & 0x02) << 2)
			colour := int(palette[idx])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			idx = ((data & 0x40) >> 6) | ((data & 0x10) >> 2) | ((data & 0x04) >> 1) | ((data & 0x01) << 3)
			colour = int(palette[idx])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
		case 1: // 2 bpp
			idx := ((data & 0x80) >> 6) | ((data & 0x08) >> 3)
			colour := int(palette[idx])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			idx = ((data & 0x40) >> 5) | ((data & 0x04) >> 2)
			colour = int(palette[idx])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			idx = ((data & 0x20) >> 4) | ((data & 0x02) >> 1)
			colour = int(palette[idx])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			idx = ((data & 0x10) >> 3) | (data & 0x01)
			colour = int(palette[idx])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
		case 2: // 1 bpp
			colour := int(palette[(data&0x80)>>7])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			colour = int(palette[(data&0x40)>>6])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			colour = int(palette[(data&0x20)>>5])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			colour = int(palette[(data&0x10)>>4])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			colour = int(palette[(data&0x08)>>3])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			colour = int(palette[(data&0x04)>>2])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			colour = int(palette[(data&0x02)>>1])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
			colour = int(palette[(data & 0x01)])
			vdu.screen.SetPixelIndex(x, y, colour)
			x++
		}
		// next byte
		row++
		if row == videoHRows {
			row, x = 0, videoHBorder
			col++
			y += 8
			if col == videoVCols {
				col = 0
				y -= (videoVLines - 1)
				addr += videoVSpare // spare bytes
			}
		}
	}
}

// paintBorder paints border
func (vdu *VduVideo) paintBorder() {
	border := vdu.cpc.gatearray.Border()
	if border == vdu.border {
		return
	}
	vdu.border = border
	colour := cpcPalette[int(border)]
	display := vdu.screen.Display()
	// Border Top, Bottom and Paper
	for y := display.Y; y < videoVBorder; y++ {
		vdu.paintLine(y, 0, videoTotalWidth-1, colour)
	}
	for y := videoVBorder; y < videoVBorder+videoScreenHeight; y++ {
		vdu.paintLine(y, 0, videoHBorder-1, colour)
		vdu.paintLine(y, videoHBorder+videoScreenWidth, videoTotalWidth-1, colour)
	}
	for y := videoVBorder + videoScreenHeight; y < display.Y+display.H; y++ {
		vdu.paintLine(y, 0, videoTotalWidth-1, colour)
	}
}

func (vdu *VduVideo) paintLine(y, x1, x2 int, colour int32) {
	for x := x1; x <= x2; x++ {
		vdu.screen.SetPixel(x, y, colour)
	}
}
