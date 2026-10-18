package io

import (
	"testing"

	"github.com/jtruco/emu8/emulator/device/io/keyboard"
)

// testReceiver records the processed keys
type testReceiver struct {
	keymap  keyboard.KeyMap
	pressed map[keyboard.Key]bool
}

func (r *testReceiver) KeyMap() keyboard.KeyMap { return r.keymap }

func (r *testReceiver) ProcessKey(key keyboard.Key, pressed bool) {
	r.pressed[key] = pressed
}

func TestKeyboardControllerEvents(t *testing.T) {
	const keyFire, keyShift, keyLeft = 1, 2, 3
	receiver := &testReceiver{
		keymap: keyboard.KeyMap{
			keyboard.KeySpace: {keyFire},
			keyboard.KeyDown:  {keyShift, keyLeft}, // cursor key down : 2 machine keys
		},
		pressed: map[keyboard.Key]bool{},
	}
	controller := NewKeyboardController()
	controller.AddReceiver(receiver)

	controller.KeyDown(keyboard.KeySpace)
	controller.KeyDown(keyboard.KeyDown)
	controller.KeyDown(keyboard.KeyA) // not mapped
	if len(receiver.pressed) != 0 {
		t.Fatal("keys processed before Flush")
	}
	controller.Flush()
	for _, key := range []keyboard.Key{keyFire, keyShift, keyLeft} {
		if !receiver.pressed[key] {
			t.Errorf("key %d not pressed after KeyDown", key)
		}
	}

	controller.KeyUp(keyboard.KeyDown)
	controller.Flush()
	if !receiver.pressed[keyFire] || receiver.pressed[keyShift] || receiver.pressed[keyLeft] {
		t.Errorf("after cursor KeyUp: %v, want only fire pressed", receiver.pressed)
	}

	// queued events are processed in order
	controller.KeyDown(keyboard.KeyUp)
	controller.KeyUp(keyboard.KeySpace)
	controller.KeyDown(keyboard.KeySpace)
	controller.Flush()
	if !receiver.pressed[keyFire] {
		t.Error("fire released, want pressed by the last queued event")
	}
}
