package spectrum

// -----------------------------------------------------------------------------
// Audio constants & vars
// -----------------------------------------------------------------------------

/* assume all three tone channels together match the beeper volume (ish).
 * Must be <=127 for all channels; 50+2+(24*3) = 124.
 * (Now scaled up for 16-bit.)
 */
const (
	amplBeeper = (50 * 256)
	amplTape   = (2 * 256)
	amplAyTone = (24 * 256)
	volumeRate = 2
)

var beeperMap = []uint16{0, amplTape >> volumeRate, amplBeeper >> volumeRate, (amplBeeper + amplTape) >> volumeRate}

// -----------------------------------------------------------------------------
// ULA
// -----------------------------------------------------------------------------

// ULA is the Unit Logic Array
type ULA struct {
	spectrum *Spectrum
}

// NewULA creates
func NewULA(spectrum *Spectrum) *ULA {
	ula := &ULA{}
	ula.spectrum = spectrum
	return ula
}

// Device

// Init initializes ULA
func (ula *ULA) Init() {}

// Reset resets ULA
func (ula *ULA) Reset() {}

// DataBus

// Access access bus
func (ula *ULA) Access(address uint16) {}

// Read bus at address
func (ula *ULA) Read(address uint16) byte {
	var result byte = 0xff
	if (address & 0x0001) == 0x0000 {
		// Read keyboard state
		var row uint
		for row = 0; row < 8; row++ {
			if (address & (1 << (uint16(row) + 8))) == 0 { // bit held low, so scan this row
				result &= ula.spectrum.keyboard.rowstates[row]
			}
		}
	}
	return result
}

// Write bus at address
func (ula *ULA) Write(address uint16, data byte) {
	if (address & 0x0001) == 0 {
		// border
		ula.spectrum.tv.SetBorder(ula.spectrum.clock.Tstates(), data&0x07)

		// beeper
		// EAR(bit 4) and MIC(bit 3) output
		tstate := ula.spectrum.clock.Tstates()
		beeper := int(data&0x18) >> 3
		ula.spectrum.beeper.SetLevel(tstate, beeper)
	}
}
