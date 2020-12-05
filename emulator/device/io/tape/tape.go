// Package tape contains tape and drive components
package tape

// -----------------------------------------------------------------------------
// Tape components
// -----------------------------------------------------------------------------

// Tape represents a tape file
type Tape interface {
	// Info gets the tape information
	Info() *Info
	// Blocks gets the tape blocks
	Blocks() []Block
	// Load tape data. Returns false on error.
	Load(data []byte) bool
	// Play tape
	Play(control *Control)
}

// BlockInfo tape block information
type BlockInfo struct {
	Type   byte // Block type
	Index  int  // Block index
	Offset int  // Block offset
	Length int  // Block length
}

// Block is a tape block
type Block interface {
	Info() *BlockInfo
	Data() []byte
}

// Info tape information
type Info struct {
	Name string // Tape name
}
