package spectrum

// -----------------------------------------------------------------------------
// ULA constants & vars
// -----------------------------------------------------------------------------

// ULA issue constants
const (
	ulaInDefault  = 0xff
	ulaIssueMask  = 0x18 // ISSUE 2
	ulaIssueMask3 = 0x10 // ISSUE 3+
)

// Audio

// Beeper + Tape + 3 * AY = 122
const (
	amplRate   = 7 // uint16
	amplBeeper = 48 << amplRate
	amplTape   = 2 << amplRate
	amplAyTone = 24 << amplRate
)

var zxBeeperMap = []uint16{0, amplTape, amplBeeper, (amplBeeper + amplTape)}

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
	spectrum.memory.Map(zxVideoMemory).OnAccess = ula.onVideoAccess
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
	defer ula.postIO(address)
	if (address & 0x0001) == 0x00 { // ULA selected
		result = ula.lastRead
		// Read keyboard state
		scan := byte(address>>8) ^ 0xff
		result &= ula.spectrum.keyboard.GetState(scan)
		// Read tape state
		if ula.spectrum.tape.IsPlaying() && ula.spectrum.tape.EarHigh() {
			result &^= 0x40
		}
	}
	if (address & 0x00e0) == 0 { // Kempston selected
		result &= ula.spectrum.joystick.State()
	}
	return result
}

// Write bus at address
func (ula *ULA) Write(address uint16, data byte) {
	ula.preIO(address)
	defer ula.postIO(address)
	if (address & 0x0001) == 0 { // ULA selected
		// border
		ula.spectrum.tv.SetBorder(data & 0x07)
		// beeper & tape output : EAR(bit 4) and MIC(bit 3) output
		tstate := ula.spectrum.clock.Tstates()
		beeper := int(data&0x18) >> 3
		if ula.spectrum.tape.IsPlaying() && ula.spectrum.tape.EarHigh() {
			beeper &^= 0x1 // Loud tape sound
		}
		ula.spectrum.beeper.SetLevel(tstate, beeper)
		// default read
		ula.lastRead = ulaInDefault
		if (data & ulaIssueMask) == 0 {
			ula.lastRead ^= 0x40
		}
	}
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
