package cpc

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// The OS ROM is selected by the machine options
func TestCpcOptionsSelectOSRom(t *testing.T) {
	roms := map[string]byte{cpcOsRomName: 0x01, cpcOsRomNameES: 0x02, cpcOsRomNameFR: 0x03}
	cases := []struct {
		options string
		marker  byte
	}{
		{"", 0x01}, {"es", 0x02}, {"fr", 0x03}, {"xx", 0x01},
	}
	for _, c := range cases {
		fs := vfs.NewMemFileSystem()
		vfs.SetFileSystem(fs)
		for name, marker := range roms {
			rom := make([]byte, 0x4000)
			rom[0] = marker
			fs.SaveFile(vfs.NewFileInfoData(name, rom))
		}
		fs.SaveFile(vfs.NewFileInfoData(cpcBasicRomName, make([]byte, 0x4000)))
		m := New(AmstradCPC464, machine.Settings{Options: c.options}).(*AmstradCPC)
		controller.New(m)
		m.Init()
		if got := m.memory.Bank(cpcLowerROM).Data()[0]; got != c.marker {
			t.Errorf("options %q: OS ROM marker %#02x, want %#02x", c.options, got, c.marker)
		}
	}
}
