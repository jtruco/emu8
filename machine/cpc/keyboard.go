package cpc

import (
	"github.com/jtruco/emu8/device/io/keyboard"
)

// -----------------------------------------------------------------------------
// Amstrad CPC Keyboard
// -----------------------------------------------------------------------------

// Keyboard is the Amstrad CPC
type Keyboard struct {
	rowstates [16]byte
	row       byte
}

// NewKeyboard creates a new keyboard
func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

// State gets current row state
func (keyboard *Keyboard) State() byte {
	return keyboard.rowstates[keyboard.row]
}

// SetRow sets current row
func (keyboard *Keyboard) SetRow(row byte) {
	keyboard.row = row & 0x0f
}

// Device

// Init initializes the keyboard
func (keyboard *Keyboard) Init() {
	for row := 0; row < 16; row++ {
		keyboard.rowstates[row] = 0xff
	}
}

// Reset resets the keyboard
func (keyboard *Keyboard) Reset() { keyboard.Init() }

// Keyboard

// ProcessKeyEvent processes the keyboard events
func (keyboard *Keyboard) ProcessKeyEvent(event *keyboard.KeyEvent) {
	row := event.Key >> 4
	bit := event.Key & 0x07
	mask := byte(1 << bit)
	if event.Pressed {
		keyboard.rowstates[row] &= ^mask
	} else {
		keyboard.rowstates[row] |= mask
	}
}

// -----------------------------------------------------------------------------
// Amstrad CPC Keys, States & Mapping
// -----------------------------------------------------------------------------

// Amstrad CPC Keyboard Keys
const (
	CpcKeyCursorUp     = 0x00 // line 0, bit 0..bit 7
	CpcKeyCursorRight  = 0x01
	CpcKeyCursorDown   = 0x02
	CpcKeyF9           = 0x03
	CpcKeyF6           = 0x04
	CpcKeyF3           = 0x05
	CpcKeyIntro        = 0x06
	CpcKeyFdot         = 0x07
	CpcKeyCursorLeft   = 0x10 // line 1, bit 0..bit 7
	CpcKeyCopy         = 0x11
	CpcKeyF7           = 0x12
	CpcKeyF8           = 0x13
	CpcKeyF5           = 0x14
	CpcKeyF1           = 0x15
	CpcKeyF2           = 0x16
	CpcKeyF0           = 0x17 // line 2, bit 0..bit 7
	CpcKeyClr          = 0x20
	CpcKeyOpenBracket  = 0x21
	CpcKeyReturn       = 0x22
	CpcKeyCloseBracket = 0x23
	CpcKeyF4           = 0x24
	CpcKeyShift        = 0x25
	CpcKeyForwardSlash = 0x26
	CpcKeyControl      = 0x27 // line 3, bit 0.. bit 7
	CpcKeyHat          = 0x30
	CpcKeyMinus        = 0x31
	CpcKeyAt           = 0x32
	CpcKeyP            = 0x33
	CpcKeySemicolon    = 0x34
	CpcKeyColon        = 0x35
	CpcKeyBackslash    = 0x36
	CpcKeyDot          = 0x37 // line 4, bit 0..bit 7
	CpcKey0            = 0x40
	CpcKey9            = 0x41
	CpcKeyO            = 0x42
	CpcKeyI            = 0x43
	CpcKeyL            = 0x44
	CpcKeyK            = 0x45
	CpcKeyM            = 0x46
	CpcKeyComma        = 0x47 // line 5, bit 0..bit 7
	CpcKey8            = 0x50
	CpcKey7            = 0x51
	CpcKeyU            = 0x52
	CpcKeyY            = 0x53
	CpcKeyH            = 0x54
	CpcKeyJ            = 0x55
	CpcKeyN            = 0x56
	CpcKeySpace        = 0x57 // line 6, bit 0..bit 7
	CpcKey6            = 0x60
	CpcKey5            = 0x61
	CpcKeyR            = 0x62
	CpcKeyT            = 0x63
	CpcKeyG            = 0x64
	CpcKeyF            = 0x65
	CpcKeyB            = 0x66
	CpcKeyV            = 0x67 // line 7, bit 0.. bit 7
	CpcKey4            = 0x70
	CpcKey3            = 0x71
	CpcKeyE            = 0x72
	CpcKeyW            = 0x73
	CpcKeyS            = 0x74
	CpcKeyD            = 0x75
	CpcKeyC            = 0x76
	CpcKeyX            = 0x77 // line 8, bit 0.. bit 7
	CpcKey1            = 0x80
	CpcKey2            = 0x81
	CpcKeyEsc          = 0x82
	CpcKeyQ            = 0x83
	CpcKeyTab          = 0x84
	CpcKeyA            = 0x85
	CpcKeyCapsLock     = 0x86
	CpcKeyZ            = 0x87 // line 9, bit 7..bit 0
	CpcKeyJoyUp        = 0x90
	CpcKeyJoyDown      = 0x91
	CpcKeyJoyLeft      = 0x92
	CpcKeyJoyRight     = 0x93
	CpcKeyJoyFire1     = 0x94
	CpcKeyJoyFire2     = 0x95
	CpcKeySpare        = 0x96
	CpcKeyDel          = 0x97
)

