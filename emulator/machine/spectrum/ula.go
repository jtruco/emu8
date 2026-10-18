package spectrum

import "github.com/jtruco/emu8/emulator/device/bus"

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

// Memory timing

// ulaMemWriteTstates is the delay from the start of a CPU memory write cycle
// until the RAM is written. The Z80 write cycle lasts 3 Ts and the data lands
// at its end (as in FUSE), but the bus access hooks run at the cycle start,
// before the CPU adds the 3 Ts. The TV beam must be rendered up to that point.
const ulaMemWriteTstates = 3

// Contention table

// IO contention pages
var ulaIoPageContention = [4]bool{false, true, false, false}

// ULA contention delay table
var ulaDelayTable [zxTStates + tvLineTstates]int

func init() {
	// contention table
	tstate := tvScreenTstate - 1
	for y := 0; y < tvPaperHeight; y++ {
		for x := 0; x < tvPaperTstates; x += 8 {
			ulaDelayTable[tstate+x+0] = 6
			ulaDelayTable[tstate+x+1] = 5
			ulaDelayTable[tstate+x+2] = 4
			ulaDelayTable[tstate+x+3] = 3
			ulaDelayTable[tstate+x+4] = 2
			ulaDelayTable[tstate+x+5] = 1
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
	// render the TV beam up to the write, while video memory holds the old byte
	if code == bus.EventWrite && address < tvVideoSize {
		ula.spectrum.tv.Update(ula.spectrum.clock.Tstates() + ulaMemWriteTstates)
	}
}

// Device

// Init initializes ULA
func (ula *ULA) Init() {
	ula.lastRead = ulaInDefault
}

// Reset resets ULA
func (ula *ULA) Reset() {
	// nothing to do
}

// DataBus

// Read bus at address. The Z80 core adds the 4 Ts of the I/O cycle; the ULA
// adds the contention of the port (preIO and postIO).
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

// Write bus at address. The Z80 core adds the 4 Ts of the I/O cycle; the ULA
// adds the contention of the port (preIO and postIO).
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

// I/O contention patterns (C : contended T-state, N : not contended), by
// high byte (contended or not) and bit 0 (ULA port when 0) :
//
//	not contended, bit 0 = 1 : N:4
//	not contended, bit 0 = 0 : N:1, C:3
//	contended,     bit 0 = 1 : C:1, C:1, C:1, C:1
//	contended,     bit 0 = 0 : C:1, C:3
//
// The Z80 core adds the 4 Ts of the cycle (1 T before the device access, 3 Ts
// after); preIO and postIO only add the wait states.

// preIO contention : wait states of the first T-state (the core has already
// added it, so the cycle started 1 T before)
func (ula *ULA) preIO(address uint16) {
	if ula.isContended(address) {
		start := ula.spectrum.clock.Tstates() - 1
		ula.spectrum.clock.Add(ulaDelayTable[start])
	}
}

// postIO contention : wait states of the last 3 T-states (added by the core
// after the device access)
func (ula *ULA) postIO(address uint16) {
	clock := ula.spectrum.clock
	switch {
	case (address & 0x0001) == 0: // C:3
		clock.Add(ulaDelayTable[clock.Tstates()])
	case ula.isContended(address): // C:1, C:1, C:1
		tstate, wait := clock.Tstates(), 0
		for i := 0; i < 3; i++ {
			delay := ulaDelayTable[tstate]
			wait += delay
			tstate += delay + 1
		}
		clock.Add(wait)
	}
}

// doContention applies clock contention
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
