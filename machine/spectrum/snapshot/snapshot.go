package snapshot

import (
	"github.com/jtruco/emu8/cpu/z80"
)

// -----------------------------------------------------------------------------
// Snapshot
// -----------------------------------------------------------------------------

// Snapshot ZX Spectrum 16k / 48k snap
type Snapshot struct {
	z80.State                 // Z80 state
	Tstates   int             // CPU tstates
	Border    byte            // ULA current border
	Memory    [48 * 1024]byte // Spectrum memory (48k)
}

// NewSnapshot returns a new ZX Spectrum snap
func NewSnapshot() *Snapshot {
	snap := &Snapshot{}
	snap.State.Init()
	return snap
}
