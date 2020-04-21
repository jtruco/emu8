package cpc

// -----------------------------------------------------------------------------
// Amstrad CPC - 8255 Parallel peripheral interface
// -----------------------------------------------------------------------------

// Ppi 8255 parallel peripheral interface
type Ppi struct {
	cpc     *AmstradCPC
	portA   byte
	portB   byte
	portC   byte
	control byte
	jumpers byte
}

// NewPpi creates new PPI
func NewPpi(cpc *AmstradCPC) *Ppi {
	ppi := new(Ppi)
	ppi.cpc = cpc
	return ppi
}

// Init the PPI
func (ppi *Ppi) Init() { ppi.Reset() }

// Reset the PPI
func (ppi *Ppi) Reset() {
	ppi.portA = 0x00
	ppi.portB = 0x00
	ppi.portC = 0x00
	ppi.control = 0x00
	ppi.jumpers = 0x00
}

// Read read from gatearray
func (ppi *Ppi) Read(port byte) byte {
	var data byte = 0xff
	switch port {
	case 0: // port A
		if (ppi.control & 0x10) != 0 { // input
			data = ppi.cpc.psg.Read()
		} else {
			data = ppi.portA
		}
	case 1: // port B
		if (ppi.control & 0x02) != 0 { // input
			data = ppi.jumpers
			if ppi.cpc.crtc.InVSync() {
				data |= 0x01
			}
			// TODO : tape,...
		} else {
			data = ppi.portB
		}
	case 2: // port C
		data = ppi.portC
		if (ppi.control & 0x08) != 0 { // upper nibble
			data &= 0x0f
			value := ppi.portC & 0xc0
			if value == 0xc0 {
				value = 0x80
			}
			data |= value | 0x20
		}
		if (ppi.control & 0x01) != 0 { // lower nibble
			data |= 0x0f
		}
	}
	return data
}

// Write sets current pen
func (ppi *Ppi) Write(port byte, data byte) {
	switch port {
	case 0: // port A
		ppi.portA = data
		if (ppi.control & 0x10) == 0 { // output
			ppi.cpc.psg.Write(data)
		}
	case 1: // port B
		ppi.portB = data
	case 2: // port C
		ppi.portC = data
		if (ppi.control & 0x01) == 0 { // lower nibble
			ppi.cpc.keyboard.SetRow(data)
		}
		if (ppi.control & 0x08) == 0 { // upper nibble
			// todo : tape
			ppi.cpc.psg.control = data
			ppi.cpc.psg.Write(ppi.portA)
		}
	case 3: // PPI control
		if (data & 0x80) != 0 {
			ppi.control = data
			ppi.portA = 0
			ppi.portB = 0
			ppi.portC = 0
		} else {
			bit := (data >> 1) & 0x7
			mask := byte(1 << bit)
			if (data & 0x01) != 0 {
				ppi.portC |= mask
			} else {
				ppi.portC &= ^mask
			}
			if (ppi.control & 0x01) == 0 { // lower nibble
				ppi.cpc.keyboard.SetRow(ppi.portC)
			}
			if (ppi.control & 0x08) == 0 { // upper nibble
				// TODO : tape control
				ppi.cpc.psg.control = ppi.portC
				ppi.cpc.psg.Write(ppi.portA)
			}
		}
	}
}
