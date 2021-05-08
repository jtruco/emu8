package audio

import (
	"github.com/jtruco/emu8/emulator/device"
)

// -----------------------------------------------------------------------------
// General Instruments AY-3-891x - Programmable Sound Generator
// -----------------------------------------------------------------------------
// Only AY-3-8912 emulation supported

// AY38910 constants
const (
	AY38910Nreg        = 0x10                      // 16 registers
	AY38910Nlevels     = 0x10                      // 16 volume levels
	AY38910Nchannels   = 0x03                      // 3 channels
	AY38910VolumeRange = 0x7fff / AY38910Nchannels // 32767 / N channels
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
	AY38910ChannelAFrequencyLow = iota
	AY38910ChannelAFrequencyHigh
	AY38910ChannelBFrequencyLow
	AY38910ChannelBFrequencyHigh
	AY38910ChannelCFrequencyLow
	AY38910ChannelCFrequencyHigh
	AY38910NoiseFrequency
	AY38910MixerControl
	AY38910ChannelAVolume
	AY38910ChannelBVolume
	AY38910ChannelCVolume
	AY38910EnvelopeFrequencyLow
	AY38910EnvelopeFrequencyHigh
	AY38910EnvelopeShape
	AY38910DataPortA
	AY38910DataPortB
)

// AY38910 register data
var ay38910Masks = [AY38910Nreg]byte{
	0xFF, 0x0F, 0xFF, 0x0F, 0xFF, 0x0F, 0x1F, 0xFF,
	0x1F, 0x1F, 0x1F, 0xFF, 0xFF, 0x0F, 0xFF, 0xFF}

var ay38910VolumeLevels = [AY38910Nlevels]float32{
	0.0, 0.00999465934234, 0.0144502937362, 0.0210574502174,
	0.0307011520562, 0.0455481803616, 0.0644998855573, 0.107362478065,
	0.126588845655, 0.20498970016, 0.292210269322, 0.372838941024,
	0.492530708782, 0.635324635691, 0.805584802014, 1.0}

// -----------------------------------------------------------------------------
// AY38910
// -----------------------------------------------------------------------------

// AY38910 Crtc Device
type AY38910 struct {
	config    *Config
	buffer    *Buffer
	registers [AY38910Nreg]*byte
	selected  byte
	control   byte
	inPortA   bool
	inPortB   bool
	counter   byte
	nsample   float32
	// registers
	ChannelAFrequencyLow  byte
	ChannelAFrequencyHigh byte
	ChannelBFrequencyLow  byte
	ChannelBFrequencyHigh byte
	ChannelCFrequencyLow  byte
	ChannelCFrequencyHigh byte
	NoiseFrequency        byte
	MixerControl          byte
	ChannelAVolume        byte
	ChannelBVolume        byte
	ChannelCVolume        byte
	EnvelopeFrequencyLow  byte
	EnvelopeFrequencyHigh byte
	EnvelopeShape         byte
	DataPortA             byte
	DataPortB             byte
	// audio
	channelA AY38910Channel
	channelB AY38910Channel
	channelC AY38910Channel
	envelope AY38910Envelope
	noise    AY38910Noise
	// callbacks
	OnReadPortA  device.ReadCallback
	OnWritePortA device.WriteCallback
	OnReadPortB  device.ReadCallback
	OnWritePortB device.WriteCallback
}

// NewAY38910 creates new PSG
func NewAY38910(config *Config) *AY38910 {
	ay := new(AY38910)
	ay.config = config
	ay.buffer = NewBuffer(config.Samples)
	ay.buffer.SetFilter(NewSmaFilter(3)) // window = 8
	ay.registers = [AY38910Nreg]*byte{
		&ay.ChannelAFrequencyLow,
		&ay.ChannelAFrequencyHigh,
		&ay.ChannelBFrequencyLow,
		&ay.ChannelBFrequencyHigh,
		&ay.ChannelCFrequencyLow,
		&ay.ChannelCFrequencyHigh,
		&ay.NoiseFrequency,
		&ay.MixerControl,
		&ay.ChannelAVolume,
		&ay.ChannelBVolume,
		&ay.ChannelCVolume,
		&ay.EnvelopeFrequencyLow,
		&ay.EnvelopeFrequencyHigh,
		&ay.EnvelopeShape,
		&ay.DataPortA,
		&ay.DataPortB}
	ay.channelA.envelope = &ay.envelope
	ay.channelA.noise = &ay.noise
	ay.channelA.initLevels(config.Rate)
	ay.channelB.envelope = &ay.envelope
	ay.channelB.noise = &ay.noise
	ay.channelB.initLevels(config.Rate)
	ay.channelC.envelope = &ay.envelope
	ay.channelC.noise = &ay.noise
	ay.channelC.initLevels(config.Rate)
	return ay
}

