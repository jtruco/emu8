// Package memory defines memory components
package memory

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

// Defaults
const (
	DataDefault byte = 0
)

// -----------------------------------------------------------------------------
// Memory device bus
// -----------------------------------------------------------------------------

// Memory is a memory structure of banks or bus devices
type Memory struct {
	size   int
	banks  []*BankMap
	mapper Mapper
}

// New creates a new memory device with a size and a number banks
func New(size int, banks int) *Memory {
	memory := new(Memory)
	memory.size = size
	memory.banks = make([]*BankMap, banks)
	memory.SetMapper(new(DefaultMapper))
	return memory
}

// Bank mapping

// Banks returns Banks
func (memory *Memory) Banks() []*BankMap {
	return memory.banks
}

// Bank returns bank mapped at index
func (memory *Memory) Bank(index int) *Bank {
	return memory.Map(index).Bank()
}

// Map returns bank map at index
func (memory *Memory) Map(index int) *BankMap {
	return memory.banks[index]
}

// SetMap sets the bank map at index
func (memory *Memory) SetMap(index int, bank *BankMap) {
	memory.banks[index] = bank
}

// Switch switches two memory banks and update its active state
func (memory *Memory) Switch(current, new int) {
	curbank, newbank := memory.banks[current], memory.banks[new]
	memory.banks[current], memory.banks[new] = newbank, curbank
	curbank.active, newbank.active = false, true
}

// Mapper returns the bank mapper
func (memory *Memory) Mapper() Mapper { return memory.mapper }

// SetMapper sets the bank mapper
func (memory *Memory) SetMapper(mapper Mapper) {
	mapper.Init(memory)
	memory.mapper = mapper
}

// LoadRAM loads data into memory starting at address
func (memory *Memory) LoadRAM(address uint16, data []byte) {
	length := len(data)
	offset, last := 0, 0
	for offset < length {
		bankaddr := address + uint16(offset)
		bank, rel := memory.mapper.SelectBankWrite(bankaddr)
		last = offset + bank.bank.size
		bank.bank.Load(rel, data[offset:last])
		offset = last
	}
}

// Device interface

// Init initializes the memory
func (memory *Memory) Init() {
	for _, b := range memory.banks {
		if b != nil {
			b.Init()
		}
	}
}

// Reset resets the memory data at initial state
func (memory *Memory) Reset() {
	for _, b := range memory.banks {
		b.Reset()
	}
}

// Bus interface

// Read reads a byte from memory
func (memory *Memory) Read(address uint16) byte {
	bank, bankAddr := memory.mapper.SelectBank(address)
	if bank != nil {
		return bank.bus.Read(bankAddr)
	}
	return DataDefault
}

// Write writes a byte to the memory
func (memory *Memory) Write(address uint16, data byte) {
	bank, bankAddr := memory.mapper.SelectBankWrite(address)
	if bank != nil {
		bank.bus.Write(bankAddr, data)
	}
	// default : no write
}
