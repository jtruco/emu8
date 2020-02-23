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
	videoScreenWidth  = videoHMode1
	videoScreenHeight = videoVLines
	videoHBorder      = 32
	videoVBorder      = 32
	videoVSpare       = 48
	videoTotalWidth   = videoScreenWidth + videoHBorder*2
	videoTotalHeight  = videoScreenHeight + videoVBorder*2
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
			idx0 := ((data & 0x80) >> 7) | ((data & 0x20) >> 3) | ((data & 0x08) >> 2) | ((data & 0x02) << 2)
			idx1 := ((data & 0x40) >> 6) | ((data & 0x10) >> 2) | ((data & 0x04) >> 1) | ((data & 0x01) << 3)
			vdu.screen.SetPixelIndex(x, y, int(palette[idx0]))
			x++
			vdu.screen.SetPixelIndex(x, y, int(palette[idx0]))
			x++
			vdu.screen.SetPixelIndex(x, y, int(palette[idx1]))
			x++
			vdu.screen.SetPixelIndex(x, y, int(palette[idx1]))
			x++
		case 1: // 2 bpp
			idx0 := ((data & 0x80) >> 6) | ((data & 0x08) >> 3)
			idx1 := ((data & 0x40) >> 5) | ((data & 0x04) >> 2)
			idx2 := ((data & 0x20) >> 4) | ((data & 0x02) >> 1)
			idx3 := ((data & 0x10) >> 3) | (data & 0x01)
			vdu.screen.SetPixelIndex(x, y, int(palette[idx0]))
			x++
			vdu.screen.SetPixelIndex(x, y, int(palette[idx1]))
			x++
			vdu.screen.SetPixelIndex(x, y, int(palette[idx2]))
			x++
			vdu.screen.SetPixelIndex(x, y, int(palette[idx3]))
			x++
		case 2: // 1 bpp
			idx0 := (data & 0x80) >> 7
			idx1 := (data & 0x40) >> 6
			idx2 := (data & 0x20) >> 5
			idx3 := (data & 0x10) >> 4
			idx4 := (data & 0x08) >> 3
			idx5 := (data & 0x04) >> 2
			idx6 := (data & 0x02) >> 1
			idx7 := (data & 0x01)
			vdu.screen.SetPixel(x, y, rgbBlend(palette[idx0], palette[idx1]))
			x++
			vdu.screen.SetPixel(x, y, rgbBlend(palette[idx2], palette[idx3]))
			x++
			vdu.screen.SetPixel(x, y, rgbBlend(palette[idx4], palette[idx5]))
			x++
			vdu.screen.SetPixel(x, y, rgbBlend(palette[idx6], palette[idx7]))
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

// mode 2 - colour helper

func rgbBlend(idx1, idx2 byte) int32 {
	const fintpl = 2
	c1 := cpcPalette[idx1]
	if idx1 == idx2 {
		return c1
	}
	c2 := cpcPalette[idx2]
	r1, g1, b1 := rgb(c1)
	r2, g2, b2 := rgb(c2)
	r3 := (r2-r1)>>fintpl + r1
	g3 := (g2-g1)>>fintpl + r1
	b3 := (b2-b1)>>fintpl + r1
	c3 := torgb(r3, g3, b3)
	return c3
}

func rgb(col int32) (byte, byte, byte) {
	r := byte(col & 0xff0000 >> 16)
	g := byte(col & 0x00ff00 >> 8)
	b := byte(col & 0x0000ff)
	return r, g, b
}

func torgb(r, g, b byte) int32 {
	return int32(r)<<16 | int32(g)<<8 | int32(b)
}