package audio

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// General Instruments AY-3-891x - Programmable Sound Generator
// -----------------------------------------------------------------------------
// Only AY-3-8912 emulation supported

// http://www.cpcwiki.eu/index.php/PSG
// http://www.cpctech.org.uk/docs/ay.html
// http://www.cpctech.org.uk/docs/psgnotes.htm
// http://www.cpctech.org.uk/docs/ay38912/psgspec.htm
// https://github.com/mamedev/mame/blob/master/src/devices/sound/ay8910.cpp

// AY38910 constants
const (
	AY38910Nreg = 0X10 // 16 registers
)

// AY38910 models
const (
	AY38910Model = iota
	AY38912Model
	AY38913Model
)

// AY38910 control functions
const (
	AY38910Inactive = iota
	AY38910ReadRegister
	AY38910WriteRegister
	AY38910SelectRegister
)

// AY38910 register constants
const (
	AY38910ChannelAToneFreqLow8Bit = iota
	AY38910ChannelAToneFreqHigh4Bit
	AY38910ChannelBToneFreqLow8Bit
	AY38910ChannelBToneFreqHigh4Bit
	AY38910ChannelCToneFreqLow8Bit
	AY38910ChannelCToneFreqHigh4Bit
	AY38910NoiseFrequency
	AY38910MixerControl
	AY38910ChannelAVolume
	AY38910ChannelBVolume
	AY38910ChannelCVolume
	AY38910VolumeEnvFreqLow
	AY38910VolumeEnvFreqHigh
	AY38910VolumeEnvShape
	AY38910ExternalDataPortA
	AY38910ExternalDataPortB
)

// AY38910 register data
var (
	ay38910Masks = [AY38910Nreg]byte{
		0xFF, 0x0F, 0xFF, 0x0F, 0xFF, 0x0F, 0x1F, 0xFF,
		0x1F, 0x1F, 0x1F, 0xFF, 0xFF, 0x0F, 0xFF, 0xFF}
)

// AY38910 Crtc Device
type AY38910 struct {
	registers [AY38910Nreg]*byte
	selected  byte
	control   byte
	// registers
	ChannelAToneFreqLow8Bit  byte
	ChannelAToneFreqHigh4Bit byte
	ChannelBToneFreqLow8Bit  byte
	ChannelBToneFreqHigh4Bit byte
	ChannelCToneFreqLow8Bit  byte
	ChannelCToneFreqHigh4Bit byte
	NoiseFrequency           byte
	MixerControl             byte
	ChannelAVolume           byte
	ChannelBVolume           byte
	ChannelCVolume           byte
	VolumeEnvFreqLow         byte
	VolumeEnvFreqHigh        byte
	VolumeEnvShape           byte
	ExternalDataPortA        byte
	ExternalDataPortB        byte
	// callbacks
	OnReadPortA  device.ReadCallback
	OnWritePortA device.WriteCallback
	OnReadPortB  device.ReadCallback
	OnWritePortB device.WriteCallback
}

// NewAY38910 creates new PSG
func NewAY38910() *AY38910 {
	ay := new(AY38910)
	ay.registers = [AY38910Nreg]*byte{
		&ay.ChannelAToneFreqLow8Bit,
		&ay.ChannelAToneFreqHigh4Bit,
		&ay.ChannelBToneFreqLow8Bit,
		&ay.ChannelBToneFreqHigh4Bit,
		&ay.ChannelCToneFreqLow8Bit,
		&ay.ChannelCToneFreqHigh4Bit,
		&ay.NoiseFrequency,
		&ay.MixerControl,
		&ay.ChannelAVolume,
		&ay.ChannelBVolume,
		&ay.ChannelCVolume,
		&ay.VolumeEnvFreqLow,
		&ay.VolumeEnvFreqHigh,
		&ay.VolumeEnvShape,
		&ay.ExternalDataPortA,
		&ay.ExternalDataPortB}
	return ay
}

// properties

// Control returns the control register
func (ay *AY38910) Control() byte { return ay.control }

// SetControl sets de control register
func (ay *AY38910) SetControl(value byte) { ay.control = (value & 0xc0) >> 6 }

// device inerface

// Init the PSG
func (ay *AY38910) Init() { ay.Reset() }

// Reset the PSG
func (ay *AY38910) Reset() {
	ay.selected = 0
	for i := byte(0); i < AY38910Nreg; i++ {
		ay.WriteRegister(i, 0)
	}
	ay.selected = 0
	ay.control = 0
}

// io operations

// Read reads data
func (ay *AY38910) Read() byte {
	var data byte = 0xff
	switch ay.selected {
	case AY38910ExternalDataPortA: // port A
		if ay.OnReadPortA != nil {
			data = ay.OnReadPortA()
		}
		if (ay.MixerControl & 0x40) != 0 {
			data &= ay.readSelected()
		}
	case AY38910ExternalDataPortB: // port B
		if ay.OnReadPortB != nil {
			data = ay.OnReadPortB()
		}
		if (ay.MixerControl & 0x80) != 0 {
			data &= ay.readSelected()
		}
	default:
		data &= ay.readSelected()
	}
	return data
}

// Write writes data
func (ay *AY38910) Write(data byte) {
	switch ay.control {
	case AY38910SelectRegister:
		ay.SelectRegister(data & 0x0f)
	case AY38910WriteRegister:
		ay.writeSelected(data)
	}
}

// register operations

// SelectRegister selects current register
func (ay *AY38910) SelectRegister(selected byte) {
	ay.selected = selected
}

// readSelected returns current register value
func (ay *AY38910) readSelected() byte {
	return ay.ReadRegister(ay.selected)
}

// ReadRegister returns register value
func (ay *AY38910) ReadRegister(register byte) byte {
	if register < AY38910Nreg {
		return *ay.registers[register]
	}
	return 0 // write only
}

// writeSelected writes value to selected register
func (ay *AY38910) writeSelected(data byte) {
	ay.WriteRegister(ay.selected, data)
}

// WriteRegister writes value to register
func (ay *AY38910) WriteRegister(register, data byte) {
	*ay.registers[register] = data & ay38910Masks[register]

	switch register {
	// TODO : psg registers
	case AY38910ExternalDataPortA:
		if ay.MixerControl&0x40 != 0 && ay.OnWritePortA != nil {
			ay.OnWritePortA(ay.ExternalDataPortA)
		}
	case AY38910ExternalDataPortB:
		if ay.MixerControl&0x80 != 0 && ay.OnWritePortB != nil {
			ay.OnWritePortB(ay.ExternalDataPortB)
		}
	default:
	}
}

// emulation

// sound generation
