package format

// -----------------------------------------------------------------------------
// ZX Spectrum tape common constants
// -----------------------------------------------------------------------------

// Tape EAR constants
const (
	TapeEarOn   = 0xff
	TapeEarOff  = 0xbf
	TapeEarMask = 0x40
)

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
	tapeStateLeader
	tapeStateSync
	tapeStateNewByte
	tapeStateNewBit
	tapeStateHalf2
	tapeStatePause
	tapeStatePauseStop
	tapeStateStop
)

// Tape tstate constants
const (
	tapeLeaderLenght  = 2168
	tapeSync1Lenght   = 667
	tapeSync2Lenght   = 735
	tapeZeroLenght    = 855
	tapeOneLenght     = 1710
	tapeHeaderPulses  = 8063
	tapeDataPulses    = 3223
	tapeEndBlockPause = 3500000
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

// readIntN reads LSB unsgined integer as integer
func readIntN(data []byte, pos int, len int) int {
	value := uint(data[pos])
	if len > 1 && len < 4 {
		lshift := 8
		for len > 1 {
			pos++
			value += uint(data[pos]) << lshift
			lshift += 8
			len--
		}
	}
	return int(value)
}
