package z80

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/cpu/z80"
)

// Memory

type Memory struct {
	data [0x10000]byte
}

func (m *Memory) Read(address uint16) byte {
	return m.data[address]
}

func (m *Memory) Write(address uint16, value byte) {
	m.data[address] = value
}

// Z80 Test

// TestZ80 runs the zexall Z80 instruction set exerciser
func TestZ80(t *testing.T) {
	if testing.Short() {
		t.Skip("zexall: slow Z80 instruction set test (~90 s)")
	}
	output := runCPM(t, "zexall.bin")
	if !strings.Contains(output, "Tests complete") {
		t.Fatal("zexall did not complete")
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "ERROR") {
			t.Error(strings.TrimSpace(line))
		}
	}
}

// runCPM runs a CP/M program and returns its console output. With -v the
// output is also shown while it runs.
func runCPM(t *testing.T, testfile string) string {
	t.Helper()

	// initialize cpu
	mem, io := new(Memory), new(Memory)
	cpu := z80.New(device.NewClock(), mem, io)

	// load testfile
	data, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	copy(mem.data[0x100:], data[:])

	// prepare test
	mem.data[0] = 0xc3 // JP 0x100 CP/M TPA
	mem.data[1] = 0x00
	mem.data[2] = 0x01
	mem.data[5] = 0xc9 // RET from BDOS call

	// console output
	var output strings.Builder
	print := func(c byte) {
		output.WriteByte(c)
		if testing.Verbose() {
			fmt.Printf("%c", c)
		}
	}

	// run test
	for {
		cpu.Execute()
		// Emulate CP/M syscall at address 5
		if cpu.PC == 0x05 {
			switch cpu.C {
			case 0: // BDOS 0 System Reset
				t.Logf("Z80 reset after %d t-states", cpu.Clock().Tstates())
				return output.String()
			case 2: // BDOS 2 console char output
				print(cpu.E)
			case 9: // BDOS 9 console string output (string terminated by "$")
				for addr := cpu.DE.Get(); mem.data[addr] != '$'; addr++ {
					print(mem.data[addr])
				}
			default:
				t.Fatalf("unsupported BDOS call %d", cpu.C)
			}
		}
	}
}
