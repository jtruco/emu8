// Package device contains common device components
package device

// -----------------------------------------------------------------------------
// Device
// -----------------------------------------------------------------------------

// Device is the base device component
type Device interface {
	Init()  // Init initializes the device
	Reset() // Reset resets the device
}

// -----------------------------------------------------------------------------
// Device callbacks
// -----------------------------------------------------------------------------

// Callback is a device callback
type Callback func()

// AckCallback device callback with ack control
type AckCallback func() bool

// ReadCallback line read byte callback
type ReadCallback func() byte

// WriteCallback line write byte callback
type WriteCallback func(byte)
