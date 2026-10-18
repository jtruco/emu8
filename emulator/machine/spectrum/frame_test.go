package spectrum

import "testing"

// runHeadlessFrame runs a full frame without the UI controllers (no Refresh)
func runHeadlessFrame(m *Spectrum) {
	m.BeginFrame()
	m.clock.Restart(zxTStates)
	for m.clock.Tstates() < zxTStates {
		m.Emulate()
	}
	m.EndFrame()
}

// The machine ends its devices frames : the emulation does not depend on
// the UI controllers being refreshed.
func TestSpectrumFrameLifecycleHeadless(t *testing.T) {
	m := newTestSpectrum(t)
	m.memory.Bank(0).Load(0, []byte{0x18, 0xfe}) // ROM : JR $ (loop)
	m.ula.Write(0x00fe, 0x12)                    // red border, beeper high
	for frame := 1; frame <= 3; frame++ {
		runHeadlessFrame(m)
		if m.tv.tstate != 0 {
			t.Fatalf("frame %d: TV beam not restarted (tstate %d)", frame, m.tv.tstate)
		}
		if m.tv.frames != frame {
			t.Errorf("frame %d: TV frame count = %d", frame, m.tv.frames)
		}
		if p := m.pixel(0, 0); p != colRed {
			t.Errorf("frame %d: border = %#x, want red %#x", frame, p, colRed)
		}
		// constant beeper level : samples don't accumulate across frames
		buffer := m.beeper.Buffer()
		if s := buffer.GetSample(buffer.Size() / 2); s != amplBeeper {
			t.Errorf("frame %d: beeper sample = %d, want %d", frame, s, amplBeeper)
		}
	}
}
