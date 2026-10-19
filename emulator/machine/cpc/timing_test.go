package cpc

import "testing"

// Instruction timings of the CPC (Z80 with Gate Array wait states), from the
// CPC instruction timing tables (NOP = 4 Ts). Known differences are marked
// with the bug or design flaw that causes them : the test logs them, and fails
// when they are fixed (to remove the mark).
func TestCpcInstructionTimings(t *testing.T) {
	cases := []struct {
		name  string
		code  []byte
		nops  int    // real CPC timing in NOPs
		known string // known difference
	}{
		{"NOP", []byte{0x00}, 1, ""},
		{"LD (HL),A", []byte{0x77}, 2, ""},
		{"LD A,(HL)", []byte{0x7e}, 2, ""},
		{"POP BC", []byte{0xc1}, 3, ""},
		{"PUSH BC", []byte{0xc5}, 4, "D7"},
		{"OUT (n),A", []byte{0xd3, 0x10}, 3, ""},
		{"IN A,(n)", []byte{0xdb, 0x10}, 3, ""},
		{"OUT (C),A", []byte{0xed, 0x79}, 4, "B17"},
		{"IN A,(C)", []byte{0xed, 0x78}, 4, "B17"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestCPC(t)
			m.memory.Bank(2).Load(0, c.code) // RAM at 4000h
			m.cpu.PC, m.cpu.SP = 0x4000, 0xbf00
			m.cpu.B, m.cpu.C = 0xf7, 0x00 // port F700h : PPI control
			m.cpu.HL.Set(0x8000)
			m.clock.SetTstates(0)
			m.Emulate()
			got, want := m.clock.Tstates(), c.nops*4
			switch {
			case got == want && c.known != "":
				t.Errorf("%s = %d Ts as on the CPC : %s fixed, remove the known mark", c.name, got, c.known)
			case got != want && c.known != "":
				t.Logf("%s = %d Ts, CPC %d Ts (known : %s)", c.name, got, want, c.known)
			case got != want:
				t.Errorf("%s = %d Ts, want %d (CPC %d NOPs)", c.name, got, want, c.nops)
			}
		})
	}
}
