package cpc

// -----------------------------------------------------------------------------
// Amstrad CPC - CRTC 6845 - Cathode Ray Tube Controller
// -----------------------------------------------------------------------------

// CRTC flag constans
const (
	CrtcVS    = 1   // VSync
	CrtcHS    = 2   // HSync
	CrtcHDT   = 4   // HorizontalDisplayedTotal
	CrtcVDT   = 8   // VerticalDisplayedTotal
	CrtcHT    = 16  // HorizontalTotal
	CrtcVT    = 32  // VerticalTotal
	CrtcMR    = 64  // MaximumRasterAddress
	CrtcVTadj = 128 // VerticalTotalAdjust
	CrtcVSf   = 256 // VerticalSyncPosition
)

// CRTC init register values
var crtcInitRegisters = []byte{0x3f, 0x28, 0x2e, 0x8e, 0x1f, 0x06, 0x19, 0x1b, 0x00, 0x07, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00}

// Crtc the Cathode Ray Tube Controller
type Crtc struct {
	cpc       *AmstradCPC
	selected  byte
	registers []byte
	flags     byte
}

// NewCrtc creates new CRTC
func NewCrtc(cpc *AmstradCPC) *Crtc {
	crtc := &Crtc{}
	crtc.cpc = cpc
	crtc.registers = make([]byte, 18)
	return crtc
}

// ReadRegister returns current register value
func (crtc *Crtc) ReadRegister() byte { return crtc.registers[crtc.selected] }

// WriteRegister returns current register value
func (crtc *Crtc) WriteRegister(data byte) { crtc.registers[crtc.selected] = data }

// Flags gets flags
func (crtc *Crtc) Flags() byte { return crtc.flags }

// SetFlags set flags
func (crtc *Crtc) SetFlags(flags byte) { crtc.flags = flags }

// AddFlags adds flags
func (crtc *Crtc) AddFlags(flags byte) { crtc.flags |= flags }

// RemoveFlags removes flags
func (crtc *Crtc) RemoveFlags(flags byte) { crtc.flags &= ^flags }

// Init the CRTC
func (crtc *Crtc) Init() { crtc.Reset() }

// Reset the CRTC
func (crtc *Crtc) Reset() {
	crtc.selected = 0
	crtc.flags = CrtcHDT | CrtcVDT
	copy(crtc.registers, crtcInitRegisters)
}

// Read reads data
func (crtc *Crtc) Read(port byte) byte {
	var data byte = 0xff
	if port == 0x03 {
		if (crtc.selected > 11) && (crtc.selected < 18) {
			data &= crtc.ReadRegister()
		} else {
			data = 0 // write only
		}
	}
	return data
}

// Write writes data
func (crtc *Crtc) Write(port byte, data byte) {
	switch port {
	case 0x00:
		if data < 18 {
			crtc.selected = data
		}
	case 0x01:
		if crtc.selected < 16 {
			crtc.WriteRegister(data)
		}
	}
}
