package memory

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Memory bank device
// -----------------------------------------------------------------------------

// Bank is a memory bank of bytes
type Bank struct {
	size         int
	readonly     bool
	data         []byte
	OnAccess     device.EventBus
	OnPostAccess device.EventBus
}

// NewBank creates a new memory bank
func NewBank(size int, readonly bool) *Bank {
	bank := new(Bank)
	bank.size = size
	bank.readonly = readonly
	bank.data = make([]byte, size)
	return bank
}

// Data gets bank data
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

// Init initializes bank data
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

// Read reads a byte from the bank address
func (bank *Bank) Read(address uint16) byte {
	bank.OnAccess.Emit(device.NewBusEvent(device.EventBusRead, address))
	data := bank.data[address]
	bank.OnPostAccess.Emit(device.NewBusEvent(device.EventBusAfterRead, address))
	return data
}

// Write writes a byte to the bank address
func (bank *Bank) Write(address uint16, data byte) {
	bank.OnAccess.Emit(device.NewBusEvent(device.EventBusWrite, address))
	if !bank.readonly {
		bank.data[address] = data
	}
	bank.OnPostAccess.Emit(device.NewBusEvent(device.EventBusAfterWrite, address))
}
