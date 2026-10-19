package cpc

import "testing"

// The lower and upper ROMs are switched over the RAM by the gate array;
// writes always go to RAM.
func TestCpcROMSwitching(t *testing.T) {
	m := newTestCPC(t)
	m.memory.Bank(cpcLowerROM).Data()[0] = 0xa1 // lower ROM at 0000h
	m.memory.Bank(cpcUpperROM).Data()[0] = 0xa2 // upper ROM at C000h
	m.memory.Write(0x0000, 0x11)                // RAM under the ROMs
	m.memory.Write(0xc000, 0x22)
	cases := []struct {
		name         string
		config       byte // gate array mode and ROM configuration
		lower, upper byte
	}{
		{"both ROMs", 0x80, 0xa1, 0xa2},
		{"lower ROM off", 0x84, 0x11, 0xa2},
		{"upper ROM off", 0x88, 0xa1, 0x22},
		{"both ROMs off", 0x8c, 0x11, 0x22},
	}
	for _, c := range cases {
		m.gatearray.Write(c.config)
		if lower, upper := m.memory.Read(0x0000), m.memory.Read(0xc000); lower != c.lower || upper != c.upper {
			t.Errorf("%s: read %#02x/%#02x, want %#02x/%#02x", c.name, lower, upper, c.lower, c.upper)
		}
	}
	if m.memory.Bank(1).Data()[0] != 0x11 || m.memory.Bank(5).Data()[0] != 0x22 {
		t.Error("writes did not go to RAM")
	}
}
