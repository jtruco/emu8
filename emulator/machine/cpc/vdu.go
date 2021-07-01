package cpc

import (
	"github.com/jtruco/emu8/emulator/device/video"
)

// -----------------------------------------------------------------------------
// Amstrad CPC - VDU Video
// -----------------------------------------------------------------------------

// CPC Video constants
const (
	videoScreenWidth  = 640
	videoScreenHeight = 200
	videoWidthScale   = 0.5
	videoHBorder      = 4 * 16
	videoVBorder      = 4 * 8
	videoTotalWidth   = videoScreenWidth + videoHBorder*2
	videoTotalHeight  = videoScreenHeight + videoVBorder*2
	videoLineBytes    = 0x800 // 2 KBytes
)

// CPC 464 RGBA colour palette (27 colors)
var cpcPaletteRGBA = []uint32{
	0xff808080, 0xff808080, 0xff80ff00, 0xff80ffff, 0xff800000, 0xff8000ff, 0xff808000, 0xff8080ff,
	0xff8000ff, 0xff80ffff, 0xff00ffff, 0xffffffff, 0xff0000ff, 0xffff00ff, 0xff0080ff, 0xffff80ff,
	0xff800000, 0xff80ff00, 0xff00ff00, 0xffffff00, 0xff000000, 0xffff0000, 0xff008000, 0xffff8000,
	0xff800080, 0xff80ff80, 0xff00ff80, 0xffffff80, 0xff000080, 0xffff0080, 0xff008080, 0xffff8080,
}

// CPC mode palette index tables
var (
	cpcMode0 [256][2]int
	cpcMode1 [256][4]int
	cpcMode2 [256][8]int
)

// VduVideo device
type VduVideo struct {
	screen    *video.Screen
	gatearray *GateArray
	crtc      *video.MC6845
	ram       [][]byte
	palette   []int
	mode      byte
	paintByte func(int, int, byte) int
	scanLine  uint16
	maxSLine  uint16
	minSLine  uint16
	lineBytes uint16
	firstX    int
	page      byte
	offset    uint16
}

// NewVduVideo creates a new vdu
func NewVduVideo(cpc *AmstradCPC) *VduVideo {
	vdu := new(VduVideo)
	vdu.screen = video.NewScreen(videoTotalWidth, videoTotalHeight, cpcPaletteRGBA)
	vdu.screen.SetScaleX(videoWidthScale)
	vdu.gatearray = cpc.gatearray
	vdu.crtc = cpc.crtc
	vdu.ram = make([][]byte, 4)
	vdu.ram[0] = cpc.memory.Bank(1).Data()
	vdu.ram[1] = cpc.memory.Bank(2).Data()
	vdu.ram[2] = cpc.memory.Bank(3).Data()
	vdu.ram[3] = cpc.memory.Bank(5).Data()
	return vdu
}

// Screen the video screen
func (vdu *VduVideo) Screen() *video.Screen { return vdu.screen }

// Init initializes video device
func (vdu *VduVideo) Init() { vdu.Reset() }

// Reset resets video device
func (vdu *VduVideo) Reset() {
	vdu.screen.Clear(0)
	vdu.updateCrtc()
}

// EndFrame updates screen video frame
func (vdu *VduVideo) EndFrame() {
	// nothing to do
}

// updateMode update gatearray
func (vdu *VduVideo) updateMode() {
	if vdu.mode == vdu.gatearray.mode {
		return
	}
	vdu.mode = vdu.gatearray.mode
	// mode : select paint byte function
	switch vdu.mode {
	case 0, 3: // 4 bpp
		vdu.paintByte = vdu.paintByte0
	case 1: // 2 bpp
		vdu.paintByte = vdu.paintByte1
	case 2: // 1 bpp
		vdu.paintByte = vdu.paintByte2
	}
	// palette
	vdu.palette = vdu.gatearray.palette
}

// updateCrtc update crtc screen options
func (vdu *VduVideo) updateCrtc() {
	// update gatearray
	vdu.updateMode()
	// scanline control
	firstY := (uint16(vdu.crtc.VerticalTotal+1)*uint16(vdu.crtc.MaxScanlineAddress+1) - videoTotalHeight) / 2
	firstY += 16 // vsync width
	vdu.minSLine = firstY
	vdu.maxSLine = firstY + videoTotalHeight - 1
	vdu.scanLine = 0
	// screen raster
	vdu.lineBytes = uint16(vdu.crtc.HorizontalDisplayed) << 1
	vdu.firstX = (videoTotalWidth - int(vdu.lineBytes)<<3) >> 1
	// memory offset
	vdu.page = (vdu.crtc.StartAddressHigh >> 4) & 0x03
	vdu.offset = (uint16(vdu.crtc.StartAddressHigh&0x03)<<8 | uint16(vdu.crtc.StartAddressLow)) << 1
}

