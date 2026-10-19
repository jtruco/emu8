package z80

import (
	"os"
	"testing"

	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/cpu/z80"
)

// BenchmarkZ80Zexall executes zexall instructions (ns per instruction)
func BenchmarkZ80Zexall(b *testing.B) {
	data, err := os.ReadFile("zexall.bin")
	if err != nil {
		b.Fatal(err)
	}
	mem, io := new(Memory), new(Memory)
	copy(mem.data[0x100:], data)
	mem.data[0], mem.data[1], mem.data[2] = 0xc3, 0x00, 0x01 // JP 0x100
	mem.data[5] = 0xc9                                       // RET from BDOS call
	cpu := z80.New(device.NewClock(), mem, io)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpu.Execute()
	}
}
