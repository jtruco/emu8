// Package z80 a Zilog Z80 CPU emulator
package z80

import (
	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/bus"
)

// -----------------------------------------------------------------------------
// Z80 - Zilog Z80 CPU
// -----------------------------------------------------------------------------

// Z80 the Zilog Z80 CPU
type Z80 struct {
	State                       // Z80 State
	clock    device.Clock       // Clock device
	mem      bus.Bus            // Memory data bus
	io       bus.Bus            // I/O data bus
	OnIntAck device.AckCallback // INT / NMI ack callback
}

// New creates a new Z80
func New(clock device.Clock, mem, io bus.Bus) *Z80 {
	z80 := new(Z80)
	z80.clock = clock
	z80.mem = mem
	z80.io = io
	z80.State.Init()
	return z80
}

// Clock gets the Cpu Clock
func (z80 *Z80) Clock() device.Clock {
	return z80.clock
}

// Memory gets the Cpu Memory bus
func (z80 *Z80) Memory() bus.Bus {
	return z80.mem
}

// IO gets the Cpu IO bus
func (z80 *Z80) IO() bus.Bus {
	return z80.io
}

// Init initializes Cpu (power-on)
func (z80 *Z80) Init() {
	z80.State.HardReset()
}

// Reset (soft) resets Cpu
func (z80 *Z80) Reset() {
	z80.State.SoftReset()
}

// Execute executes one instruction
func (z80 *Z80) Execute() int {
	tstate := z80.clock.Tstates()
	if z80.NmiRq {
		z80.NMInterrupt()
	} else if z80.IntRq && z80.IFF1 {
		z80.Interrupt()
	} else {
		z80.fetchAndExecute(z80.execute)
	}
	return z80.clock.Tstates() - tstate
}

// InterruptRequest request a maskable interrupt
func (z80 *Z80) InterruptRequest(request bool) {
	z80.IntRq = request
}

// Interrupt a maskable interrupt
func (z80 *Z80) Interrupt() {
	// Check maskable interrupts enabled
	if !z80.IFF1 {
		return
	}
	// Check EI activate
	for z80.ActiveEI {
		z80.fetchAndExecute(z80.execute)
	}
	// Check NMOS IFF2 parity bug
	if z80.ReadIFF2 {
		z80.F &= ^FlagP
	}
	// Accept interrupt
	z80.IntRq = z80.acceptInterrupt()
	// Process interrupt
	z80.IFF1, z80.IFF2 = false, false
	z80.incR()
	z80.clock.Add(7) // 7 tstate
	z80.push8(highbyte(z80.PC))
	z80.push8(lowbyte(z80.PC))
	// Case Interrupt Mode
	switch z80.IM {
	case 0, 1:
		// RST 0x38
		z80.PC = 0x0038
	case 2:
		tmp := (uint16(z80.I) << 8) | 0xff
		pcl := z80.readByte(tmp)
		tmp++
		pch := z80.readByte(tmp)
		z80.PC = toword(pcl, pch)
	}
	z80.Memptr.Set(z80.PC)
}

// NMInterruptRequest request a maskable interrupt
func (z80 *Z80) NMInterruptRequest(request bool) {
	z80.NmiRq = request
}

// NMInterrupt a not maskable interrupt
func (z80 *Z80) NMInterrupt() {
	// Accept interrupt
	z80.NmiRq = z80.acceptInterrupt()
	// Process interrupt
	z80.IFF1 = false
	z80.incR()
	z80.clock.Add(5)
	z80.push8(highbyte(z80.PC))
	z80.push8(lowbyte(z80.PC))
	// NMI : set PC address at 0x0066
	z80.PC = 0x0066
}

// fetchAndExecute fetchs and executes an opcode
func (z80 *Z80) fetchAndExecute(execute func(byte)) {
	opcode := z80.readByte(z80.PC)
	z80.clock.Inc() // +1 tstate opcode execution
	z80.incPC()
	z80.incR()
	z80.ActiveEI = false
	z80.ReadIFF2 = false
	execute(opcode)
}

// fetchAndExecute fetchs and executes an opcode
func (z80 *Z80) acceptInterrupt() bool {
	// Check halted state
	if z80.Halted {
		z80.incPC()
		z80.Halted = false
	}
	// Ack interrupt
	if z80.OnIntAck != nil {
		return z80.OnIntAck()
	}
	return false // don't ack interrupt
}
