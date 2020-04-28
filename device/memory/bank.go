package memory

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Memory bank device
// -----------------------------------------------------------------------------

// Bank is a memory bank of bytes
type Bank struct {
	size      int
	readonly  bool
	data      []byte
	listeners []device.BusListener
	OnBus     device.EventBus // On bus access events
	OnPostBus device.EventBus // Post bus access events
}

// NewBank creates a new memory bank
func NewBank(size int, readonly bool) *Bank {
	bank := new(Bank)
	bank.size = size
	bank.readonly = readonly
	bank.data = make([]byte, size)
	bank.listeners = make([]device.BusListener, 0)
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
	if len(bank.listeners) > 0 {
		bank.notifyListeners(device.EventBusRead, address)
	}
	bank.OnBus.EmitEvent(device.NewBusEvent(device.EventBusRead, address))
	data := bank.data[address]
	if len(bank.listeners) > 0 {
		bank.notifyListeners(device.EventBusAfterRead, address)
	}
	return data
}

// Write writes a byte to the bank address
func (bank *Bank) Write(address uint16, data byte) {
	if len(bank.listeners) > 0 {
		bank.notifyListeners(device.EventBusWrite, address)
	}
	if !bank.readonly {
		bank.data[address] = data
	}
	if len(bank.listeners) > 0 {
		bank.notifyListeners(device.EventBusAfterWrite, address)
	}
}

// Events

// AddBusListener adds a listener tu memory bank events
func (bank *Bank) AddBusListener(listener device.BusListener) {
	bank.listeners = append(bank.listeners, listener)
}

// notifyListeners emits event and notify listeners
func (bank *Bank) notifyListeners(code int, address uint16) {
	busevent := device.NewBusEvent(code, address)
	for _, l := range bank.listeners {
		l.ProcessBusEvent(busevent)
	}
}
