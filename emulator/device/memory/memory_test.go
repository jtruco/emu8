package memory

import (
	"bytes"
	"testing"

	"github.com/jtruco/emu8/emulator/device/bus"
)

// newTestMemory creates a 64k memory : 16k ROM + 48k RAM, or 32k RAM only
// (16k ROM + 16k RAM, unmapped from 0x8000) when small.
func newTestMemory(small bool) *Memory {
	banks := 4
	if small {
		banks = 2
	}
	m := New(banks)
	m.SetMap(0, NewROM(0x0000, Size16K))
	for i := 1; i < banks; i++ {
		m.SetMap(i, NewRAM(uint16(i)*Size16K, Size16K))
	}
	m.SetMapper(bus.NewMaskMapper(14))
	m.Init()
	return m
}

func TestLoadRAM(t *testing.T) {
	seq := func(n int) []byte {
		d := make([]byte, n)
		for i := range d {
			d[i] = byte(i%251 + 1)
		}
		return d
	}
	cases := []struct {
		name    string
		small   bool
		address uint16
		data    []byte
		loaded  int // bytes expected in memory from address
	}{
		{"partial bank", false, 0x8000, seq(3), 3},
		{"unaligned in bank", false, 0x8123, seq(100), 100},
		{"across banks", false, 0x7ff0, seq(64), 64},
		{"whole RAM", false, 0x4000, seq(0xc000), 0xc000},
		{"past end of memory", false, 0xfff0, seq(64), 16},
		{"past mapped memory", true, 0x7ff0, seq(64), 16},
		{"empty", false, 0x4000, nil, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestMemory(c.small)
			m.LoadRAM(c.address, c.data)
			got := make([]byte, c.loaded)
			for i := range got {
				got[i] = m.Read(c.address + uint16(i))
			}
			if !bytes.Equal(got, c.data[:c.loaded]) {
				t.Errorf("memory from %#04x does not match the loaded data", c.address)
			}
			if c.loaded < len(c.data) && m.Read(0x0000) != 0 {
				t.Errorf("data past the end wrapped around to 0x0000")
			}
		})
	}
}

func TestLoadRAMSkipsROM(t *testing.T) {
	m := newTestMemory(false)
	m.LoadRAM(0x3ffe, []byte{1, 2, 3, 4})
	if m.Read(0x3ffe) != 0 || m.Read(0x3fff) != 0 {
		t.Error("LoadRAM wrote into ROM")
	}
}
