package memory

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Bank map
// -----------------------------------------------------------------------------

// BankMap contains the bank mapping information
type BankMap struct {
	bus        device.BusDevice // The bus device
	bank       *Bank            // The memory bank
	address    uint16           // Base memory address of the map
	endaddress uint16           // End memory address of the map
	readonly   bool             // Is readonly access
	active     bool             // Is bank active
	init       bool             // Initial active state
}

// NewMap creates a bank mapping from a bus device
func NewMap(bus device.BusDevice, address uint16, size int, readonly bool, active bool) *BankMap {
	bmap := new(BankMap)
	bmap.bus = bus
	bmap.address = address
	bmap.endaddress = address + uint16(size) - 1
	bmap.readonly = readonly
	bmap.active = active
	bmap.init = active
	return bmap
}

// NewBankMap creates a new memory bank mapping
func NewBankMap(address uint16, size int, readonly bool, active bool) *BankMap {
	bank := NewBank(size, readonly)
	bmap := NewMap(bank, address, size, readonly, active)
	bmap.bank = bank
	return bmap
}

// NewROM creates a ROM bank mapping
func NewROM(address uint16, size int) *BankMap {
	return NewBankMap(address, size, true, true)
}

// NewRAM creates a RAM bank mapping
func NewRAM(address uint16, size int) *BankMap {
	return NewBankMap(address, size, false, true)
}

// Init inits the bankmap
func (bmap *BankMap) Init() {
	bmap.Bus().Init()
	bmap.active = bmap.init
}

// Reset resets the bankmap
func (bmap *BankMap) Reset() {
	bmap.Bus().Reset()
	bmap.active = bmap.init
}

// Bank returns the memory bank device
func (bmap *BankMap) Bank() *Bank { return bmap.bank }

// Bus returns the bus device
func (bmap *BankMap) Bus() device.BusDevice { return bmap.bus }

// Active returns if bank is active
func (bmap *BankMap) Active() bool { return bmap.active }

// SetActive sets bank active/inactive
func (bmap *BankMap) SetActive(active bool) { bmap.active = active }
