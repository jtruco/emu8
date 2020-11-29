package memory

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// BankMap & Mapper
// -----------------------------------------------------------------------------

// BankMap contains a bank bus and mapping information
type BankMap struct {
	bus        device.BusDevice
	bank       *Bank
	address    uint16
	endaddress uint16
	active     bool
	init       bool
	write      bool
}

// Mapper memory bank mapper interface
type Mapper interface {
	// Init inits the mapper
	Init(memory *Memory)
	// SelectBank for Read access
	SelectBank(address uint16) (*BankMap, uint16)
	// SelectBank for Write access
	SelectBankWrite(address uint16) (*BankMap, uint16)
}

// NewBankMap creates a memory bank map
func NewBankMap(address uint16, size int, readonly bool, active bool) *BankMap {
	bmap := new(BankMap)
	bmap.bank = NewBank(size, readonly)
	bmap.bus = bmap.bank
	bmap.address = address
	bmap.endaddress = address + uint16(size) - 1
	bmap.active = active
	bmap.init = active
	bmap.write = !readonly
	return bmap
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
func NewBusMap(bus device.BusDevice, address uint16, size int, readonly bool, active bool) *BankMap {
	bmap := new(BankMap)
	bmap.bus = bus
	bmap.address = address
	bmap.endaddress = address + uint16(size) - 1
	bmap.active = active
	bmap.init = active
	return bmap
}

// Init inits bank
func (bmap *BankMap) Init() {
	bmap.Bus().Init()
	bmap.active = bmap.init
}

// Reset resets bank
func (bmap *BankMap) Reset() {
	bmap.Bus().Reset()
	bmap.active = bmap.init
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
func (bmap *BankMap) Bus() device.BusDevice {
	return bmap.bus
}

// SetActive is bank active
func (bmap *BankMap) SetActive(active bool) {
	bmap.active = active
}

// -----------------------------------------------------------------------------
// Default Mapper implementation
// -----------------------------------------------------------------------------

// DefaultMapper is a simple but inefficient memory mapper
type DefaultMapper struct {
	memory *Memory
}

// Init inits the mapper
func (mapper *DefaultMapper) Init(memory *Memory) {
	mapper.memory = memory
}

// SelectBank selects the first active bank mapped at address
func (mapper *DefaultMapper) SelectBank(address uint16) (*BankMap, uint16) {
	return mapper.selectInternal(address, false)
}

// SelectBankWrite selects the first active bank mapped at address for write
func (mapper *DefaultMapper) SelectBankWrite(address uint16) (*BankMap, uint16) {
	return mapper.selectInternal(address, true)
}

// selectInternal internal bank selection
func (mapper *DefaultMapper) selectInternal(address uint16, write bool) (*BankMap, uint16) {
	for _, bank := range mapper.memory.banks {
		if bank != nil && bank.active && (!write || bank.write == write) {
			if address >= bank.address && address <= bank.endaddress {
				return bank, address - bank.address
			}
		}
	}
	return nil, 0
}

// -----------------------------------------------------------------------------
// MaskMapper
// -----------------------------------------------------------------------------

// MaskMapper is a efficient memory mapper based on address bits (shift and mask)
type MaskMapper struct {
	memory *Memory
	Shift  uint
	Mask   uint16
}

// Init inits the mapper
func (mapper *MaskMapper) Init(memory *Memory) {
	mapper.memory = memory
}

// SelectBank selects read bank mapped at address
func (mapper *MaskMapper) SelectBank(address uint16) (*BankMap, uint16) {
	bank := int(address >> mapper.Shift)
	if bank < len(mapper.memory.banks) {
		return mapper.memory.banks[bank], address & mapper.Mask
	}
	return nil, 0
}

// SelectBankWrite selects write bank mapped at address
func (mapper *MaskMapper) SelectBankWrite(address uint16) (*BankMap, uint16) {
	return mapper.SelectBank(address)
}
