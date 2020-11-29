package tape

// -----------------------------------------------------------------------------
// Tape Control
// -----------------------------------------------------------------------------

// EAR level constants
const (
	LevelLow  = 0
	LevelHigh = 0x80
	LevelMask = 0x80
)

// Control struct for tape playback
type Control struct {
	Playing    bool  // Tape drive is playing
	Ear        byte  // Tape EAR level
	State      int   // Playback state
	Tstate     int64 // Last clock Tstate
	Timeout    int   // Timeout of current state
	Block      Block // Current tape block
	NumBlocks  int   // Total number of blocks on tape
	BlockIndex int   // Current block index
	BlockPos   int   // Curren block pos
}

// DataAtPos returns data byte at current block position
func (control *Control) DataAtPos() byte {
	return control.Block.Data()[control.BlockPos]
}

// EndOfBlock position at end of block data
func (control *Control) EndOfBlock() bool {
	return control.BlockPos >= control.Block.Info().Length
}

// EndOfTape blockindex at end of tape blocks
func (control *Control) EndOfTape() bool {
	return control.BlockIndex >= control.NumBlocks
}

func (control *Control) reset() {
	control.Playing = false
	control.Ear = LevelLow
	control.State = 0
	control.Tstate = 0
	control.Timeout = 0
	control.BlockIndex = 0
	control.BlockPos = 0
}
