package cpc

// -----------------------------------------------------------------------------
// Amstrad CPC - Joystick emulation
// -----------------------------------------------------------------------------

// Joystick emulation
type Joystick struct {
	id       byte
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

// ID returns the joystick ID
func (joy *Joystick) ID() byte { return joy.id }

// SetAxis sets axis value
func (joy *Joystick) SetAxis(axis byte, value byte) {
	if axis == 0 { // right / left
		if value == 0 {
			joy.keyboard.ProcessKey(CpcKeyJoyRight, false)
			joy.keyboard.ProcessKey(CpcKeyJoyLeft, false)
		} else if value < 128 {
			joy.keyboard.ProcessKey(CpcKeyJoyRight, true)
		} else {
			joy.keyboard.ProcessKey(CpcKeyJoyLeft, true)
		}
	} else if axis == 1 { // down / up
		if value == 0 {
			joy.keyboard.ProcessKey(CpcKeyJoyDown, false)
			joy.keyboard.ProcessKey(CpcKeyJoyUp, false)
		} else if value < 128 {
			joy.keyboard.ProcessKey(CpcKeyJoyDown, true)
		} else {
			joy.keyboard.ProcessKey(CpcKeyJoyUp, true)
		}
	}
}

// SetButton sets button state
func (joy *Joystick) SetButton(button byte, state byte) {
	switch button {
	case 0:
		joy.keyboard.ProcessKey(CpcKeyJoyFire1, state > 0)
	case 1:
		joy.keyboard.ProcessKey(CpcKeyJoyFire2, state > 0)
	}
}
