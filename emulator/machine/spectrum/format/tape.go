package format

// -----------------------------------------------------------------------------
// ZX Spectrum tape common constants
// -----------------------------------------------------------------------------

// Tape file types
const (
	TapeFileProgram        = 0
	TapeFileNumberArray    = 1
	TapeFileCharacterArray = 2
	TapeFileCode           = 3
)

// Tape play states
const (
	tapeStateStart = iota
	tapeStatePilot
	tapeStateSync
	tapeStateByte
	tapeStateBit1
	tapeStateBit2
	tapeStatePause
	tapeStatePauseStop
	tapeStateStop
)

// Tape tstate constants
const (
	tapeTimingPilot   = 2168
	tapeTimingSync1   = 667
	tapeTimingSync2   = 735
	tapeTimingZero    = 855
	tapeTimingOne     = 1710
	tapeHeaderPulses  = 8063
	tapeDataPulses    = 3223
	tapeEndBlockPause = 3494400 // 3494400 Ts/s
	tapeTimingEoB     = tapeEndBlockPause / 1000
)

// -----------------------------------------------------------------------------
// Format common functions
// -----------------------------------------------------------------------------

// readInt reads a 16 bit LSB unsgined integer as integer
func readInt(data []byte, pos int) int {
	return int(readWord(data, pos))
}

// readWord reads a 16 bit LSB unsgined integer
func readWord(data []byte, pos int) uint16 {
	return uint16(data[pos]) | (uint16(data[pos+1]) << 8)
}

// writeWord writes a 16 bit LSB unsgined integer
func writeWord(data []byte, pos int, value uint16) {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
}

// readIntN reads LSB unsgined integer as integer
func readIntN(data []byte, pos int, len int) int {
	value := uint(data[pos])
	if len > 1 && len < 4 {
		lshift := uint(8)
		for len > 1 {
			pos++
			value += uint(data[pos]) << lshift
			lshift += 8
			len--
		}
	}
	return int(value)
}

// readString reads valid characters to string
func readString(data []byte, pos int, len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		char := data[pos+i]
		if char > 31 || char < 128 {
			bytes[i] = char
		} else {
			bytes[i] = '?'
		}
	}
	return string(bytes)
}
