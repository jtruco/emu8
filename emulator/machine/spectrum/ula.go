package spectrum

// -----------------------------------------------------------------------------
// ULA constants & vars
// -----------------------------------------------------------------------------

const ulaInDefault = 0xff

// Audio

/* assume all three tone channels together match the beeper volume (ish).
 * Must be <=127 for all channels; 50+2+(24*3) = 124.
 * (Now scaled up for 16-bit.)
 */
const (
	amplRate   = 1
	amplBeeper = (50 * 256) >> amplRate
	amplTape   = (2 * 256) >> amplRate
	amplAyTone = (24 * 256) >> amplRate
)

var beeperMap = []uint16{0, amplTape, amplBeeper, (amplBeeper + amplTape)}

// Contention table

// IO contention pages
var ulaIoPageContention = [4]bool{false, true, false, false}

// ULA contention delay table
var ulaDelayTable [zxTStates + tvLineTstates]int

func init() {
	// contention table
	tstate := tvFirstScreenTstate - 1
	for y := 0; y < tvScreenHeight; y++ {
		for x := 0; x < tvScreenWidth; x += 16 {
			tstatex := x / tvTstatePixels
			ulaDelayTable[tstate+tstatex+0] = 6
			ulaDelayTable[tstate+tstatex+1] = 5
			ulaDelayTable[tstate+tstatex+2] = 4
			ulaDelayTable[tstate+tstatex+3] = 3
			ulaDelayTable[tstate+tstatex+4] = 2
			ulaDelayTable[tstate+tstatex+5] = 1
		}
		tstate += tvLineTstates
	}
}

// -----------------------------------------------------------------------------
// ULA
// -----------------------------------------------------------------------------

// ULA is the Unit Logic Array
type ULA struct {
	spectrum *Spectrum // The spectrum machine
	lastRead byte      // Last read value
}

// NewULA creates
func NewULA(spectrum *Spectrum) *ULA {
	ula := new(ULA)
	ula.spectrum = spectrum
	spectrum.VideoMemory().OnAccess = ula.onVideoAccess
	return ula
}

// onVideoAccess processes the bus event
func (ula *ULA) onVideoAccess(code int, address uint16) {
	ula.doContention(0)
}

// Device

// Init initializes ULA
func (ula *ULA) Init() {
	ula.lastRead = ulaInDefault
}

// Reset resets ULA
func (ula *ULA) Reset() {
	// nothing to to
}

// DataBus

// Read bus at address
func (ula *ULA) Read(address uint16) byte {
	var result byte = 0xff
	ula.preIO(address)
	ula.postIO(address)
	if (address & 0x0001) == 0x00 { // ULA selected
		result = ula.lastRead
		// Read keyboard state
		scan := byte(address>>8) ^ 0xff
		mask := byte(1)
		for row := 0; row < 8; row++ {
			if (scan & mask) != 0 { // scan row
				result &= ula.spectrum.keyboard.rowstates[row]
			}
			mask <<= 1
		}
		// Read tape state
		if ula.spectrum.tape.IsPlaying() && ula.spectrum.tape.EarHigh() {
			result &^= 0x40
		}
	}
	if (address & 0x00e0) == 0 { // Kempston
		result &= ula.spectrum.joystick.GetKempston()
	}
	return result
}

// Write bus at address
func (ula *ULA) Write(address uint16, data byte) {
	ula.preIO(address)
	if (address & 0x0001) == 0 { // ULA selected
		// border
		ula.spectrum.tv.SetBorder(data & 0x07)
		// beeper & tape output : EAR(bit 4) and MIC(bit 3) output
		tstate := ula.spectrum.clock.Tstates()
		beeper := int(data&0x18) >> 3
		if ula.spectrum.tape.IsPlaying() {
			if ula.spectrum.tape.EarHigh() {
				beeper |= 2
			} else {
				beeper &^= 2
			}
		}
		ula.spectrum.beeper.SetLevel(tstate, beeper)
		// default read
		ula.lastRead = ulaInDefault
		if (data & 0x18) == 0 { // ISSUE 2
			ula.lastRead ^= 0x40
		}
	}
	ula.postIO(address)
}

// preIO contention
func (ula *ULA) preIO(address uint16) {
	if ula.isContended(address) {
		ula.doContention(1)
	} else {
		ula.spectrum.clock.Add(1)
	}
}

// postIO contention
func (ula *ULA) postIO(address uint16) {
	if (address & 0x0001) != 0 {
		if ula.isContended(address) {
			ula.doContention(1)
			ula.doContention(1)
			ula.doContention(1)
		} else {
			ula.spectrum.clock.Add(3)
		}
	} else {
		ula.doContention(3)
	}
}

// doContention aplies clock contention
func (ula *ULA) doContention(tstates int) {
	delay := ulaDelayTable[ula.spectrum.clock.Tstates()] + tstates
	if delay > 0 {
		ula.spectrum.clock.Add(delay)
	}
}

// isContended true if address access is contended
func (ula *ULA) isContended(address uint16) bool {
	page := address >> 14
	return ulaIoPageContention[page]
}
