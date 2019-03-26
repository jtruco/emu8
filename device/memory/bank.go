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
}

// NewBank creates a new memory bank
func NewBank(size int, readonly bool) *Bank {
	bank := &Bank{}
	bank.size = size
	bank.readonly = readonly
	bank.data = make([]byte, size)
	bank.listeners = make([]device.BusListener, 0)
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
func (bank *Bank) Access(address uint16) {
	if len(bank.listeners) > 0 {
		bank.notifyListeners(address, device.EventBusAccess, device.OrderBefore)
		bank.notifyListeners(address, device.EventBusAccess, device.OrderAfter)
	}
}

// Read reads a byte from the bank address
func (bank *Bank) Read(address uint16) byte {
	if len(bank.listeners) > 0 {
		bank.notifyListeners(address, device.EventBusRead, device.OrderBefore)
	}
	data := bank.data[address]
	if len(bank.listeners) > 0 {
		bank.notifyListeners(address, device.EventBusRead, device.OrderAfter)
	}
	return data
}

// Write writes a byte to the bank address
func (bank *Bank) Write(address uint16, data byte) {
	if len(bank.listeners) > 0 {
		bank.notifyListeners(address, device.EventBusWrite, device.OrderBefore)
	}
	if !bank.readonly {
		bank.data[address] = data
	}
	if len(bank.listeners) > 0 {
		bank.notifyListeners(address, device.EventBusWrite, device.OrderAfter)
	}
}

// Events

// AddBusListener adds a listener tu memory bank events
func (bank *Bank) AddBusListener(listener device.BusListener) {
	bank.listeners = append(bank.listeners, listener)
}

// notifyListeners emits event and notify listeners
func (bank *Bank) notifyListeners(address uint16, event, order int) {
	busevent := device.BusEvent{
		Event:   device.Event{Type: event, Order: order},
		Address: address}
	for _, l := range bank.listeners {
		l.ProcessBusEvent(&busevent)
	}
}
