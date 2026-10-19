package emulator

import (
	"strings"
	"sync"
	"testing"

	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// spectrumSNA is a 48K SNA snapshot with a marker byte at address
func spectrumSNA(address int, marker byte) []byte {
	data := make([]byte, 49179)
	data[27+address-0x4000] = marker
	return data
}

// cpcSNA is a 64K CPC SNA snapshot with a marker byte at address
func cpcSNA(address int, marker byte) []byte {
	data := make([]byte, 0x100+0x10000)
	copy(data, "MV - SNA")
	data[0x10], data[0x6b] = 1, 64
	data[0x100+address] = marker
	return data
}

// romFileSystem returns a memory file system with zeroed ROM files
func romFileSystem() *vfs.MemFileSystem {
	fs := vfs.NewMemFileSystem()
	for _, name := range []string{"zxspectrum.rom", "cpc464_os.rom", "cpc464_basic.rom"} {
		fs.SaveFile(vfs.NewFileInfoData(name, make([]byte, 0x4000)))
	}
	return fs
}

// Emulators are independent : each one has its own machine, settings and
// file system, and they can run at the same time.
func TestEmulatorInstancesAreIndependent(t *testing.T) {
	type instance struct {
		model     string
		frequency int
		snapshot  []byte
		marker    int // marker position in the saved snapshot data
		samples   int // audio samples per frame
	}
	instances := []instance{
		{"zx48k", 44100, spectrumSNA(0x8000, 0x5a), 27 + 0x4000, 44100 / 50},
		{"cpc464", 22050, cpcSNA(0x8000, 0xa5), 0x100 + 0x8000, 22050 / 50},
	}
	emus := make([]*Emulator, len(instances))
	for i, in := range instances {
		fs := romFileSystem()
		fs.SaveFile(vfs.NewFileInfoData("game.sna", in.snapshot)) // same name in each file system
		emu, err := Create(machine.Settings{ModelID: in.model, AudioFrequency: in.frequency}, fs)
		if err != nil {
			t.Fatal(err)
		}
		if err := emu.LoadFile("game.sna"); err != nil {
			t.Fatal(err)
		}
		emus[i] = emu
	}
	for i, in := range instances {
		emu := emus[i]
		if id := emu.Machine().Config().ModelID; id != in.model {
			t.Errorf("Config().ModelID = %q, want %q", id, in.model)
		}
		if emu.Machine().Config().Name == "" {
			t.Errorf("%s: Config().Name is empty", in.model)
		}
		if n := emu.Control().Audio().Device().Buffer().Size(); n != in.samples {
			t.Errorf("%s: %d audio samples per frame, want %d", in.model, n, in.samples)
		}
		data := emu.Machine().SaveState().Data
		if got, want := data[in.marker], in.snapshot[in.marker]; got != want {
			t.Errorf("%s: snapshot marker %#02x, want %#02x from its own file system", in.model, got, want)
		}
	}
	// run both emulators at the same time (go test -race checks shared state)
	var wg sync.WaitGroup
	for _, emu := range emus {
		wg.Add(1)
		go func(emu *Emulator) {
			defer wg.Done()
			emu.Start()
			for frame := 0; frame < 20; frame++ {
				emu.Emulate()
			}
			emu.Stop()
		}(emu)
	}
	wg.Wait()
}

// Create returns a ready to use emulator : the machine is initialized with
// its ROM loaded from the file system, and the emulation timings are set.
func TestCreateReadyEmulator(t *testing.T) {
	// ROM : LD A,0x42 ; LD (0x8000),A ; JR $
	rom := make([]byte, 0x4000)
	copy(rom, []byte{0x3e, 0x42, 0x32, 0x00, 0x80, 0x18, 0xfe})
	fs := vfs.NewMemFileSystem()
	fs.SaveFile(vfs.NewFileInfoData("zxspectrum.rom", rom))
	emu, err := Create(machine.Settings{ModelID: "zx48k"}, fs)
	if err != nil {
		t.Fatal(err)
	}
	emu.Start()
	emu.Emulate() // one frame
	emu.Stop()
	data := emu.Machine().SaveState().Data // SNA : RAM from offset 27 (0x4000)
	if got := data[27+0x8000-0x4000]; got != 0x42 {
		t.Errorf("RAM 0x8000 = %#02x, want 0x42 written by the ROM", got)
	}
}

// Create fails when the machine can't be initialized (missing ROM)
func TestCreateMissingROM(t *testing.T) {
	for _, model := range []string{"zx48k", "cpc464"} {
		emu, err := Create(machine.Settings{ModelID: model}, vfs.NewMemFileSystem())
		if emu != nil || err == nil {
			t.Errorf("%s: Create without ROMs = %v, %v; want nil and an error", model, emu, err)
		}
	}
}

// LoadFile returns an error, naming the file, when it can't be loaded
func TestLoadFileErrors(t *testing.T) {
	cases := []struct {
		model, name string
		data        []byte // nil : missing file
		reason      string // expected in the error
	}{
		{"zx48k", "game.sna", nil, "not found"},
		{"zx48k", "game.xyz", []byte{0}, "not supported format"},
		{"zx48k", "game.dsk", []byte{0}, "not supported format"},
		{"zx48k", "game.sna", make([]byte, 100), "invalid sna snapshot"},
		{"zx48k", "game.z80", make([]byte, 10), "invalid z80 snapshot"},
		{"zx48k", "game.tap", []byte{0x10, 0x00, 0x00}, "invalid tap tape"},
		{"zx48k", "game.tzx", []byte("ZXTape?"), "invalid tzx tape"},
		{"cpc464", "game.sna", make([]byte, 100), "invalid sna snapshot"},
		{"cpc464", "game.cdt", []byte("ZXTape?"), "invalid cdt tape"},
		{"cpc464", "game.z80", []byte{0}, "not supported format"},
	}
	for _, c := range cases {
		fs := romFileSystem()
		if c.data != nil {
			fs.SaveFile(vfs.NewFileInfoData(c.name, c.data))
		}
		emu, err := Create(machine.Settings{ModelID: c.model}, fs)
		if err != nil {
			t.Fatal(err)
		}
		err = emu.LoadFile(c.name)
		if err == nil || !strings.Contains(err.Error(), c.name) || !strings.Contains(err.Error(), c.reason) {
			t.Errorf("%s: LoadFile(%s) = %v, want an error naming the file and %q", c.model, c.name, err, c.reason)
		}
	}
}