// Config returns the audio configuration
func (ay *AY38910) Config() *Config { return ay.config }

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
	ay.counter = 0
	ay.nsample = 0
	ay.channelA.reset()
	ay.channelB.reset()
	ay.channelC.reset()
	ay.noise.reset()
	ay.envelope.reset()
}

// Audio interface

// Buffer gets audio buffer
func (ay *AY38910) Buffer() *Buffer {
	return ay.buffer
}

// EndFrame ends audio frame
func (ay *AY38910) EndFrame() {
	ay.nsample = 0
}

// io operations

// Read reads data
func (ay *AY38910) Read() byte {
	var data byte = 0xff
	switch ay.selected {
	case AY38910DataPortA: // port A
		if ay.OnReadPortA != nil {
			data &= ay.OnReadPortA()
		}
		if !ay.inPortA {
			data &= ay.ReadRegister(ay.selected)
		}
	case AY38910DataPortB: // port B
		if ay.OnReadPortB != nil {
			data &= ay.OnReadPortB()
		}
		if !ay.inPortB {
			data &= ay.ReadRegister(ay.selected)
		}
	default:
		data &= ay.ReadRegister(ay.selected)
	}
	return data
}

// Write writes data
func (ay *AY38910) Write(data byte) {
	switch ay.control {
	case AY38910SelectRegister:
		ay.SelectRegister(data & 0x0f)
	case AY38910WriteRegister:
		ay.WriteRegister(ay.selected, data)
	}
}

// register operations

// Selected selected register
func (ay *AY38910) Selected() byte { return ay.selected }

// Register gets register value at index
func (ay *AY38910) Register(index byte) byte { return *ay.registers[index] }

// SelectRegister selects current register
func (ay *AY38910) SelectRegister(selected byte) {
	ay.selected = selected
}

// ReadRegister returns register value
func (ay *AY38910) ReadRegister(register byte) byte {
	if register < AY38910Nreg {
		return *ay.registers[register]
	}
	return 0 // write only
}

// WriteRegister writes value to register
func (ay *AY38910) WriteRegister(register, data byte) {
	*ay.registers[register] = data & ay38910Masks[register]

	switch register {
	case AY38910ChannelAFrequencyLow, AY38910ChannelAFrequencyHigh:
		ay.channelA.setPeriod(ay.ChannelAFrequencyHigh, ay.ChannelAFrequencyLow)
	case AY38910ChannelBFrequencyLow, AY38910ChannelBFrequencyHigh:
		ay.channelB.setPeriod(ay.ChannelBFrequencyHigh, ay.ChannelBFrequencyLow)
	case AY38910ChannelCFrequencyLow, AY38910ChannelCFrequencyHigh:
		ay.channelC.setPeriod(ay.ChannelCFrequencyHigh, ay.ChannelCFrequencyLow)
	case AY38910NoiseFrequency:
		ay.noise.period = data
	case AY38910MixerControl:
		ay.channelA.toneEnabled = ((data & 0x01) == 0)
		ay.channelB.toneEnabled = ((data & 0x02) == 0)
		ay.channelC.toneEnabled = ((data & 0x04) == 0)
		ay.channelA.noiseEnabled = ((data & 0x08) == 0)
		ay.channelB.noiseEnabled = ((data & 0x10) == 0)
		ay.channelC.noiseEnabled = ((data & 0x20) == 0)
		ay.inPortA = ((data & 0x40) == 0)
		ay.inPortB = ((data & 0x80) == 0)
	case AY38910ChannelAVolume:
		ay.channelA.volume = (data & 0x0f)
		ay.channelA.useEnvelope = ((data & 0x10) != 0)
	case AY38910ChannelBVolume:
		ay.channelB.volume = (data & 0x0f)
		ay.channelB.useEnvelope = ((data & 0x10) != 0)
	case AY38910ChannelCVolume:
		ay.channelC.volume = (data & 0x0f)
		ay.channelC.useEnvelope = ((data & 0x10) != 0)
	case AY38910EnvelopeFrequencyLow, AY38910EnvelopeFrequencyHigh:
		ay.envelope.setPeriod(ay.EnvelopeFrequencyHigh, ay.EnvelopeFrequencyLow)
	case AY38910EnvelopeShape:
		ay.envelope.setShape(data)
	case AY38910DataPortA:
		if !ay.inPortA && ay.OnWritePortA != nil {
			ay.OnWritePortA(ay.DataPortA)
		}
	case AY38910DataPortB:
		if !ay.inPortB && ay.OnWritePortB != nil {
			ay.OnWritePortB(ay.DataPortB)
		}
	default:
	}
}

// emulation

// Emulate emulates Tstates
func (ay *AY38910) Emulate(tstates int) {
	for i := 0; i < tstates; i++ {
		ay.OnClock()
	}
}

