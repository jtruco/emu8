package spectrum

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// newBenchSpectrum builds a 48K Spectrum with the ROM of the repository
func newBenchSpectrum(b *testing.B) *Spectrum {
	b.Helper()
	m := New(ZXSpectrum48K, machine.Settings{}).(*Spectrum)
	controller.New(m, vfs.NewDesktopFileSystem("../../.."))
	if err := m.Init(); err != nil {
		b.Skip(err)
	}
	return m
}

// benchFrames runs b.N frames (headless) after warm up frames
func benchFrames(b *testing.B, m *Spectrum, warmup int) {
	b.Helper()
	for i := 0; i < warmup; i++ {
		runHeadlessFrame(m)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runHeadlessFrame(m)
	}
}

// BenchmarkSpectrumFrameROM : BASIC idle loop of the 48K ROM
func BenchmarkSpectrumFrameROM(b *testing.B) {
	benchFrames(b, newBenchSpectrum(b), 100)
}

// BenchmarkSpectrumFrameScreenWrite : contended screen memory writes (LDIR)
func BenchmarkSpectrumFrameScreenWrite(b *testing.B) {
	m := newBenchSpectrum(b)
	for i := 0; i < 100; i++ {
		runHeadlessFrame(m)
	}
	// DI ; loop: LD HL,4000h ; LD DE,4001h ; LD BC,1AFFh ; LD (HL),A ; LDIR ; INC A ; JR loop
	m.memory.Bank(2).Load(0, []byte{0xf3, 0x21, 0x00, 0x40, 0x11, 0x01, 0x40, 0x01, 0xff, 0x1a,
		0x77, 0xed, 0xb0, 0x3c, 0x18, 0xf1})
	m.cpu.PC = 0x8000
	benchFrames(b, m, 0)
}
