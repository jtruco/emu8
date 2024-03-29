package cpc

// -----------------------------------------------------------------------------
// Amstrad CPC - Gate Array
// -----------------------------------------------------------------------------

// GateArray constants
const (
	gaTotalPens    = 0x11
	gaBorderPen    = 0x10
	gaSlVsyncDelay = 0x02
	gaSlIntMax     = 0x34 // 52
	gaSlIntLimit   = 0x20 // 32
)

// GateArray for the CPC
type GateArray struct {
	cpc          *AmstradCPC
	palette      []int
	mode         byte
	pen          byte
	countSlInt   int
	countSlVsync int
}

// NewGateArray creates a GA
func NewGateArray(cpc *AmstradCPC) *GateArray {
	ga := new(GateArray)
	ga.cpc = cpc
	ga.palette = make([]int, gaTotalPens)
	ga.cpc.cpu.OnIntAck = ga.onInterruptAck
	ga.cpc.crtc.OnHSync = ga.onHSync
	ga.cpc.crtc.OnVSync = ga.onVSync
	return ga
}

// Init the GA
func (ga *GateArray) Init() {
	ga.Reset()
}

// Reset the GA
func (ga *GateArray) Reset() {
	ga.mode = 1
	ga.pen = 0
	for i := 0; i < gaTotalPens; i++ {
		ga.palette[i] = 0
	}
	// vdu scanline control
	ga.countSlInt = 0
	ga.countSlVsync = 0
}

// Config gets current config
func (ga *GateArray) Config() byte {
	var data byte = 0x80
	// mode / rom selection / interrupts
	data |= (ga.mode & 0x03)
	if !ga.cpc.memory.Map(cpcLowerROM).IsActive() {
		data |= 0x04
	}
	if !ga.cpc.memory.Map(cpcUpperROM).IsActive() {
		data |= 0x08
	}
	if ga.cpc.cpu.IntRq {
		data |= 0x10
	}
	return data
}

// Mode gets current mode
func (ga *GateArray) Mode() byte { return ga.mode }

// SetMode sets current mode
func (ga *GateArray) SetMode(mode byte) { ga.mode = mode }

// Pen gets current pen
func (ga *GateArray) Pen() byte { return ga.pen }

// SetPen sets current pen
func (ga *GateArray) SetPen(pen byte) { ga.pen = pen }

// Border returns the border color
func (ga *GateArray) Border() int { return ga.palette[gaBorderPen] }

// Palette returns the active pen colors
func (ga *GateArray) Palette() []int { return ga.palette }

// SetInk set ink colour & palette
func (ga *GateArray) SetInk(ink byte) {
	ga.palette[ga.pen] = int(ink)
}

// Bus Input / Output

// Read read from gatearray
func (ga *GateArray) Read() byte {
	return 0xff
}

// Write data to gatearray
func (ga *GateArray) Write(data byte) {
	switch data >> 6 {
	case 0: // select pen
		if (data & 0x10) == 0x00 {
			ga.SetPen(data & 0x0f)
		} else {
			ga.SetPen(gaBorderPen) // border
		}
	case 1: // set colur
		ga.SetInk(data & 0x1f)
	case 2: // Video Mode & ROM
		// mode
		ga.SetMode(data & 0x03)
		// rom selection
		ga.cpc.memory.Map(cpcLowerROM).SetActive(data&0x04 == 0)
		ga.cpc.memory.Map(cpcUpperROM).SetActive(data&0x08 == 0)
		// interrupts
		if (data & 0x10) != 0 {
			// clear pending interrupts
			ga.cpc.cpu.InterruptRequest(false)
			ga.countSlInt = 0
		}
	case 3:
		// RAM memory management (not implemented)
	}
}

// Emulation

// Emulate gate array
func (ga *GateArray) Emulate(tstates int) {
	// 4 MHz gatearray emulation

	// 1 MHz clock emulation
	gaTstates := tstates >> 2
	ga.cpc.crtc.Emulate(gaTstates)
	ga.cpc.psg.Emulate(gaTstates)
}

// onHSync on CRTC hsync callback
func (ga *GateArray) onHSync() {
	ga.countSlInt++
	if ga.countSlVsync == 0 {
		if ga.countSlInt == gaSlIntMax {
			ga.cpc.cpu.InterruptRequest(true)
			ga.countSlInt = 0
		}
	} else {
		ga.countSlVsync--
		if ga.countSlInt >= gaSlIntLimit {
			ga.cpc.cpu.InterruptRequest(true)
		}
		ga.countSlInt = 0
	}
	ga.cpc.video.OnHSync()
}

// onVSync on CRTC vsync callback
func (ga *GateArray) onVSync() {
	ga.countSlVsync = gaSlVsyncDelay
	ga.cpc.video.OnVSync()
}

// onInterruptAck interrupt ack
func (ga *GateArray) onInterruptAck() bool {
	ga.countSlInt &= 0x01F // Unset bit 5
	ga.cpc.video.OnIntAck()
	return false
}
