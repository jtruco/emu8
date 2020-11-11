package spectrum

import "github.com/jtruco/emu8/emulator/device/io/keyboard"

// -----------------------------------------------------------------------------
// ZX Spectrum Keyboard
// -----------------------------------------------------------------------------

// Keyboard is the ZX Spectrum Keyboard
type Keyboard struct {
	rowstates [8]byte // The keyboard row states
}

// NewKeyboard creates a new keyboard
func NewKeyboard() *Keyboard {
	return new(Keyboard)
}

// Device

// Init initializes the keyboard
func (keyboard *Keyboard) Init() {
	for row := 0; row < 8; row++ {
		keyboard.rowstates[row] = 0xff
	}
}

// Reset resets the keyboard
func (keyboard *Keyboard) Reset() { keyboard.Init() }

// Keyboard

// ProcessKeyEvent processes the keyboard events
func (keyboard *Keyboard) ProcessKeyEvent(event *keyboard.KeyEvent) {
	state, ok := keyStates[event.Key]
	if ok {
		if event.Pressed {
			keyboard.rowstates[state[0]] &= ^(state[1])
		} else {
			keyboard.rowstates[state[0]] |= state[1]
		}
	}
}

// -----------------------------------------------------------------------------
// ZX Spectrum Keys, States & Mapping
// -----------------------------------------------------------------------------

// ZX Spectrum Keys
const (
	ZxKey1 = iota
	ZxKey2
	ZxKey3
	ZxKey4
	ZxKey5
	ZxKey6
	ZxKey7
	ZxKey8
	ZxKey9
	ZxKey0

	ZxKeyQ
	ZxKeyW
	ZxKeyE
	ZxKeyR
	ZxKeyT
	ZxKeyY
	ZxKeyU
	ZxKeyI
	ZxKeyO
	ZxKeyP

	ZxKeyA
	ZxKeyS
	ZxKeyD
	ZxKeyF
	ZxKeyG
	ZxKeyH
	ZxKeyJ
	ZxKeyK
	ZxKeyL
	ZxKeyEnter

	ZxKeyCapsShift
	ZxKeyZ
	ZxKeyX
	ZxKeyC
	ZxKeyV
	ZxKeyB
	ZxKeyN
	ZxKeyM
	ZxKeySymbolShift
	ZxKeySpace
)

// keyStates ZX Spectrum Keyboard states
var keyStates = map[keyboard.Key][2]byte{
	ZxKey1: {3, 0x01},
	ZxKey2: {3, 0x02},
	ZxKey3: {3, 0x04},
	ZxKey4: {3, 0x08},
	ZxKey5: {3, 0x10},
	ZxKey6: {4, 0x10},
	ZxKey7: {4, 0x08},
	ZxKey8: {4, 0x04},
	ZxKey9: {4, 0x02},
	ZxKey0: {4, 0x01},

	ZxKeyQ: {2, 0x01},
	ZxKeyW: {2, 0x02},
	ZxKeyE: {2, 0x04},
	ZxKeyR: {2, 0x08},
	ZxKeyT: {2, 0x10},
	ZxKeyY: {5, 0x10},
	ZxKeyU: {5, 0x08},
	ZxKeyI: {5, 0x04},
	ZxKeyO: {5, 0x02},
	ZxKeyP: {5, 0x01},

	ZxKeyA:     {1, 0x01},
	ZxKeyS:     {1, 0x02},
	ZxKeyD:     {1, 0x04},
	ZxKeyF:     {1, 0x08},
	ZxKeyG:     {1, 0x10},
	ZxKeyH:     {6, 0x10},
	ZxKeyJ:     {6, 0x08},
	ZxKeyK:     {6, 0x04},
	ZxKeyL:     {6, 0x02},
	ZxKeyEnter: {6, 0x01},

	ZxKeyCapsShift:   {0, 0x01},
	ZxKeyZ:           {0, 0x02},
	ZxKeyX:           {0, 0x04},
	ZxKeyC:           {0, 0x08},
	ZxKeyV:           {0, 0x10},
	ZxKeyB:           {7, 0x10},
	ZxKeyN:           {7, 0x08},
	ZxKeyM:           {7, 0x04},
	ZxKeySymbolShift: {7, 0x02},
	ZxKeySpace:       {7, 0x01},
}

