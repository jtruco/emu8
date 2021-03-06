package keyboard

// -----------------------------------------------------------------------------
// Keyboard KeyCodes & Machine Keys
// -----------------------------------------------------------------------------

// KeyCode is a keyboard key code
type KeyCode = int

// Key is a machine keyboard key
type Key = int

// KeyMap is a keyboard code to machine mapping
type KeyMap = map[KeyCode][]Key

// Keyboard keys avaible to map machine keyboards
// These codes are equals to SDL scancodes and USB standards
const (
	// Unknown
	KeyUnknown = 0
	// Alpha
	KeyA = 4
	KeyB = 5
	KeyC = 6
	KeyD = 7
	KeyE = 8
	KeyF = 9
	KeyG = 10
	KeyH = 11
	KeyI = 12
	KeyJ = 13
	KeyK = 14
	KeyL = 15
	KeyM = 16
	KeyN = 17
	KeyO = 18
	KeyP = 19
	KeyQ = 20
	KeyR = 21
	KeyS = 22
	KeyT = 23
	KeyU = 24
	KeyV = 25
	KeyW = 26
	KeyX = 27
	KeyY = 28
	KeyZ = 29
	// Numeric
	Key1 = 30
	Key2 = 31
	Key3 = 32
	Key4 = 33
	Key5 = 34
	Key6 = 35
	Key7 = 36
	Key8 = 37
	Key9 = 38
	Key0 = 39
	// Special
	KeyReturn       = 40
	KeyEscape       = 41
	KeyBackspace    = 42
	KeyTab          = 43
	KeySpace        = 44
	KeyMinus        = 45
	KeyEquals       = 46
	KeyLeftBracket  = 47
	KeyRightBracket = 48
	KeyBackSlash    = 49
	KeyNonUsHash    = 50
	KeySemicolon    = 51
	KeyApostrophe   = 52
	KeyGrave        = 53
	KeyComma        = 54
	KeyPeriod       = 55
	KeySlash        = 56
	KeyCapsLock     = 57
	// Function
	KeyF1  = 58
	KeyF2  = 59
	KeyF3  = 60
	KeyF4  = 61
	KeyF5  = 62
	KeyF6  = 63
	KeyF7  = 64
	KeyF8  = 65
	KeyF9  = 66
	KeyF10 = 67
	KeyF11 = 68
	KeyF12 = 69
	// Special function
	KeyPrintScreen = 70
	KeyScrollLock  = 71
	KeyPause       = 72
	// Cursor & Edit
	KeyInsert   = 73
	KeyHome     = 74
	KeyPageUp   = 75
	KeyDelete   = 76
	KeyEnd      = 77
	KeyPageDown = 78
	KeyRight    = 79
	KeyLeft     = 80
	KeyDown     = 81
	KeyUp       = 82
	// Keypad
	KeyPadDivide   = 84
	KeyPadMultiply = 85
	KeyPadMinus    = 86
	KeyPadPlus     = 87
	KeyPadEnter    = 88
	KeyPad1        = 89
	KeyPad2        = 90
	KeyPad3        = 91
	KeyPad4        = 92
	KeyPad5        = 93
	KeyPad6        = 94
	KeyPad7        = 95
	KeyPad8        = 96
	KeyPad9        = 97
	KeyPad0        = 98
	KeyPadPeriod   = 99
	// Meta
	KeyLCtrl  = 224
	KeyLShift = 225
	KeyLAlt   = 226
	KeyLGui   = 227
	KeyRCtrl  = 228
	KeyRShift = 229
	KeyRAlt   = 230
	KeyRGui   = 231
)
