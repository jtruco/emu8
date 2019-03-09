package memory

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// BankMap mapping information
// -----------------------------------------------------------------------------

// BankMap contains a bank bus and mapping information
type BankMap struct {
	bus        device.Bus
	bank       *Bank
	address    uint16
	endaddress uint16
	active     bool
	init       bool
}

// NewBankMap creates a memory bank map
func NewBankMap(address uint16, size int, readonly bool, active bool) *BankMap {
	bmap := BankMap{}
	bmap.bank = NewBank(size, readonly)
	bmap.bus = bmap.bank
	bmap.address = address
	bmap.endaddress = address + uint16(size) - 1
	bmap.active = active
	bmap.init = active
	return &bmap
}

// NewROM creates a ROM bank map
func NewROM(address uint16, size int) *BankMap {
	return NewBankMap(address, size, true, true)
}

// NewRAM creates a RAM bank map
func NewRAM(address uint16, size int) *BankMap {
	return NewBankMap(address, size, false, true)
}

// NewBusMap creates a bank map from a device bus
func NewBusMap(bus device.Bus, address uint16, size int, readonly bool, active bool) *BankMap {
	bmap := BankMap{}
	bmap.bus = bus
	bmap.address = address
	bmap.endaddress = address + uint16(size) - 1
	bmap.active = active
	bmap.init = active
	return &bmap
}

// Active is bank active
func (bmap *BankMap) Active() bool {
	return bmap.active
}

// Bank gets the bank
func (bmap *BankMap) Bank() *Bank {
	return bmap.bank
}

// Bus gets the device bus of the bank
func (bmap *BankMap) Bus() device.Bus {
	return bmap.bus
}

// SetActive is bank active
func (bmap *BankMap) SetActive(active bool) {
	bmap.active = active
}

// -----------------------------------------------------------------------------
// Mapper memory bank mapper
// -----------------------------------------------------------------------------

// Mapper memory bank mapper interface
type Mapper interface {
	SelectBank(m *Memory, address uint16) (*BankMap, uint16)
}

// DefaultMapper is a simple but inefficent memory mapper
type DefaultMapper struct{}

// SelectBank selects the first active bank mapped at address
func (mapper *DefaultMapper) SelectBank(memory *Memory, address uint16) (*BankMap, uint16) {
	for _, bank := range memory.banks {
		if bank != nil && bank.active {
			if address >= bank.address && address <= bank.endaddress {
				return bank, address - bank.address
			}
		}
	}
	return nil, 0
}

// BusMapper is a efficent memory mapper based on address bits (shift and mask)
type BusMapper struct {
	Shift uint
	Mask  uint16
}

// SelectBank selects bank mapped by high bits of bus address
func (mapper *BusMapper) SelectBank(memory *Memory, address uint16) (*BankMap, uint16) {
	return memory.banks[address>>mapper.Shift], address & mapper.Mask
}
