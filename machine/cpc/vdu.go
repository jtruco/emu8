package cpc

import (
	"github.com/jtruco/emu8/device/video"
)

// -----------------------------------------------------------------------------
// Amstrad CPC - VDU Video
// -----------------------------------------------------------------------------

// CPC Video constants
const (
	videoScreenWidth  = 640
	videoScreenHeight = 200
	videoHRows        = videoScreenWidth >> 3 // bytes * row
	videoVLines       = videoScreenHeight
	videoVCols        = videoVLines >> 3
	videoHBorder      = 64
	videoVBorder      = 32
	videoVSpare       = 48
	videoTotalWidth   = videoScreenWidth + videoHBorder*2
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

// CPC mode palette index tables
var (
	cpcMode0 [256][2]int
	cpcMode1 [256][4]int
)

// VduVideo device
type VduVideo struct {
	screen    *video.Screen
	gatearray *GateArray
	crtc      *video.MC6845
	banks     [][]byte
	border    int
	paintByte func(int, int, byte) (int, int)
}

// NewVduVideo creates a new vdu
func NewVduVideo(cpc *AmstradCPC) *VduVideo {
	vdu := new(VduVideo)
	vdu.screen = video.NewScreen(videoTotalWidth, videoTotalHeight, cpcPalette)
	vdu.screen.SetWScale(videoWidthScale)
	vdu.gatearray = cpc.gatearray
	vdu.crtc = cpc.crtc
	vdu.banks = make([][]byte, 4)
	vdu.banks[0] = cpc.memory.Map(1).Bank().Data()
	vdu.banks[1] = cpc.memory.Map(2).Bank().Data()
	vdu.banks[2] = cpc.memory.Map(3).Bank().Data()
	vdu.banks[3] = cpc.memory.Map(5).Bank().Data()
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
	vdu.updateOptions()
	vdu.paintScreen()
	vdu.paintBorder()
}

// updateOptions update screen options
func (vdu *VduVideo) updateOptions() {
	// select paint byte function
	switch vdu.gatearray.mode {
	case 0, 3: // 4 bpp
		vdu.paintByte = vdu.paintByte0
	case 1: // 2 bpp
		vdu.paintByte = vdu.paintByte1
	case 2: // 1 bpp
		vdu.paintByte = vdu.paintByte2
	}
}

// paintScreen paints screen
func (vdu *VduVideo) paintScreen() {
	// select crtc bank & offset
	r12 := vdu.crtc.ReadRegister(video.MC6845StartAddressHigh)
	r13 := vdu.crtc.ReadRegister(video.MC6845StartAddressLow)
	// fixme : crtc bank switch
	bank := vdu.banks[(r12>>4)&0x03]
	offset := (((uint16(r12) & 0x03) << 8) | uint16(r13)) << 1
	// paint screen data
	x, y := videoHBorder, videoVBorder
	row, col := 0, 0
	for addr := offset; addr < videoTotalBytes; addr++ {
		x, y = vdu.paintByte(x, y, bank[addr])
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

// paintByte0 paints mode 0 screen byte
func (vdu *VduVideo) paintByte0(x, y int, data byte) (int, int) {
	palette := vdu.gatearray.palette

	colour := palette[cpcMode0[data][0]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[cpcMode0[data][1]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++

	return x, y
}

// paintByte1 paints mode 1 screen byte
func (vdu *VduVideo) paintByte1(x, y int, data byte) (int, int) {
	palette := vdu.gatearray.palette

	colour := palette[cpcMode1[data][0]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[cpcMode1[data][1]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[cpcMode1[data][2]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[cpcMode1[data][3]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++

	return x, y
}

// paintByte2 paints mode 2 screen byte
func (vdu *VduVideo) paintByte2(x, y int, data byte) (int, int) {
	palette := vdu.gatearray.palette

	colour := palette[(data&0x80)>>7]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[(data&0x40)>>6]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[(data&0x20)>>5]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[(data&0x10)>>4]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[(data&0x08)>>3]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[(data&0x04)>>2]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[(data&0x02)>>1]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = palette[(data & 0x01)]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++

	return x, y
}

// paintBorder paints border
func (vdu *VduVideo) paintBorder() {
	border := vdu.gatearray.Border()
	if border == vdu.border {
		return
	}
	vdu.border = border
	colour := vdu.screen.GetColour(int(border))
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

func init() {
	// mode0 palette index table
	for i := 0; i < 256; i++ {
		cpcMode0[i][0] = ((i & 0x80) >> 7) | ((i & 0x20) >> 3) | ((i & 0x08) >> 2) | ((i & 0x02) << 2)
		cpcMode0[i][1] = ((i & 0x40) >> 6) | ((i & 0x10) >> 2) | ((i & 0x04) >> 1) | ((i & 0x01) << 3)
	}
	// mode1 palette index table
	for data := 0; data < 256; data++ {
		cpcMode1[data][0] = ((data & 0x80) >> 7) | ((data & 0x08) >> 2)
		cpcMode1[data][1] = ((data & 0x40) >> 6) | ((data & 0x04) >> 1)
		cpcMode1[data][2] = ((data & 0x20) >> 5) | (data & 0x02)
		cpcMode1[data][3] = ((data & 0x10) >> 4) | ((data & 0x01) << 1)
	}
}
