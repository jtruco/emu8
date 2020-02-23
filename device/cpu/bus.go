package cpu

// -----------------------------------------------------------------------------
// DataBus
// -----------------------------------------------------------------------------

// DataBus is a 8 bit data bus of 16 bit address
type DataBus interface {
	// Read reads one byte from address
	Read(address uint16) byte
	// Write writes a byte at address
	Write(address uint16, data byte)
}