// OnIntAck on interrupt ack
func (vdu *VduVideo) OnIntAck() {
	vdu.updateMode()
}

// OnVSync init a new screen
func (vdu *VduVideo) OnVSync() {
	vdu.updateCrtc()
}

// OnHSync renders a scanline
func (vdu *VduVideo) OnHSync() {

	// scanline control
	scanLine := vdu.scanLine
	vdu.scanLine++
	if scanLine < vdu.minSLine || scanLine > vdu.maxSLine {
		return // vertical border paddings
	}
	y := int(scanLine - vdu.minSLine)

	// vertical border scanlines
	border := vdu.screen.GetColour(vdu.gatearray.Border())
	if vdu.crtc.CurrentRow() >= vdu.crtc.VerticalDisplayed {
		vdu.paintLine(y, 0, videoTotalWidth, border)
		return
	}

	// render screen
	defer vdu.updateMode()

	// border left & right
	vdu.paintLine(y, 0, vdu.firstX, border)
	vdu.paintLine(y, videoTotalWidth-vdu.firstX, vdu.firstX, border)

	// render screen rasterline
	rasteraddr := uint16(vdu.crtc.CurrentLine()) << 11
	offset := vdu.offset + vdu.lineBytes*uint16(vdu.crtc.CurrentRow())
	x := vdu.firstX
	for i := uint16(0); i < vdu.lineBytes; i++ {
		offset &= 0x7ff
		addr := rasteraddr | offset
		x = vdu.paintByte(x, y, vdu.ram[vdu.page][addr])
		offset++
	}
}

// render functions

// paintByte0 paints mode 0 screen byte
func (vdu *VduVideo) paintByte0(x, y int, data byte) int {
	colour := vdu.screen.Palette()[vdu.palette[cpcMode0[data][0]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode0[data][1]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	return x
}

// paintByte1 paints mode 1 screen byte
func (vdu *VduVideo) paintByte1(x, y int, data byte) int {
	colour := vdu.screen.Palette()[vdu.palette[cpcMode1[data][0]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode1[data][1]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode1[data][2]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode1[data][3]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	vdu.screen.SetPixel(x, y, colour)
	x++
	return x
}

// paintByte2 paints mode 2 screen byte
func (vdu *VduVideo) paintByte2(x, y int, data byte) int {
	colour := vdu.screen.Palette()[vdu.palette[cpcMode2[data][0]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode2[data][1]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode2[data][2]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode2[data][3]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode2[data][4]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode2[data][5]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode2[data][6]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	colour = vdu.screen.Palette()[vdu.palette[cpcMode2[data][7]]]
	vdu.screen.SetPixel(x, y, colour)
	x++
	return x
}

// paintLine render a line of colour
func (vdu *VduVideo) paintLine(y, x1, width int, colour uint32) {
	x2 := x1 + width
	for x := x1; x < x2; x++ {
		vdu.screen.SetPixel(x, y, colour)
	}
}

// initialization

func init() {
	// mode color palette tables
	for data := 0; data < 256; data++ {
		// mode0 palette index table
		cpcMode0[data][0] = ((data & 0x80) >> 7) | ((data & 0x20) >> 3) | ((data & 0x08) >> 2) | ((data & 0x02) << 2)
		cpcMode0[data][1] = ((data & 0x40) >> 6) | ((data & 0x10) >> 2) | ((data & 0x04) >> 1) | ((data & 0x01) << 3)
		// mode1 palette index table
		cpcMode1[data][0] = ((data & 0x80) >> 7) | ((data & 0x08) >> 2)
		cpcMode1[data][1] = ((data & 0x40) >> 6) | ((data & 0x04) >> 1)
		cpcMode1[data][2] = ((data & 0x20) >> 5) | (data & 0x02)
		cpcMode1[data][3] = ((data & 0x10) >> 4) | ((data & 0x01) << 1)
		// mode2 palette index table
		cpcMode2[data][0] = (data & 0x80) >> 7
		cpcMode2[data][1] = (data & 0x40) >> 6
		cpcMode2[data][2] = (data & 0x20) >> 5
		cpcMode2[data][3] = (data & 0x10) >> 4
		cpcMode2[data][4] = (data & 0x08) >> 3
		cpcMode2[data][5] = (data & 0x04) >> 2
		cpcMode2[data][6] = (data & 0x02) >> 1
		cpcMode2[data][7] = (data & 0x01)
	}
}
