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

// Joystick emulation
type Joystick struct {
	model    byte // Joystick model
	kempston byte // Kempston state
}

// NewJoystick creates the spectrum joystick
func NewJoystick() *Joystick {
	joy := new(Joystick)
	joy.model = JoystickKempston
	return joy
}

// GetKempston gets kempston status
func (joy *Joystick) GetKempston() byte { return joy.kempston }

// Init initializes the device
func (joy *Joystick) Init() { joy.Reset() }

// Reset resets the device
func (joy *Joystick) Reset() { joy.kempston = 0x00 }

// SetAxis sets axis value
func (joy *Joystick) SetAxis(axis byte, value byte) {
	if axis == 0 { // right / left
		if value == 0 {
			joy.kempston &= ^_KempstonRight
			joy.kempston &= ^_KempstonLeft
		} else if value < 128 {
			joy.kempston |= _KempstonRight
		} else {
			joy.kempston |= _KempstonLeft
		}
	} else if axis == 1 { // down / up
		if value == 0 {
			joy.kempston &= ^_KempstonDown
			joy.kempston &= ^_KempstonUp
		} else if value < 128 {
			joy.kempston |= _KempstonDown
		} else {
			joy.kempston |= _KempstonUp
		}
	}
}

// SetButton sets button state
func (joy *Joystick) SetButton(button byte, state byte) {
	switch button {
	case 0:
		if state > 0 {
			joy.kempston |= _KempstonButton1
		} else {
			joy.kempston &= ^_KempstonButton1
		}
	case 1:
		if state > 0 {
			joy.kempston |= _KempstonButton2
		} else {
			joy.kempston &= ^_KempstonButton2
		}
	case 2:
		if state > 0 {
			joy.kempston |= _KempstonButton3
		} else {
			joy.kempston &= ^_KempstonButton3
		}
	}
}
