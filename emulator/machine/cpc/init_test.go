package cpc

import (
	"strings"
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// Init fails, naming the ROM, when a ROM file is missing or too short
func TestCpcInitROMErrors(t *testing.T) {
	cases := []struct {
		name    string
		os      []byte
		basic   []byte
		options string
		romname string // ROM named by the error
	}{
		{"missing OS", nil, make([]byte, 0x4000), "", cpcOsRomName},
		{"missing ES OS", make([]byte, 0x4000), make([]byte, 0x4000), "es", cpcOsRomNameES},
		{"missing BASIC", make([]byte, 0x4000), nil, "", cpcBasicRomName},
		{"short OS", make([]byte, 0x3fff), make([]byte, 0x4000), "", cpcOsRomName},
		{"short BASIC", make([]byte, 0x4000), make([]byte, 10), "", cpcBasicRomName},
	}
	for _, c := range cases {
		fs := vfs.NewMemFileSystem()
		if c.os != nil {
			fs.SaveFile(vfs.NewFileInfoData(cpcOsRomName, c.os))
		}
		if c.basic != nil {
			fs.SaveFile(vfs.NewFileInfoData(cpcBasicRomName, c.basic))
		}
		m := New(AmstradCPC464, machine.Settings{Options: c.options}).(*AmstradCPC)
		controller.New(m, fs)
		if err := m.Init(); err == nil || !strings.Contains(err.Error(), c.romname) {
			t.Errorf("%s: Init() = %v, want an error naming %s", c.name, err, c.romname)
		}
	}
}

// Reset clears the RAM and keeps the ROMs loaded by Init and the PPI jumpers
func TestCpcResetKeepsROMs(t *testing.T) {
	m := newTestCPC(t)
	m.memory.Bank(cpcLowerROM).Data()[0] = 0xa1
	m.memory.Bank(cpcUpperROM).Data()[0] = 0xa2
	m.memory.Bank(3).Data()[0] = 0x55 // RAM at 8000h
	m.Reset()
	if lo, up := m.memory.Bank(cpcLowerROM).Data()[0], m.memory.Bank(cpcUpperROM).Data()[0]; lo != 0xa1 || up != 0xa2 {
		t.Errorf("ROMs after Reset = %#02x, %#02x; want 0xa1, 0xa2", lo, up)
	}
	if got := m.memory.Bank(3).Data()[0]; got != 0 {
		t.Errorf("RAM after Reset = %#02x, want 0", got)
	}
	if m.ppi.jumpers != cpcJumpers {
		t.Errorf("PPI jumpers after Reset = %#02x, want %#02x", m.ppi.jumpers, cpcJumpers)
	}
}
