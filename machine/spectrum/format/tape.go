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
