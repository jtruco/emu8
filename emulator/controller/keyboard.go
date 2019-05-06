package controller

import (
	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/io/keyboard"
)

// -----------------------------------------------------------------------------
// Keyboard Controller
// -----------------------------------------------------------------------------

// KeyboardController is the emulator keyboard controller
type KeyboardController struct {
	receivers map[keyboard.Receiver]keyboard.KeyMap
}

// NewKeyboardController creates a new controller with an empty keymap
func NewKeyboardController() *KeyboardController {
	controller := KeyboardController{}
	controller.receivers = make(map[keyboard.Receiver]keyboard.KeyMap)
	return &controller
}

// Receivers

// AddReceiver adds a keyboard events Receiver to the controller and asociated keymap
func (controller *KeyboardController) AddReceiver(receiver keyboard.Receiver, keymap keyboard.KeyMap) {
	controller.receivers[receiver] = keymap
}

// RemoveReceiver removes the Receiver
func (controller *KeyboardController) RemoveReceiver(receiver keyboard.Receiver) {
	delete(controller.receivers, receiver)
}

// Key events

// KeyDown emits a keyboard keydown event
func (controller *KeyboardController) KeyDown(keycode keyboard.KeyCode) {
	controller.KeyEvent(keycode, keyboard.KeyDown)
}

// KeyUp emits a keyboard keyup event
func (controller *KeyboardController) KeyUp(keycode keyboard.KeyCode) {
	controller.KeyEvent(keycode, keyboard.KeyUp)
}

// KeyEvent emits a keyboard event
func (controller *KeyboardController) KeyEvent(keycode keyboard.KeyCode, eventType int) {
	// For every receiver checks if keycode is mapped
	for receiver, keymap := range controller.receivers {
		keys, ok := keymap[keycode]
		if ok {
			// For each key emit event to receiver
			for _, key := range keys {
				keyevent := keyboard.KeyEvent{
					Event:   device.Event{Type: eventType, Order: device.OrderAfter},
					Key:     key,
					Pressed: eventType == keyboard.KeyDown}
				receiver.ProcessKeyEvent(&keyevent)
			}
		}
	}
}
