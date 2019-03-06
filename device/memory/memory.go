// Package memory defines memory components
package memory

// Common memory sizes
const (
	Size128B = 0x0080
	Size256B = 0x0100
	Size512B = 0x0200
	Size1K   = 0x0400
	Size2K   = 0x0800
	Size4K   = 0x1000
	Size8K   = 0x2000
	Size16K  = 0x4000
	Size32K  = 0x8000
	Size64K  = 0x10000
)

// Defaults
const (
	DataDefault = byte(0)
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
	memory := &Memory{}
	memory.size = size
	memory.banks = make([]*BankMap, banks)
	memory.mapper = &DefaultMapper{}
	return memory
}

// Banks returns Banks
func (memory *Memory) Banks() []*BankMap {
	return memory.banks
}

// GetBankMap returns bank map at index
func (memory *Memory) GetBankMap(index int) *BankMap {
	return memory.banks[index]
}

// Load data to memory
func (memory *Memory) Load(address uint16, data []byte) {
	length := uint16(len(data))
	offset := uint16(0)
	last := uint16(0)
	for offset < length {
		info, rel := memory.mapper.SelectBank(memory, address+offset)
		last = offset + uint16(info.bank.size)
		info.bank.Load(rel, data[offset:last])
		offset = last
	}
}

// SetMap sets the bank map at index
func (memory *Memory) SetMap(index int, bank *BankMap) {
	memory.banks[index] = bank
}

// SetMapper sets the bank mapper
func (memory *Memory) SetMapper(mapper Mapper) {
	memory.mapper = mapper
}

// Switch switchs two memory banks and update its active state
func (memory *Memory) Switch(current, new int) {
	curbank, newbank := memory.banks[current], memory.banks[new]
	memory.banks[current], memory.banks[new] = newbank, curbank
	curbank.active, newbank.active = false, true
}

// Device interface

// Init initializes the memory
func (memory *Memory) Init() {
	for _, b := range memory.banks {
		if b != nil {
			b.active = b.init
			b.bus.Init()
		}
	}
}

// Reset resets the memory data at initial state
func (memory *Memory) Reset() {
	for _, b := range memory.banks {
		b.active = b.init
		b.bus.Reset()
	}
}

// Bus interface

// Access access memory at address
func (memory *Memory) Access(address uint16) {
	bank, bankAddr := memory.mapper.SelectBank(memory, address)
	if bank != nil {
		bank.bus.Access(bankAddr)
	}
}

// Read reads a byte from memory
func (memory *Memory) Read(address uint16) byte {
	bank, bankAddr := memory.mapper.SelectBank(memory, address)
	if bank != nil {
		return bank.bus.Read(bankAddr)
	}
	return DataDefault
}

// Write writes a byte to the memory
func (memory *Memory) Write(address uint16, data byte) {
	bank, bankAddr := memory.mapper.SelectBank(memory, address)
	if bank != nil {
		bank.bus.Write(bankAddr, data)
	}
	// default : no write
}