// CPC Keyboard map
var cpcKeyboardMap = map[keyboard.KeyCode][]keyboard.Key{
	// alphanum
	keyboard.Key0: {CpcKey0},
	keyboard.Key1: {CpcKey1},
	keyboard.Key2: {CpcKey2},
	keyboard.Key3: {CpcKey3},
	keyboard.Key4: {CpcKey4},
	keyboard.Key5: {CpcKey5},
	keyboard.Key6: {CpcKey6},
	keyboard.Key7: {CpcKey7},
	keyboard.Key8: {CpcKey8},
	keyboard.Key9: {CpcKey9},
	keyboard.KeyA: {CpcKeyA},
	keyboard.KeyB: {CpcKeyB},
	keyboard.KeyC: {CpcKeyC},
	keyboard.KeyD: {CpcKeyD},
	keyboard.KeyE: {CpcKeyE},
	keyboard.KeyF: {CpcKeyF},
	keyboard.KeyG: {CpcKeyG},
	keyboard.KeyH: {CpcKeyH},
	keyboard.KeyI: {CpcKeyI},
	keyboard.KeyJ: {CpcKeyJ},
	keyboard.KeyK: {CpcKeyK},
	keyboard.KeyL: {CpcKeyL},
	keyboard.KeyM: {CpcKeyM},
	keyboard.KeyN: {CpcKeyN},
	keyboard.KeyO: {CpcKeyO},
	keyboard.KeyP: {CpcKeyP},
	keyboard.KeyQ: {CpcKeyQ},
	keyboard.KeyR: {CpcKeyR},
	keyboard.KeyS: {CpcKeyS},
	keyboard.KeyT: {CpcKeyT},
	keyboard.KeyU: {CpcKeyU},
	keyboard.KeyV: {CpcKeyV},
	keyboard.KeyW: {CpcKeyW},
	keyboard.KeyX: {CpcKeyX},
	keyboard.KeyY: {CpcKeyY},
	keyboard.KeyZ: {CpcKeyZ},
	// special
	keyboard.KeySpace:        {CpcKeySpace},
	keyboard.KeyComma:        {CpcKeyComma},
	keyboard.KeyPeriod:       {CpcKeyDot},
	keyboard.KeySemicolon:    {CpcKeyColon},
	keyboard.KeyMinus:        {CpcKeyMinus},
	keyboard.KeyEquals:       {CpcKeyHat},
	keyboard.KeyLeftBracket:  {CpcKeyAt},
	keyboard.KeyRightBracket: {CpcKeyOpenBracket},
	keyboard.KeyTab:          {CpcKeyTab},
	keyboard.KeyReturn:       {CpcKeyReturn},
	keyboard.KeyBackspace:    {CpcKeyDel},
	keyboard.KeyEscape:       {CpcKeyEsc},
	// cursors
	keyboard.KeyUp:    {CpcKeyCursorUp},
	keyboard.KeyDown:  {CpcKeyCursorDown},
	keyboard.KeyLeft:  {CpcKeyCursorLeft},
	keyboard.KeyRight: {CpcKeyCursorRight},
	// keypad
	keyboard.KeyPad0:      {CpcKeyF0},
	keyboard.KeyPad1:      {CpcKeyF1},
	keyboard.KeyPad2:      {CpcKeyF2},
	keyboard.KeyPad3:      {CpcKeyF3},
	keyboard.KeyPad4:      {CpcKeyF4},
	keyboard.KeyPad5:      {CpcKeyF5},
	keyboard.KeyPad6:      {CpcKeyF6},
	keyboard.KeyPad7:      {CpcKeyF7},
	keyboard.KeyPad8:      {CpcKeyF8},
	keyboard.KeyPad9:      {CpcKeyF9},
	keyboard.KeyPadEnter:  {CpcKeyIntro},
	keyboard.KeyPadPeriod: {CpcKeyFdot},
	keyboard.KeyDelete:    {CpcKeyClr},
	// shift & control
	keyboard.KeyLShift:   {CpcKeyShift},
	keyboard.KeyRShift:   {CpcKeyShift},
	keyboard.KeyLCtrl:    {CpcKeyControl},
	keyboard.KeyRCtrl:    {CpcKeyControl},
	keyboard.KeyCapsLock: {CpcKeyCapsLock},
	// other
	keyboard.KeyLAlt:       {CpcKeyCopy},
	keyboard.KeyRAlt:       {CpcKeyCopy},
	keyboard.KeyGrave:      {CpcKeyForwardSlash},
	keyboard.KeySlash:      {CpcKeyBackslash},
	keyboard.KeyApostrophe: {CpcKeySemicolon},
	keyboard.KeyBackSlash:  {CpcKeyCloseBracket},
}
