// Package io contains io controller components
package io

import (
	"sync"

	"github.com/jtruco/emu8/emulator/device/io/keyboard"
)

// -----------------------------------------------------------------------------
// Keyboard Controller
// -----------------------------------------------------------------------------

// KeyboardController is the emulator keyboard controller
type KeyboardController struct {
	receivers  map[keyboard.Receiver]keyboard.KeyMap // Keyboard receiver devices
	eventQueue []keyEvent                            // Keyboard event queue
	mtx        sync.Mutex                            // Sync
}

// Keyboard key event
type keyEvent struct {
	Keycode   keyboard.KeyCode
	EventType int
}

// NewKeyboardController creates a new controller with an empty keymap
func NewKeyboardController() *KeyboardController {
	controller := new(KeyboardController)
	controller.receivers = make(map[keyboard.Receiver]keyboard.KeyMap)
	controller.eventQueue = make([]keyEvent, 0, 5)
	return controller
}

// Receivers

// AddReceiver adds a keyboard events Receiver to the controller and associated keymap
func (controller *KeyboardController) AddReceiver(receiver keyboard.Receiver) {
	controller.receivers[receiver] = receiver.KeyMap()
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
	controller.mtx.Lock()
	defer controller.mtx.Unlock()

	event := keyEvent{keycode, eventType}
	controller.eventQueue = append(controller.eventQueue, event)
}

// Flush flushes keyboard event queue
func (controller *KeyboardController) Flush() {
	controller.mtx.Lock()
	defer controller.mtx.Unlock()

	for _, e := range controller.eventQueue {
		controller.emitEvent(e)
	}
	controller.eventQueue = controller.eventQueue[:0]
}

// emitEvent emits a keyboard event
func (controller *KeyboardController) emitEvent(e keyEvent) {
	// For every receiver checks if keycode is mapped
	for receiver, keymap := range controller.receivers {
		keys, ok := keymap[e.Keycode]
		if ok {
			pressed := e.EventType == keyboard.KeyDown
			// For each key emit event to receiver
			for _, key := range keys {
				receiver.ProcessKey(key, pressed)
			}
		}
	}
}
