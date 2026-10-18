package spectrum

import (
	"testing"

	"github.com/jtruco/emu8/emulator/machine"
	"github.com/jtruco/emu8/emulator/machine/spectrum/format"
)

// A saved SNA snapshot resumes at the same PC and SP once loaded (the 48K
// ROM RETN at 0x0072 pops the PC pushed by SaveSNA).
func TestSpectrumSNASaveLoadResume(t *testing.T) {
	m := newTestSpectrum(t)
	m.cpu.PC = 0x8000
	m.cpu.SP = 0xff00
	m.cpu.A = 0x42
	m.memory.Write(0xc000, 0x99)
	state := m.SaveState()
	if state.Format != format.SNA {
		t.Fatalf("snapshot format = %q, want %q", state.Format, format.SNA)
	}

	m2 := newTestSpectrum(t)
	m2.memory.Bank(0).Load(0x0072, []byte{0xed, 0x45}) // ROM : RETN
	m2.LoadState(machine.State{Format: state.Format, Data: state.Data})
	m2.cpu.Execute()

	if m2.cpu.PC != 0x8000 || m2.cpu.SP != 0xff00 {
		t.Errorf("resumed at PC/SP = %#04x/%#04x, want 0x8000/0xff00", m2.cpu.PC, m2.cpu.SP)
	}
	if m2.cpu.A != 0x42 {
		t.Errorf("A = %#02x, want 0x42", m2.cpu.A)
	}
	if v := m2.memory.Read(0xc000); v != 0x99 {
		t.Errorf("memory 0xc000 = %#02x, want 0x99", v)
	}
}
