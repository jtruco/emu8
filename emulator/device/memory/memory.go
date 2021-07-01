// Package memory defines memory components
package memory

import "github.com/jtruco/emu8/emulator/device/bus"

// Common memory sizes
const (
	Size00K  = 0x0000
	Size128B = 0x0080
	Size256B = 0x0100
	Size512B = 0x0200
	Size1K   = 0x0400
	Size2K   = 0x0800
	Size4K   = 0x1000
	Size8K   = 0x2000
	Size16K  = 0x4000
	Size32K  = 0x8000
	Size48K  = 0xC000
	Size64K  = 0x10000
)

// -----------------------------------------------------------------------------
// Memory banks mapping
// -----------------------------------------------------------------------------

// NewRAM creates a RAM bank mapping
func NewRAM(address uint16, size uint16) *bus.Map {
	return NewMemoryBank(address, size, true, false)
}

// NewROM creates a ROM bank mapping
func NewROM(address uint16, size uint16) *bus.Map {
	return NewMemoryBank(address, size, true, true)
}

// NewMemoryBank creates a new memory bank mapping
func NewMemoryBank(address, size uint16, active, readonly bool) *bus.Map {
	bank := NewBank(size, readonly)
	bankmap := bus.NewMap(bank, address, size, active, readonly)
	return bankmap
}

// -----------------------------------------------------------------------------
// Memory device bus
// -----------------------------------------------------------------------------

// Memory is a composite bus device of mapped banks
type Memory struct {
	bus.Composite
}

// New creates a new memory device with a size and a number banks
func New(banks int) *Memory {
	memory := new(Memory)
	memory.Build(banks)
	return memory
}

// Bank returns memory bank mapped at index
func (memory *Memory) Bank(index int) *Bank {
	return memory.Map(index).Device().(*Bank)
}

// LoadRAM loads a data chunk of bytes into memory starting at address
func (memory *Memory) LoadRAM(address uint16, data []byte) {
	var offset, last uint
	length := uint(len(data))
	for offset < length {
		bankaddr := address + uint16(offset)
		m, rel := memory.Mapper().SelectWrite(bankaddr)
		last = offset + uint(m.Size())
		bank := m.Device().(*Bank)
		if bank != nil {
			bank.Load(rel, data[offset:last])
		}
		offset = last
	}
}
