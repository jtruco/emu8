package z80

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/jtruco/emu8/cpu"
	"github.com/jtruco/emu8/cpu/z80"
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

// TestZexall test Z80 emulator
func TestZ80(t *testing.T) {
	testFile("zexall.bin")
}

func testFile(testfile string) {

	// initialize cpu
	mem, io := &Memory{}, &Memory{}
	cpu := z80.New(cpu.NewClock(), mem, io)

	// load testfile
	data, err := ioutil.ReadFile(testfile)
	if err != nil {
		log.Println(err.Error())
		return
	}
	copy(mem.data[0x100:], data[:])

	// prepare test
	mem.data[0] = 0xc3 // JP 0x100 CP/M TPA
	mem.data[1] = 0x00
	mem.data[2] = 0x01
	mem.data[5] = 0xc9 // RET from BDOS call

	// run test
	done := false
	for !done {
		cpu.Execute()
		// Emulate CP/M syscall at address 5
		if cpu.PC == 0x05 {
			switch cpu.C {
			case 0: // BDOS 0 System Reset
				fmt.Println("Z80 reset after ", cpu.Clock().Tstates(), " t-states")
				done = true
			case 2: // BDOS 2 console char output
				fmt.Printf("%c", cpu.E)
			case 9: // BDOS 9 console string output (string terminated by "$")
				addr := cpu.DE.Get()
				for mem.data[addr] != '$' {
					fmt.Printf("%c", mem.data[addr])
					addr++
				}
			default:
				fmt.Println("BDOS Call ", cpu.C)
				done = true
			}
		}
	}
}
