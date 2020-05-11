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
	videoWidthScale   = 0.5
	videoHBorder      = 4 * 16
	videoVBorder      = 4 * 8
	videoTotalWidth   = videoScreenWidth + videoHBorder*2
	videoTotalHeight  = videoScreenHeight + videoVBorder*2
	videoLineBytes    = 0x800 // 2 KBytes
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
	startX    int
	page      byte
	offset    uint16
}

// NewVduVideo creates a new vdu
func NewVduVideo(cpc *AmstradCPC) *VduVideo {
	vdu := new(VduVideo)
	vdu.screen = video.NewScreen(videoTotalWidth, videoTotalHeight, cpcPalette)
	vdu.screen.SetWScale(videoWidthScale)
	vdu.gatearray = cpc.gatearray
	vdu.crtc = cpc.crtc
	vdu.ram = make([][]byte, 4)
	vdu.ram[0] = cpc.memory.Map(1).Bank().Data()
	vdu.ram[1] = cpc.memory.Map(2).Bank().Data()
	vdu.ram[2] = cpc.memory.Map(3).Bank().Data()
	vdu.ram[3] = cpc.memory.Map(5).Bank().Data()
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
	// crtc screen params
	vdu.lineBytes = uint16(vdu.crtc.HorizontalDisplayed) << 1
	vdu.startX = (videoTotalWidth - int(vdu.lineBytes)<<3) >> 1
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
	vdu.paintLine(y, 0, vdu.startX, border)
	vdu.paintLine(y, videoTotalWidth-vdu.startX, vdu.startX, border)

	// render screen rasterline
	rasteraddr := uint16(vdu.crtc.CurrentLine()) << 11
	offset := vdu.offset + vdu.lineBytes*uint16(vdu.crtc.CurrentRow())
	x := vdu.startX
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
	colour := vdu.palette[cpcMode0[data][0]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode0[data][1]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	return x
}

// paintByte1 paints mode 1 screen byte
func (vdu *VduVideo) paintByte1(x, y int, data byte) int {
	colour := vdu.palette[cpcMode1[data][0]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode1[data][1]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode1[data][2]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode1[data][3]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	return x
}

// paintByte2 paints mode 2 screen byte
func (vdu *VduVideo) paintByte2(x, y int, data byte) int {
	colour := vdu.palette[cpcMode2[data][0]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode2[data][1]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode2[data][2]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode2[data][3]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode2[data][4]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode2[data][5]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode2[data][6]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	colour = vdu.palette[cpcMode2[data][7]]
	vdu.screen.SetPixelIndex(x, y, colour)
	x++
	return x
}

// paintLine render a line of colour
func (vdu *VduVideo) paintLine(y, x1, width int, colour int32) {
	x2 := x1 + width
	for x := x1; x < x2; x++ {
		vdu.screen.SetPixel(x, y, colour)
	}
}

// initialization

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
	// mode2 palette index table
	for data := 0; data < 256; data++ {
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
