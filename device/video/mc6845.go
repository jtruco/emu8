package video

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Motorola MC6845 - CRTC controller
// -----------------------------------------------------------------------------

// MC6845 constants
const (
	MC6845Nreg = 0X12 // 18 registers
)

// MC6845 register constants
const (
	MC6845HorizontalTotal = iota
	MC6845HorizontalDisplayed
	MC6845HorizontalSyncPosition
	MC6845SyncWidths
	MC6845VerticalTotal
	MC6845VerticalTotalAdjust
	MC6845VerticalDisplayed
	MC6845VerticalSyncPosition
	MC6845InterlaceAndSkew
	MC6845MaxScanlineAddress
	MC6845CursorStart
	MC6845CursorEnd
	MC6845StartAddressHigh
	MC6845StartAddressLow
	MC6845CursorHigh
	MC6845CursorLow
	MC6845LightPenHigh
	MC6845LightPenLow
)

// MC6845 register data
var (
	MC6845Defaults = [MC6845Nreg]byte{ // Amstrad CPC 464 default values
		0x3f, 0x28, 0x2e, 0x8e, 0x26, 0x00, 0x19, 0x1e, 0x00,
		0x07, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00}
	mc6845Masks = [MC6845Nreg]byte{
		0xff, 0xff, 0xff, 0xff, 0x7f, 0x1f, 0x7f, 0x7f, 0x03,
		0x1f, 0x1f, 0x1f, 0x3f, 0xff, 0x3f, 0xff, 0x3f, 0xff}
)

// MC6845 Crtc Device
type MC6845 struct {
	registers [MC6845Nreg]*byte
	defaults  [MC6845Nreg]byte
	selected  byte
	// registers
	HorizontalTotal        byte
	HorizontalDisplayed    byte
	HorizontalSyncPosition byte
	SyncWidths             byte
	VerticalTotal          byte
	VerticalTotalAdjust    byte
	VerticalDisplayed      byte
	VerticalSyncPosition   byte
	InterlaceAndSkew       byte
	MaxScanlineAddress     byte
	CursorStart            byte
	CursorEnd              byte
	StartAddressHigh       byte
	StartAddressLow        byte
	CursorHigh             byte
	CursorLow              byte
	LightPenHigh           byte
	LightPenLow            byte
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
	mc.defaults = MC6845Defaults
	mc.registers = [MC6845Nreg]*byte{
		&mc.HorizontalTotal,
		&mc.HorizontalDisplayed,
		&mc.HorizontalSyncPosition,
		&mc.SyncWidths,
		&mc.VerticalTotal,
		&mc.VerticalTotalAdjust,
		&mc.VerticalDisplayed,
		&mc.VerticalSyncPosition,
		&mc.InterlaceAndSkew,
		&mc.MaxScanlineAddress,
		&mc.CursorStart,
		&mc.CursorEnd,
		&mc.StartAddressHigh,
		&mc.StartAddressLow,
		&mc.CursorHigh,
		&mc.CursorLow,
		&mc.LightPenHigh,
		&mc.LightPenLow,
	}
	return mc
}

// Properties

// CurrentCol current row column
func (mc *MC6845) CurrentCol() byte { return mc.currentCol }

// CurrentRow current row
func (mc *MC6845) CurrentRow() byte { return mc.currentRow }

// CurrentLine current line in row
func (mc *MC6845) CurrentLine() byte { return mc.currentLine }

// InHSync in HSync
func (mc *MC6845) InHSync() bool { return mc.inHSync }

// InVSync in VSync
func (mc *MC6845) InVSync() bool { return mc.inVSync }

// SetDefaults sets default register values
func (mc *MC6845) SetDefaults(defaults [MC6845Nreg]byte) {
	mc.defaults = defaults
}

// Device interface

// Init the CRTC
func (mc *MC6845) Init() { mc.Reset() }

// Reset the CRTC
func (mc *MC6845) Reset() {
	mc.selected = 0
	for i := byte(0); i < MC6845Nreg; i++ {
		mc.WriteRegister(i, mc.defaults[i])
	}
	mc.currentCol = 0
	mc.currentLine = 0
	mc.currentRow = 0
	mc.hSyncCount = 0
	mc.vSyncCount = 0
	mc.inHSync = false
	mc.inVSync = false
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
		mc.SelectRegister(data & 0x1f)
	case 0x01:
		mc.writeSelected(data)
	}
}

// register operations

// SelectRegister selects current register
func (mc *MC6845) SelectRegister(selected byte) {
	mc.selected = selected
}

// readSelected returns current register value
func (mc *MC6845) readSelected() byte {
	return mc.ReadRegister(mc.selected)
}

// ReadRegister returns register value
func (mc *MC6845) ReadRegister(register byte) byte {
	if (register > 11) && (register < MC6845Nreg) {
		return *mc.registers[register]
	}
	return 0 // write only
}

// writeSelected writes value to selected register
func (mc *MC6845) writeSelected(data byte) {
	mc.WriteRegister(mc.selected, data)
}

// WriteRegister writes value to register
func (mc *MC6845) WriteRegister(register, data byte) {
	*mc.registers[register] = data & mc6845Masks[register]

	// HSync & VSync widths
	if register == 0x03 {
		mc.hSyncWidth = mc.SyncWidths & 0x0f
		if mc.hSyncWidth == 0 {
			mc.hSyncWidth = 0x10
		}
		mc.vSyncWidth = (mc.SyncWidths >> 4) & 0x0f
		if mc.vSyncWidth == 0 {
			mc.vSyncWidth = 0x10
		}
	}
}

// emulation

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
	if mc.currentCol > mc.HorizontalTotal {
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
		if mc.currentLine > mc.MaxScanlineAddress {
			mc.currentLine = 0
			mc.currentRow++
			if mc.currentRow > mc.VerticalTotal {
				mc.currentRow = 0
			}
		}
		// vsync control
		if !mc.inVSync && mc.currentRow == mc.VerticalSyncPosition {
			mc.inVSync = true
			mc.vSyncCount = mc.vSyncWidth
			if mc.OnVSync != nil {
				mc.OnVSync()
			}
		}
	} else {
		// hsync control
		if !mc.inHSync && mc.currentCol == mc.HorizontalSyncPosition {
			mc.inHSync = true
			mc.hSyncCount = mc.hSyncWidth
			if mc.OnHSync != nil {
				mc.OnHSync()
			}
		}
	}
}
