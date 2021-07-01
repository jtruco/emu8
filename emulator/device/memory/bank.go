package memory

// -----------------------------------------------------------------------------
// Memory bank device
// -----------------------------------------------------------------------------

// Bank is a memory bank
type Bank struct {
	data     []byte // The bank bytes
	readOnly bool   // Is a r/w or ro bank
}

// NewBank creates a new memory bank
func NewBank(size uint16, readonly bool) *Bank {
	bank := new(Bank)
	bank.data = make([]byte, size)
	bank.readOnly = readonly
	return bank
}

// Data gets bank data
func (bank *Bank) Data() []byte { return bank.data }

// IsReadOnly returns if is a read only bank
func (bank *Bank) IsReadOnly() bool { return bank.readOnly }

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
	return bank.data[address]
}

// Write writes a byte to the bank address
func (bank *Bank) Write(address uint16, data byte) {
	if !bank.readOnly {
		bank.data[address] = data
	}
}