// ZX Keyboard map
var zxKeyboardMap = map[keyboard.KeyCode][]keyboard.Key{
	// standar mapping
	keyboard.Key0: {ZxKey0},
	keyboard.Key1: {ZxKey1},
	keyboard.Key2: {ZxKey2},
	keyboard.Key3: {ZxKey3},
	keyboard.Key4: {ZxKey4},
	keyboard.Key5: {ZxKey5},
	keyboard.Key6: {ZxKey6},
	keyboard.Key7: {ZxKey7},
	keyboard.Key8: {ZxKey8},
	keyboard.Key9: {ZxKey9},

	keyboard.KeyA: {ZxKeyA},
	keyboard.KeyB: {ZxKeyB},
	keyboard.KeyC: {ZxKeyC},
	keyboard.KeyD: {ZxKeyD},
	keyboard.KeyE: {ZxKeyE},
	keyboard.KeyF: {ZxKeyF},
	keyboard.KeyG: {ZxKeyG},
	keyboard.KeyH: {ZxKeyH},
	keyboard.KeyI: {ZxKeyI},
	keyboard.KeyJ: {ZxKeyJ},
	keyboard.KeyK: {ZxKeyK},
	keyboard.KeyL: {ZxKeyL},
	keyboard.KeyM: {ZxKeyM},
	keyboard.KeyN: {ZxKeyN},
	keyboard.KeyO: {ZxKeyO},
	keyboard.KeyP: {ZxKeyP},
	keyboard.KeyQ: {ZxKeyQ},
	keyboard.KeyR: {ZxKeyR},
	keyboard.KeyS: {ZxKeyS},
	keyboard.KeyT: {ZxKeyT},
	keyboard.KeyU: {ZxKeyU},
	keyboard.KeyV: {ZxKeyV},
	keyboard.KeyW: {ZxKeyW},
	keyboard.KeyX: {ZxKeyX},
	keyboard.KeyY: {ZxKeyY},
	keyboard.KeyZ: {ZxKeyZ},

	keyboard.KeyReturn: {ZxKeyEnter},
	keyboard.KeySpace:  {ZxKeySpace},
	keyboard.KeyLShift: {ZxKeyCapsShift},
	keyboard.KeyRShift: {ZxKeyCapsShift},
	keyboard.KeyLCtrl:  {ZxKeySymbolShift},
	keyboard.KeyRCtrl:  {ZxKeySymbolShift},

	// cursors
	keyboard.KeyLeft:  {ZxKeyCapsShift, ZxKey5},
	keyboard.KeyDown:  {ZxKeyCapsShift, ZxKey6},
	keyboard.KeyUp:    {ZxKeyCapsShift, ZxKey7},
	keyboard.KeyRight: {ZxKeyCapsShift, ZxKey8},

	// keypad
	keyboard.KeyPad1:        {ZxKey1},
	keyboard.KeyPad2:        {ZxKey2},
	keyboard.KeyPad3:        {ZxKey3},
	keyboard.KeyPad4:        {ZxKey4},
	keyboard.KeyPad5:        {ZxKey5},
	keyboard.KeyPad6:        {ZxKey6},
	keyboard.KeyPad7:        {ZxKey7},
	keyboard.KeyPad8:        {ZxKey8},
	keyboard.KeyPad9:        {ZxKey9},
	keyboard.KeyPad0:        {ZxKey0},
	keyboard.KeyPadMultiply: {ZxKeySymbolShift, ZxKeyB},
	keyboard.KeyPadDivide:   {ZxKeySymbolShift, ZxKeyV},
	keyboard.KeyPadPlus:     {ZxKeySymbolShift, ZxKeyK},
	keyboard.KeyPadMinus:    {ZxKeySymbolShift, ZxKeyJ},
	keyboard.KeyPadEnter:    {ZxKeyEnter},

	// other keyboard maps
	keyboard.KeyBackspace: {ZxKeyCapsShift, ZxKey0},
	keyboard.KeyEscape:    {ZxKeyCapsShift, ZxKey1},
}
