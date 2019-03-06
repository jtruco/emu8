package memory

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Memory bank device
// -----------------------------------------------------------------------------

// Bank is a memory bank of bytes
type Bank struct {
	size     int
	readonly bool
	data     []byte
}

// NewBank creates a new memory bank
func NewBank(size int, readonly bool) *Bank {
	bank := &Bank{}
	bank.size = size
	bank.readonly = readonly
	bank.data = make([]byte, size)
	return bank
}

// Data initialices bank data
func (bank *Bank) Data() []byte {
	return bank.data
}

// Load loads data at address
func (bank *Bank) Load(address uint16, data []byte) {
	copy(bank.data[address:], data[:])
}

// ReadOnly is a read only bank
func (bank *Bank) ReadOnly() bool {
	return bank.readonly
}

// Device interface

// Init initialices bank data
func (bank *Bank) Init() {
	bank.Reset()
}

// Reset resets bank data
func (bank *Bank) Reset() {
	for i := 0; i < bank.size; i++ {
		bank.data[i] = 0
	}
}

// Bus interface

// Access access to data address
func (bank *Bank) Access(address uint16) {}

// Read reads a byte from the bank address
func (bank *Bank) Read(address uint16) byte {
	return bank.data[address]
}

// Write writes a byte to the bank address
func (bank *Bank) Write(address uint16, data byte) {
	if !bank.readonly {
		bank.data[address] = data
	}
}

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
