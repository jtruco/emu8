package cpc

// -----------------------------------------------------------------------------
// Amstrad CPC - Programmable Sound Generator (AY-3-8912)
// -----------------------------------------------------------------------------

// Psg Programmable Sound Generator
type Psg struct {
	cpc       *AmstradCPC
	control   byte
	selected  byte
	registers []byte
}

// NewPsg creates new PSG
func NewPsg(cpc *AmstradCPC) *Psg {
	psg := new(Psg)
	psg.cpc = cpc
	psg.registers = make([]byte, 16)
	return psg
}

// Init the PSG
func (psg *Psg) Init() { psg.Reset() }

// Reset the PSG
func (psg *Psg) Reset() {
	psg.control = 0
	psg.selected = 0
	for n := 0; n < 16; n++ {
		psg.registers[n] = 0
	}
}

// Read reads data
func (psg *Psg) Read() byte {
	var data byte = 0xff
	if psg.selected == 14 { // port A
		data &= psg.cpc.keyboard.State()
		if (psg.registers[7] & 0x40) != 0 {
			data &= psg.registers[14]
		}
	} else if psg.selected == 15 { // port B
		if (psg.registers[7] & 0x80) != 0 {
			data &= psg.registers[15]
		}
	} else if psg.selected < 14 {
		data &= psg.registers[psg.selected]
	}
	return data
}

// Write writes data
func (psg *Psg) Write(data byte) {
	control := psg.control & 0xc0
	if control == 0xc0 {
		psg.selected = data
	} else if control == 0x80 {
		if psg.selected < 16 {
			psg.registers[psg.selected] = data
		}
	}
}
