package cpu

// DataBus is a 8 bit data bus of 16 bit address
type DataBus interface {
	// Access bus access without RW/RD request
	Access(address uint16)
	// Read reads one byte from address
	Read(address uint16) byte
	// Write writes a byte at address
	Write(address uint16, data byte)
}
