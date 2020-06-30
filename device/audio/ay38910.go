package audio

import (
	"math"

	"github.com/jtruco/emu8/device"
)

// -----------------------------------------------------------------------------
// General Instruments AY-3-891x - Programmable Sound Generator
// -----------------------------------------------------------------------------
// Only AY-3-8912 emulation supported

// AY38910 constants
const (
	AY38910Nreg         = 0X10   // 16 registers
	AY38910VolumeLevels = 0X10   // 16 volume levels
	AY38910VolumeRange  = 0x7fff // 32767
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
var (
	ay38910Masks = [AY38910Nreg]byte{
		0xFF, 0x0F, 0xFF, 0x0F, 0xFF, 0x0F, 0x1F, 0xFF,
		0x1F, 0x1F, 0x1F, 0xFF, 0xFF, 0x0F, 0xFF, 0xFF}
)

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
	counter   int
	nsample   int
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
	ay.channelA.initLevels(float64(config.Rate) * 2 / 3)
	ay.channelB.envelope = &ay.envelope
	ay.channelB.noise = &ay.noise
	ay.channelB.initLevels(float64(config.Rate) * 1 / 3)
	ay.channelC.envelope = &ay.envelope
	ay.channelC.noise = &ay.noise
	ay.channelC.initLevels(float64(config.Rate) * 2 / 3)
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
			data &= ay.readSelected()
		}
	case AY38910DataPortB: // port B
		if ay.OnReadPortB != nil {
			data &= ay.OnReadPortB()
		}
		if !ay.inPortB {
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
	// update envelope every 16 clocks
	if ay.counter&0xff == 0 {
		ay.envelope.onClock()
	}
	// update tone every 8 clocks (125 Khz)
	ay.noise.onClock()
	ay.channelA.onClock()
	ay.channelB.onClock()
	ay.channelC.onClock()
	// create audio sample
	left := ay.channelC.level + ay.channelB.level
	right := ay.channelA.level + ay.channelB.level
	sample := left + right // mono output
	index := int(float32(ay.nsample) * ay.config.Rate)
	ay.buffer.AddSample(index, sample)
	ay.nsample++
}

// -----------------------------------------------------------------------------
// AY38910 - Channel
// -----------------------------------------------------------------------------

// AY38910Channel audio channel
type AY38910Channel struct {
	volume       uint8
	period       uint16
	output       uint8
	counter      uint16
	level        uint16
	toneEnabled  bool
	noiseEnabled bool
	useEnvelope  bool
	levels       [AY38910VolumeLevels]uint16
	envelope     *AY38910Envelope
	noise        *AY38910Noise
}

func (c *AY38910Channel) initLevels(factor float64) {
	// level = max / sqrt(2)^(15-nn)
	factor /= 2 // mono output
	for l := 0; l < AY38910VolumeLevels; l++ {
		val := float64(AY38910VolumeRange) / math.Pow(math.Sqrt(2), float64(AY38910VolumeLevels-1-l))
		c.levels[l] = uint16(val * factor)
	}
}

func (c *AY38910Channel) reset() {
	c.volume = 0
	c.period = 1
	c.output = 0
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
		c.output ^= 0xff
		c.counter = 0
	}
	var volume uint8
	if c.useEnvelope {
		volume = c.envelope.volume
	} else {
		volume = c.volume
	}
	output := (c.noiseEnabled && c.noise.output != 0) || (c.toneEnabled && c.output != 0)
	if output {
		c.level = c.levels[volume]
	} else {
		c.level = -c.levels[volume]
	}
}

// -----------------------------------------------------------------------------
// AY38910 - Envelope
// -----------------------------------------------------------------------------

var (
	ay38910Shapes = [][]uint8{
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, // 0x00
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, // 0x01
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, // 0x02
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, // 0x03
		//
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}, // 0x04
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}, // 0x05
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}, // 0x06
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}, // 0x07
		//
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, // 0x08
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, // 0x09
		//
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, // 0x0a
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 15},                                                   // 0x0b
		//
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 15}, // 0x0c
		{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 15}, // 0x0d
		//
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, // 0x0e
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0},                                                    // 0x0f
	}
)

// AY38910Envelope audio envelope
type AY38910Envelope struct {
	volume     uint8
	period     uint16
	counter    uint16
	shouldHold bool
	shape      []uint8
	shapePos   int
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
	e.shapePos = 0
	e.shape = ay38910Shapes[shape]
	e.volume = e.shape[0]
	if shape < 0x08 {
		e.shouldHold = true
	} else {
		e.shouldHold = (shape & 0x01) != 0
	}
}

func (e *AY38910Envelope) onClock() {
	e.counter++
	if e.counter == e.period {
		e.counter = 0
		e.volume = e.shape[e.shapePos]
		e.shapePos++
		if e.shapePos == len(e.shape) {
			if e.shouldHold {
				e.shapePos--
			} else {
				e.shapePos = 0
			}
		}
	}
}

// -----------------------------------------------------------------------------
// AY38910 - Noise
// -----------------------------------------------------------------------------

// AY38910Noise audio noise
type AY38910Noise struct {
	output  uint8
	period  uint8
	counter uint8
	rng     int
}

func (n *AY38910Noise) reset() {
	n.output = 0xff
	n.counter = 0
	n.period = 0
	n.rng = 1
}

func (n *AY38910Noise) onClock() {
	n.counter++
	if n.counter == n.period {
		n.counter = 0
		if ((n.rng + 1) & 0x02) != 0 {
			n.output ^= 0xff
		}
		if (n.rng & 0x01) != 0 {
			n.rng ^= 0x24000
		}
		n.rng >>= 1
	}
}
