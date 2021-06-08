package spectrum

// -----------------------------------------------------------------------------
// ZX Spectrum - Joystick emulation
// -----------------------------------------------------------------------------
// Only Kempston emulation supported

// Joystick models
const (
	JoystickNone = iota
	JoystickKempston
)

// Kempston constants
const (
	_KempstonDefault = byte(0x00)
	_KempstonRight   = byte(0x01)
	_KempstonLeft    = byte(0x02)
	_KempstonDown    = byte(0x04)
	_KempstonUp      = byte(0x08)
	_KempstonButton1 = byte(0x10)
	_KempstonButton2 = byte(0x20)
	_KempstonButton3 = byte(0x40)
)

// Joystick Kempston emulation
type Joystick struct {
	id    byte // ID
	model byte // Joystick model
	state byte // Kempston state
}

// NewJoystick creates a new Kempston Joystick
func NewJoystick() *Joystick {
	joy := new(Joystick)
	joy.model = JoystickKempston
	return joy
}

// State gets kempston status
func (joy *Joystick) State() byte { return joy.state }

// Init initializes the device
func (joy *Joystick) Init() { joy.Reset() }

// Reset resets the device
func (joy *Joystick) Reset() { joy.state = 0x00 }

// ID returns the joystick ID
func (joy *Joystick) ID() byte { return joy.id }

// SetAxis sets axis value
func (joy *Joystick) SetAxis(axis byte, value byte) {
	if axis == 0 { // right / left
		if value == 0 {
			joy.state &^= _KempstonRight
			joy.state &^= _KempstonLeft
		} else if value < 128 {
			joy.state |= _KempstonRight
		} else {
			joy.state |= _KempstonLeft
		}
	} else if axis == 1 { // down / up
		if value == 0 {
			joy.state &^= _KempstonDown
			joy.state &^= _KempstonUp
		} else if value < 128 {
			joy.state |= _KempstonDown
		} else {
			joy.state |= _KempstonUp
		}
	}
}

// SetButton sets button state
func (joy *Joystick) SetButton(button byte, state byte) {
	switch button {
	case 0:
		if state > 0 {
			joy.state |= _KempstonButton1
		} else {
			joy.state &^= _KempstonButton1
		}
	case 1:
		if state > 0 {
			joy.state |= _KempstonButton2
		} else {
			joy.state &^= _KempstonButton2
		}
	case 2:
		if state > 0 {
			joy.state |= _KempstonButton3
		} else {
			joy.state &^= _KempstonButton3
		}
	}
}
