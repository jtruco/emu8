package cpc

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// newBenchCpc builds a CPC 464 with the ROMs of the repository
func newBenchCpc(b *testing.B) *AmstradCPC {
	b.Helper()
	m := New(AmstradCPC464, machine.Settings{}).(*AmstradCPC)
	controller.New(m, vfs.NewDesktopFileSystem("../../.."))
	m.Init()
	if m.memory.Bank(cpcLowerROM).Data()[0] == 0 {
		b.Skip("ROM not found")
	}
	return m
}

// runCpcFrame runs one frame (headless)
func runCpcFrame(m *AmstradCPC) {
	m.BeginFrame()
	m.clock.Restart(cpcTStates)
	for m.clock.Tstates() < cpcTStates {
		m.Emulate()
	}
	m.EndFrame()
}

// benchCpcFrames runs b.N frames after warm up frames
func benchCpcFrames(b *testing.B, m *AmstradCPC, warmup int) {
	b.Helper()
	for i := 0; i < warmup; i++ {
		runCpcFrame(m)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runCpcFrame(m)
	}
}

// BenchmarkCpcFrameROM : BASIC idle loop of the CPC 464 ROM
func BenchmarkCpcFrameROM(b *testing.B) {
	benchCpcFrames(b, newBenchCpc(b), 200)
}

// BenchmarkCpcFrameScreenWrite : screen memory writes (LDIR to C000h) and
// gate array I/O (OUT (C),C to the border pen)
func BenchmarkCpcFrameScreenWrite(b *testing.B) {
	m := newBenchCpc(b)
	for i := 0; i < 200; i++ {
		runCpcFrame(m)
	}
	// DI ; loop: LD HL,0C000h ; LD DE,0C001h ; LD BC,3FFFh ; LD (HL),A ; LDIR ;
	//            LD BC,7F10h ; OUT (C),C ; INC A ; JR loop
	m.memory.Bank(2).Load(0, []byte{0xf3, 0x21, 0x00, 0xc0, 0x11, 0x01, 0xc0, 0x01, 0xff, 0x3f,
		0x77, 0xed, 0xb0, 0x01, 0x10, 0x7f, 0xed, 0x49, 0x3c, 0x18, 0xec})
	m.cpu.PC = 0x4000
	benchCpcFrames(b, m, 0)
}
