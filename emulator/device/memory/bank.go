package memory

import "github.com/jtruco/emu8/emulator/device/bus"

// -----------------------------------------------------------------------------
// Memory bank device
// -----------------------------------------------------------------------------

// Bank is a memory bank
type Bank struct {
	data         []byte       // The bank bytes
	readonly     bool         // Is a r/w or ro bank
	OnAccess     bus.Callback // On bus access callback
	OnPostAccess bus.Callback // On bus post access callback
}

// NewBank creates a new memory bank
func NewBank(size uint16, readonly bool) *Bank {
	bank := new(Bank)
	bank.data = make([]byte, size)
	bank.readonly = readonly
	return bank
}

// Data gets bank data
func (bank *Bank) Data() []byte { return bank.data }

// ReadOnly returns if is a read only bank
func (bank *Bank) ReadOnly() bool { return bank.readonly }

// Size return bank size
func (bank *Bank) Size() uint16 { return uint16(len(bank.data)) }

// Load loads data at address
func (bank *Bank) Load(address uint16, data []byte) {
	copy(bank.data[address:], data[:])
}

// Save saves bank data to slice at address
func (bank *Bank) Save(data []byte) {
	copy(data[:], bank.data[:])
}

// Device interface

// Init initializes bank data
func (bank *Bank) Init() {
	bank.Reset()
}

// Reset resets bank data
func (bank *Bank) Reset() {
	for i := 0; i < len(bank.data); i++ {
		bank.data[i] = 0
	}
}

// Bus interface

// Read reads a byte from the bank address
func (bank *Bank) Read(address uint16) byte {
	// on access
	if bank.OnAccess != nil {
		bank.OnAccess(bus.EventRead, address)
	}

	// memory read
	data := bank.data[address]

	// on post access
	if bank.OnPostAccess != nil {
		bank.OnPostAccess(bus.EventAfterRead, address)
	}

	return data
}

// Write writes a byte to the bank address
func (bank *Bank) Write(address uint16, data byte) {
	// on access
	if bank.OnAccess != nil {
		bank.OnAccess(bus.EventWrite, address)
	}

	// memory read
	if !bank.readonly {
		bank.data[address] = data
	}

	// on post access
	if bank.OnPostAccess != nil {
		bank.OnPostAccess(bus.EventAfterWrite, address)
	}
}
