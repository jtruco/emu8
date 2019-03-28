package keyboard

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Keyboard & Events
// -----------------------------------------------------------------------------

// Keyboard event types
const (
	EventKeyDown = 21 // KeyDown is a key down event
	EventKeyUp   = 22 // KeyDown is a key up event
)

// KeyEvent is a keyboard event
type KeyEvent struct {
	device.Event      // Is a device event
	Key          Key  // Machine key pressed Down/Up
	Pressed      bool // True if key is pressed down, False if is up
}

// Receiver is a component that process keyboard events
type Receiver interface {
	ProcessKeyEvent(event *KeyEvent)
}

// Keyboard is a keyboard (receiver) device
type Keyboard interface {
	device.Device // Is a device
	Receiver      // Is a keyboard event receiver
}

// -----------------------------------------------------------------------------
// Controller
// -----------------------------------------------------------------------------

// Controller is the emulator keyboard controller
type Controller struct {
	receivers map[Receiver]KeyMap
}

// NewController creates a new controller with an empty keymap
func NewController() *Controller {
	controller := Controller{}
	controller.receivers = make(map[Receiver]KeyMap)
	return &controller
}

// Receivers

// AddReceiver adds a keyboard events Receiver to the controller and asociated keymap
func (controller *Controller) AddReceiver(receiver Receiver, keymap KeyMap) {
	controller.receivers[receiver] = keymap
}

// RemoveReceiver removes the Receiver
func (controller *Controller) RemoveReceiver(receiver Receiver) {
	delete(controller.receivers, receiver)
}

// Key events

// KeyDown emits a keyboard keydown event
func (controller *Controller) KeyDown(keycode KeyCode) {
	controller.KeyEvent(keycode, KeyDown)
}

// KeyUp emits a keyboard keyup event
func (controller *Controller) KeyUp(keycode KeyCode) {
	controller.KeyEvent(keycode, KeyUp)
}

// KeyEvent emits a keyboard event
func (controller *Controller) KeyEvent(keycode KeyCode, eventType int) {
	// For every receiver checks if keycode is mapped
	for receiver, keymap := range controller.receivers {
		keys, ok := keymap[keycode]
		if ok {
			// For each key emit event to receiver
			for _, key := range keys {
				keyevent := KeyEvent{
					Event:   device.Event{Type: eventType, Order: device.OrderAfter},
					Key:     key,
					Pressed: eventType == KeyDown}
				receiver.ProcessKeyEvent(&keyevent)
			}
		}
	}
}
