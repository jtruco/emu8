package spectrum

import (
	"strings"
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// Init fails, naming the ROM, when the ROM file is missing or too short
func TestSpectrumInitROMErrors(t *testing.T) {
	cases := []struct {
		name string
		rom  []byte
	}{
		{"missing", nil},
		{"short", make([]byte, 100)},
	}
	for _, c := range cases {
		fs := vfs.NewMemFileSystem()
		if c.rom != nil {
			fs.SaveFile(vfs.NewFileInfoData(zxRomName, c.rom))
		}
		m := New(ZXSpectrum48K, machine.Settings{}).(*Spectrum)
		controller.New(m, fs)
		if err := m.Init(); err == nil || !strings.Contains(err.Error(), zxRomName) {
			t.Errorf("%s ROM: Init() = %v, want an error naming %s", c.name, err, zxRomName)
		}
	}
}

// Reset clears the RAM and keeps the ROM loaded by Init
func TestSpectrumResetKeepsROM(t *testing.T) {
	fs := vfs.NewMemFileSystem()
	rom := make([]byte, 0x4000)
	rom[0] = 0xf3
	fs.SaveFile(vfs.NewFileInfoData(zxRomName, rom))
	m := New(ZXSpectrum48K, machine.Settings{}).(*Spectrum)
	controller.New(m, fs)
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	m.memory.Bank(2).Data()[0] = 0x55 // RAM at 8000h
	m.Reset()
	if got := m.memory.Bank(0).Data()[0]; got != 0xf3 {
		t.Errorf("ROM after Reset = %#02x, want 0xf3", got)
	}
	if got := m.memory.Bank(2).Data()[0]; got != 0 {
		t.Errorf("RAM after Reset = %#02x, want 0", got)
	}
}
