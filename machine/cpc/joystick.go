package cpc

// -----------------------------------------------------------------------------
// Amstrad CPC - Joystick emulation
// -----------------------------------------------------------------------------

// Joystick emulation
type Joystick struct {
	keyboard *Keyboard
}

// NewJoystick creates the amstrad cpc joystick
func NewJoystick(keyboard *Keyboard) *Joystick {
	joy := new(Joystick)
	joy.keyboard = keyboard
	return joy
}

// Init initializes the device
func (joy *Joystick) Init() { joy.Reset() }

// Reset resets the device
func (joy *Joystick) Reset() {}

// SetAxis sets axis value
func (joy *Joystick) SetAxis(axis byte, value byte) {
	if axis == 0 { // right / left
		if value == 0 {
			joy.keyboard.processKey(CpcKeyJoyRight, false)
			joy.keyboard.processKey(CpcKeyJoyLeft, false)
		} else if value < 128 {
			joy.keyboard.processKey(CpcKeyJoyRight, true)
		} else {
			joy.keyboard.processKey(CpcKeyJoyLeft, true)
		}
	} else if axis == 1 { // down / up
		if value == 0 {
			joy.keyboard.processKey(CpcKeyJoyDown, false)
			joy.keyboard.processKey(CpcKeyJoyUp, false)
		} else if value < 128 {
			joy.keyboard.processKey(CpcKeyJoyDown, true)
		} else {
			joy.keyboard.processKey(CpcKeyJoyUp, true)
		}
	}
}

// SetButton sets button state
func (joy *Joystick) SetButton(button byte, state byte) {
	switch button {
	case 0:
		joy.keyboard.processKey(CpcKeyJoyFire1, state > 0)
	case 1:
		joy.keyboard.processKey(CpcKeyJoyFire2, state > 0)
	}
}
