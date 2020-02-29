package cpc

// -----------------------------------------------------------------------------
// Amstrad CPC - Gate Array
// -----------------------------------------------------------------------------

// GateArray for the CPC
type GateArray struct {
	cpc     *AmstradCPC
	pen     byte
	palette []byte
	mode    byte
	slCount int
	slTotal int
	ts      int
}

// NewGateArray creates a GA
func NewGateArray(cpc *AmstradCPC) *GateArray {
	ga := &GateArray{}
	ga.cpc = cpc
	ga.palette = make([]byte, 17)
	return ga
}

// Init the GA
func (ga *GateArray) Init() {
	ga.Reset()
}

// Reset the GA
func (ga *GateArray) Reset() {
	ga.pen = 0
	for i := 1; i < 17; i++ {
		ga.palette[i] = 0
	}
	ga.mode = 1
	ga.slCount = 0
	ga.slTotal = 0
	ga.ts = 0
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
func (ga *GateArray) Border() byte { return ga.palette[0x10] }

// Palette returns the active pen colors
func (ga *GateArray) Palette() []byte { return ga.palette }

// SetInk set ink colour & palette
func (ga *GateArray) SetInk(ink byte) {
	if ga.pen == 0x10 {
		ga.palette[ga.pen] = ink
	} else {
		ga.palette[ga.pen] = ink
	}
}

// Emulate gate array
func (ga *GateArray) Emulate(tstates int) {
	// FIXME : Simple interrupt emulation
	// 52 scanlines / 1sl ~ 256 Ts
	// 312 sl -> vsync
	ga.ts += tstates
	if ga.ts >= 256 { // 1sl
		ga.ts &= 0xff
		ga.slCount++
		ga.slTotal++
		if ga.slCount == 52 {
			ga.slCount = 0
			ga.cpc.cpu.Interrupt()
		}
		if ga.slTotal == 4 { // end vsync
			ga.cpc.crtc.RemoveFlags(CrtcVS)
		} else if ga.slTotal == 312 { // Frame
			ga.slTotal = 0
			ga.cpc.crtc.AddFlags(CrtcVS)
		}
	}
}

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
			ga.SetPen(0x10) // border
		}
	case 1: // set colur
		ga.SetInk(data & 0x1f)
	case 2: // Video Mode & ROM
		// mode
		ga.SetMode(data & 0x03)
		if (data & 0x10) != 0 {
			ga.slCount = 0
			ga.slCount = 0
		}
		// rom selection
		ga.cpc.lowerRom.SetActive(data&0x04 == 0)
		ga.cpc.upperRom.SetActive(data&0x08 == 0)
		// interrupts
		if (data & 0x10) != 0 {
			// TODO : clear pending interrupts
			ga.slCount = 0 // reset GA scanline counter
		}
	case 3:
		// RAM memory management (not implemented)
	}
}
