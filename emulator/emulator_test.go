package emulator

import (
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
		fs := vfs.NewMemFileSystem()
		fs.SaveFile(vfs.NewFileInfoData("game.sna", in.snapshot)) // same name in each file system
		emu, err := Create(machine.Settings{ModelID: in.model, AudioFrequency: in.frequency}, fs)
		if err != nil {
			t.Fatal(err)
		}
		emu.Init()
		emu.LoadFile("game.sna")
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
