package cpc

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// The machine ends its devices frames : the emulation does not depend on
// the UI controllers being refreshed.
func TestCpcFrameLifecycleHeadless(t *testing.T) {
	m := New(AmstradCPC464, machine.Settings{}).(*AmstradCPC)
	controller.New(m, vfs.NewMemFileSystem())
	m.Init()
	m.psg.WriteRegister(7, 0x3e) // tone A
	m.psg.WriteRegister(0, 0x40) // tone A period
	m.psg.WriteRegister(8, 0x0f) // tone A max volume
	first := uint16(0)
	for frame := 1; frame <= 3; frame++ {
		m.BeginFrame()
		m.clock.Restart(cpcTStates)
		for m.clock.Tstates() < cpcTStates {
			m.Emulate()
		}
		m.EndFrame()
		buffer := m.psg.Buffer()
		max := uint16(0)
		for i := 0; i < buffer.Size(); i++ {
			if s := buffer.GetSample(i); s > max {
				max = s
			}
		}
		if frame == 1 {
			first = max
		}
		// same constant tone every frame : samples don't accumulate
		if max == 0 || max > first+first/8 {
			t.Errorf("frame %d: max sample %d, want 1..%d", frame, max, first+first/8)
		}
	}
}