// OnClock emulates one clock cycle (1MHz)
func (ay *AY38910) OnClock() {
	ay.counter++
	if ay.counter&0x07 != 0 {
		return
	}
	// update noise every 8 clocks (125 Khz)
	ay.noise.onClock()
	// update tone every 8 clocks (125 Khz)
	ay.channelA.onClock()
	ay.channelB.onClock()
	ay.channelC.onClock()
	// update envelope every 16 clocks
	if ay.counter&0xff == 0 {
		ay.envelope.onClock()
	}
	ay.counter = 0
	// create audio sample
	mix := ay.channelA.level + ay.channelB.level + ay.channelC.level
	index := int(ay.nsample)
	ay.buffer.AddSample(index, mix)
	ay.nsample += ay.config.Rate
}

// -----------------------------------------------------------------------------
// AY38910 - Channel
// -----------------------------------------------------------------------------

// AY38910Channel audio channel
type AY38910Channel struct {
	volume       uint8
	period       uint16
	output       bool
	counter      uint16
	level        uint16
	toneEnabled  bool
	noiseEnabled bool
	useEnvelope  bool
	levels       [AY38910Nlevels]uint16
	envelope     *AY38910Envelope
	noise        *AY38910Noise
}

func (c *AY38910Channel) initLevels(factor float32) {
	for l := 0; l < AY38910Nlevels; l++ {
		c.levels[l] = uint16(ay38910VolumeLevels[l] * AY38910VolumeRange * factor)
	}
}

func (c *AY38910Channel) reset() {
	c.volume = 0
	c.period = 1
	c.output = false
	c.counter = 0
	c.toneEnabled = false
	c.noiseEnabled = false
	c.useEnvelope = false
	c.level = 0
}

func (c *AY38910Channel) setPeriod(high, low uint8) {
	c.period = uint16(high)<<8 | uint16(low)
	if c.period == 0 {
		c.period = 1
	}
	if c.counter >= (c.period << 1) {
		c.counter %= (c.period << 1)
	}
}

func (c *AY38910Channel) onClock() {
	c.counter++
	if c.counter == c.period {
		c.output = !c.output
		c.counter = 0
	}
	enable := (c.toneEnabled && c.output) || (c.noiseEnabled && c.noise.output)
	if enable {
		if c.useEnvelope {
			c.level = c.levels[c.envelope.volume]
		} else {
			c.level = c.levels[c.volume]
		}
	} else {
		c.level = 0
	}
}

// -----------------------------------------------------------------------------
// AY38910 - Envelope
// -----------------------------------------------------------------------------

var ay38910Shapes = [][]uint8{
	/* 0 0 X X */
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	/* 0 1 X X */
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	/* 1 0 0 0 */
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	/* 1 0 0 1 */
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	/* 1 0 1 0 */
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	/* 1 0 1 1 */
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15},
	/* 1 1 0 0 */
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	/* 1 1 0 1 */
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15},
	/* 1 1 1 0 */
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	/* 1 1 1 1 */
	{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

// AY38910Envelope audio envelope
type AY38910Envelope struct {
	volume  uint8
	period  uint16
	counter uint16
	shape   []uint8
	hold    bool
	pos     int
}

func (e *AY38910Envelope) reset() {
	e.volume = 0
	e.period = 1
	e.counter = 0
	e.setShape(0)
}

func (e *AY38910Envelope) setPeriod(high, low uint8) {
	e.period = uint16(high)<<8 | uint16(low)
}

func (e *AY38910Envelope) setShape(shape uint8) {
	e.shape = ay38910Shapes[shape]
	e.pos = 0
	e.hold = (shape&0x08 == 0) || (shape&0x01 != 0)
}

func (e *AY38910Envelope) onClock() {
	e.counter++
	if e.counter == e.period {
		e.counter = 0
		e.volume = e.shape[e.pos]
		if e.pos == 0x1f {
			if !e.hold {
				e.pos = 0
			}
		} else {
			e.pos++
		}
	}
}

// -----------------------------------------------------------------------------
// AY38910 - Noise
// -----------------------------------------------------------------------------

// AY38910Noise audio noise
type AY38910Noise struct {
	output   bool
	period   uint8
	counter  uint8
	prescale bool
	rng      uint32
}

func (n *AY38910Noise) reset() {
	n.output = false
	n.period = 0
	n.counter = 0
	n.prescale = false
	n.rng = 0x01
}

func (n *AY38910Noise) onClock() {
	n.counter++
	if n.counter == n.period {
		n.counter = 0
		n.prescale = !n.prescale
		if !n.prescale {
			// https://github.com/mamedev/mame/blob/master/src/devices/sound/ay8910.cpp
			// The Random Number Generator of the 8910 is a 17-bit shift
			// register. The input to the shift register is bit0 XOR bit3
			// (bit0 is the output). This was verified on AY-3-8910 and YM2149 chips.
			n.rng ^= (((n.rng & 1) ^ ((n.rng >> 3) & 1)) << 17)
			n.rng >>= 1
		}
		n.output = (n.rng & 0x01) != 0
	}
}
