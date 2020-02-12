package controller

import (
	"sync"

	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/io/joystick"
)

// -----------------------------------------------------------------------------
// Joystick Controller
// -----------------------------------------------------------------------------

// JoystickController is the emulator joystick(s) controller
type JoystickController struct {
	receivers  map[byte]joystick.Joystick // Joystick devices mapped by ID
	eventQueue []*joystick.JoyEvent       // Events to dispatch
	mtx        sync.Mutex                 // Sync
}

// NewJoystickController creates a new controller
func NewJoystickController() *JoystickController {
	controller := JoystickController{}
	controller.receivers = make(map[byte]joystick.Joystick)
	controller.eventQueue = make([]*joystick.JoyEvent, 0, 5)
	return &controller
}

// Receivers

// AddReceiver adds a joystick receiver to the controller at associated ID
func (controller *JoystickController) AddReceiver(receiver joystick.Joystick, id byte) {
	controller.receivers[id] = receiver
}

// RemoveReceiver removes the joystick receiver associated
func (controller *JoystickController) RemoveReceiver(receiver joystick.Joystick) {
	for id, joy := range controller.receivers {
		if joy == receiver {
			delete(controller.receivers, id)
		}
	}
}

// Events

// AxisEvent emits a joystick axis event
func (controller *JoystickController) AxisEvent(id byte, axis byte, value byte) {
	joyevent := joystick.JoyEvent{
		Event:     device.Event{Type: joystick.EventJoyAxis, Order: device.OrderAfter},
		ID:        id,
		Axis:      axis,
		AxisValue: value}
	controller.appendEvent(&joyevent)
}

// ButtonEvent emits a joystick button event
func (controller *JoystickController) ButtonEvent(id byte, button byte, state byte) {
	joyevent := joystick.JoyEvent{
		Event:       device.Event{Type: joystick.EventJoyBotton, Order: device.OrderAfter},
		ID:          id,
		Button:      button,
		ButtonState: state}
	controller.appendEvent(&joyevent)
}

// appendEvent adds event to queue
func (controller *JoystickController) appendEvent(joyEvent *joystick.JoyEvent) {
	controller.mtx.Lock()
	defer controller.mtx.Unlock()

	controller.eventQueue = append(controller.eventQueue, joyEvent)
}

// Flush flushes joystick events
func (controller *JoystickController) Flush() {
	controller.mtx.Lock()
	defer controller.mtx.Unlock()

	for _, e := range controller.eventQueue {
		controller.emitEvent(e)
	}
	controller.eventQueue = controller.eventQueue[:0]
}

// processEvent process a joystick event
func (controller *JoystickController) emitEvent(joyEvent *joystick.JoyEvent) {
	if joy, ok := controller.receivers[joyEvent.ID]; ok {
		switch joyEvent.Event.Type {
		case joystick.EventJoyAxis:
			joy.SetAxis(joyEvent.Axis, joyEvent.AxisValue)
		case joystick.EventJoyBotton:
			joy.SetButton(joyEvent.Button, joyEvent.ButtonState)
		}
	}
}
