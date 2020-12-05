// Package format implements CPC file formats
package format

import "github.com/jtruco/emu8/emulator/device/cpu/z80"

// -----------------------------------------------------------------------------
// CPC 464/646 Snapshot
// -----------------------------------------------------------------------------

// Snapshot CPC snapshot version 1
type Snapshot struct {
	z80.State                  // Z80 state
	gaState                    // GateArray state
	crtcState                  // CRTC state
	ppiState                   // PPI state
	psgState                   // PSG state
	RomSelect byte             // ROM selection
	Memory    [64 * 0x400]byte // CPC RAM (64k)
}

// gaState gatearray state
type gaState struct {
	GaSelectedPen byte
	GaPenColours  [17]byte
	GaMultiConfig byte
	GaRAMSelect   byte
}

type crtcState struct {
	CrtcSelected  byte
	CrtcRegisters [18]byte
}

type ppiState struct {
	PpiPortA   byte
	PpiPortB   byte
	PpiPortC   byte
	PpiControl byte
}

type psgState struct {
	PsgSelected  byte
	PsgRegisters [16]byte
}

// NewSnapshot returns a new snapshop
func NewSnapshot() *Snapshot {
	snap := new(Snapshot)
	snap.State.Init()
	return snap
}

// -----------------------------------------------------------------------------
// Format common functions
// -----------------------------------------------------------------------------

// readWord reads a 16 bit LSB unsgined integer
func readWord(data []byte, pos int) uint16 {
	return uint16(data[pos]) | (uint16(data[pos+1]) << 8)
}

// writeWord writes a 16 bit LSB unsgined integer
func writeWord(data []byte, pos int, value uint16) {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
}
