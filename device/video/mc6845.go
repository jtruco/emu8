package video

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Motorola MC6845 - CRTC controller
// -----------------------------------------------------------------------------

// MC6845 constants
const (
	mc6845Nreg = 18 // 18 registers
)

// MC6845 register data
var (
	mc6845Defaults = [mc6845Nreg]byte{ // Amstrad CPC 464 default values
		0x3f, 0x28, 0x34, 0x34, 0x14, 0x08, 0x10, 0x13, 0x00,
		0x0b, 0x49, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	mc6845Masks = [mc6845Nreg]byte{
		0xff, 0xff, 0xff, 0xff, 0x7f, 0x1f, 0x7f, 0x7f, 0x03,
		0x1f, 0x1f, 0x1f, 0x3f, 0xff, 0x3f, 0xff, 0x3f, 0xff}
)

// MC6845 Crtc Device
type MC6845 struct {
	registers   [mc6845Nreg]*byte
	regdefaults [mc6845Nreg]byte
	selected    byte
	// registers
	rHorizontalTotal        byte
	rHorizontalDisplayed    byte
	rHorizontalSyncPosition byte
	rSyncWidths             byte
	rVerticalTotal          byte
	rVerticalTotalAdjust    byte
	rVerticalDisplayed      byte
	rVerticalSyncPosition   byte
	rInterlaceAndSkew       byte
	rMaxScanlineAddress     byte
	rCursorStart            byte
	rCursorEnd              byte
	rStartAddressHigh       byte
	rStartAddressLow        byte
	rCursorHigh             byte
	rCursorLow              byte
	rLightPenHigh           byte
	rLightPenLow            byte
	// control
	currentCol  byte
	currentRow  byte
	currentLine byte
	hSyncWidth  byte
	vSyncWidth  byte
	hSyncCount  byte
	vSyncCount  byte
	inHSync     bool
	inVSync     bool
	// callbacks
	OnHSync device.Callback
	OnVSync device.Callback
}

// NewMC6845 creates new CRTC
func NewMC6845() *MC6845 {
	mc := new(MC6845)
	mc.regdefaults = mc6845Defaults
	mc.registers = [mc6845Nreg]*byte{
		&mc.rHorizontalTotal,
		&mc.rHorizontalDisplayed,
		&mc.rHorizontalSyncPosition,
		&mc.rSyncWidths,
		&mc.rVerticalTotal,
		&mc.rVerticalTotalAdjust,
		&mc.rVerticalDisplayed,
		&mc.rVerticalSyncPosition,
		&mc.rInterlaceAndSkew,
		&mc.rMaxScanlineAddress,
		&mc.rCursorStart,
		&mc.rCursorEnd,
		&mc.rStartAddressHigh,
		&mc.rStartAddressLow,
		&mc.rCursorHigh,
		&mc.rCursorLow,
		&mc.rLightPenHigh,
		&mc.rLightPenLow,
	}
	return mc
}

// Properties

// InHSync in HSync
func (mc *MC6845) InHSync() bool { return mc.inHSync }

// InVSync in VSync
func (mc *MC6845) InVSync() bool { return mc.inVSync }

// Device interface

// Init the CRTC
func (mc *MC6845) Init() { mc.Reset() }

// Reset the CRTC
func (mc *MC6845) Reset() {
	mc.selected = 0
	for i := byte(0); i < mc6845Nreg; i++ {
		mc.writeRegister(i, mc.regdefaults[i])
	}
	mc.currentCol = 0
	mc.currentLine = 0
	mc.currentRow = 0
	mc.hSyncCount = 0
	mc.vSyncCount = 0
	mc.inHSync = false
	mc.inVSync = false
}

// Emulate emulates Tstates
func (mc *MC6845) Emulate(tstates int) {
	for i := 0; i < tstates; i++ {
		mc.OnClock()
	}
}

// OnClock emulates one clock cycle
func (mc *MC6845) OnClock() {
	// hsync duration control
	if mc.hSyncCount > 0 {
		mc.hSyncCount--
		if mc.hSyncCount == 0 {
			mc.inHSync = false
		}
	}
	// onclock moves one character
	mc.currentCol++
	// scanline control
	if mc.currentCol > mc.rHorizontalTotal {
		mc.currentCol = 0
		// vsync duration control
		if mc.vSyncCount > 0 {
			mc.vSyncCount--
			if mc.vSyncCount == 0 {
				mc.inVSync = false
			}
		}
		// new line
		mc.currentLine++
		if mc.currentLine > mc.rMaxScanlineAddress {
			mc.currentLine = 0
			mc.currentRow++
			if mc.currentRow > mc.rVerticalTotal {
				mc.currentRow = 0
			}
		}
		// vsync control
		if !mc.inVSync && mc.currentRow == mc.rVerticalSyncPosition {
			mc.inVSync = true
			mc.vSyncCount = mc.vSyncWidth
			if mc.OnVSync != nil {
				mc.OnVSync()
			}
		}
	} else if !mc.inHSync && mc.currentCol == mc.rHorizontalSyncPosition {
		// hsync control
		mc.inHSync = true
		mc.hSyncCount = mc.hSyncWidth
		if mc.OnHSync != nil {
			mc.OnHSync()
		}
	}
}

// IO operations

// Read reads data
func (mc *MC6845) Read(port byte) byte {
	var data byte = 0xff
	if port == 0x03 {
		data &= mc.readSelected()
	}
	return data
}

// Write writes data
func (mc *MC6845) Write(port byte, data byte) {
	switch port {
	case 0x00:
		mc.selectRegister(data & 0x1f)
	case 0x01:
		mc.writeSelected(data)
	}
}

// register operations

// readRegister selects current register
func (mc *MC6845) selectRegister(selected byte) {
	mc.selected = selected
}

// readSelected returns current register value
func (mc *MC6845) readSelected() byte {
	return mc.readRegister(mc.selected)
}

// readRegister returns register value
func (mc *MC6845) readRegister(register byte) byte {
	if (register > 11) && (register < mc6845Nreg) {
		return *mc.registers[register]
	}
	return 0 // write only
}

// writeSelected writes value to selected register
func (mc *MC6845) writeSelected(data byte) {
	mc.writeRegister(mc.selected, data)
}

// writeRegister writes value to register
func (mc *MC6845) writeRegister(register, data byte) {
	*mc.registers[register] = data & mc6845Masks[register]

	// HSync & VSync widths
	if register == 0x03 {
		mc.hSyncWidth = mc.rSyncWidths & 0x0f
		if mc.hSyncWidth == 0 {
			mc.hSyncWidth = 0x10
		}
		mc.vSyncWidth = (mc.rSyncWidths >> 4) & 0x0f
		if mc.vSyncWidth == 0 {
			mc.vSyncWidth = 0x10
		}
	}
}
