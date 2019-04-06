package keyboard

import "github.com/jtruco/emu8/device"

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
