package spectrum

import "testing"

// Instruction timings through the CPU (code in uncontended RAM at 8000h)
func TestSpectrumInstructionTimings(t *testing.T) {
	cases := []struct {
		name   string
		code   []byte
		hl     uint16
		tstate int // start Tstate
		want   int // Z80 Ts (plus contention)
	}{
		{"NOP", []byte{0x00}, 0x9000, 1000, 4},
		{"LD (HL),A", []byte{0x77}, 0x9000, 1000, 7},
		{"LD A,(HL)", []byte{0x7e}, 0x9000, 1000, 7},
		{"PUSH BC", []byte{0xc5}, 0x9000, 1000, 11},
		{"POP BC", []byte{0xc1}, 0x9000, 1000, 10},
		{"OUT (n),A", []byte{0xd3, 0xfe}, 0x9000, 1000, 11},
		{"IN A,(n)", []byte{0xdb, 0xfe}, 0x9000, 1000, 11},
		{"OUT (C),A", []byte{0xed, 0x79}, 0x9000, 1000, 12},
		{"IN A,(C)", []byte{0xed, 0x78}, 0x9000, 1000, 12},
		{"OUTI", []byte{0xed, 0xa3}, 0x9000, 1000, 16},
		{"INI", []byte{0xed, 0xa2}, 0x9000, 1000, 16},
		// contended write : fetch 14335..14339, write delay 2 at 14339
		{"LD (HL),A contended", []byte{0x77}, 0x4000, 14335, 7 + 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.memory.Bank(2).Load(0, c.code)
			m.cpu.PC, m.cpu.SP = 0x8000, 0xff00
			m.cpu.B, m.cpu.C = 0x80, 0xfe // port 80FEh : ULA, uncontended high byte
			m.cpu.HL.Set(c.hl)
			m.at(c.tstate)
			if got := m.cpu.Execute(); got != c.want {
				t.Errorf("%s = %d Ts, want %d", c.name, got, c.want)
			}
		})
	}
}
