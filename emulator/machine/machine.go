// Package machine contains 8bit machines and components
package machine

import (
	"fmt"
	"time"

	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/jtruco/emu8/emulator/device/io/joystick"
	"github.com/jtruco/emu8/emulator/device/io/keyboard"
	"github.com/jtruco/emu8/emulator/device/io/tape"
	"github.com/jtruco/emu8/emulator/device/memory"
	"github.com/jtruco/emu8/emulator/device/video"
)

// -----------------------------------------------------------------------------
// Machine
// -----------------------------------------------------------------------------

// Machine is a 8bit machine
type Machine interface {
	Init() error                    // Init initializes the machine and loads its ROMs
	Reset()                         // Reset resets the machine (ROMs are kept)
	Config() *Config                // Config gets the machine configuration
	Clock() device.Clock            // Clock the machine main clock
	Components() *device.Components // Components the machine components
	InitControl(Control)            // InitControl connects the machine to the emulator controller
	Emulate()                       // Emulate one machine step
	BeginFrame()                    // BeginFrame begins the frame of the machine and its devices
	EndFrame()                      // EndFrame ends the frame of the machine and its devices (video and audio output ready)
	LoadState(State) error          // LoadState loads machine state
	SaveState() State               // SaveState saves machine state
}

// Control is the machine control interface
type Control interface {
	// Device binding
	BindVideo(video.Video)          // BindVideo sets the video device
	BindAudio(audio.Audio)          // BindAudio sets the audio device
	BindKeyboard(keyboard.Keyboard) // BindKeyboard adds a keyboard device
	BindJoystick(joystick.Joystick) // BindJoystick adds a joystick device
	BindTapeDrive(*tape.Drive)      // BindTapeDrive sets the tape drive
	// File management
	LoadROM(string) ([]byte, error)    // Loads a ROM file
	RegisterSnapshot(string)           // RegisterSnapshot adds a snapshot format
	RegisterTape(string, tape.Builder) // RegisterTape ads a tape format and its builder
}

// LoadROM loads a ROM file into a read only memory bank. The file must fill
// the bank; bigger files are truncated.
func LoadROM(control Control, name string, bank *memory.Bank) error {
	data, err := control.LoadROM(name)
	if err != nil {
		return fmt.Errorf("ROM %s : %w", name, err)
	}
	if size := int(bank.Size()); len(data) < size {
		return fmt.Errorf("ROM %s : %d bytes, want %d", name, len(data), size)
	}
	bank.Load(0, data[:bank.Size()])
	return nil
}

// DefaultAudioFrequency is the default audio output frequency (48 KHz)
const DefaultAudioFrequency = 48000

// Settings are the machine construction settings
type Settings struct {
	ModelID        string // Model ID or name (e.g. "zx48k", "cpc464"), case insensitive
	AudioFrequency int    // Audio output frequency in Hz (0 : DefaultAudioFrequency)
	Options        string // Machine specific options (e.g. CPC "es", "fr" ROMs)
}

// Frequency returns the audio output frequency, or the default if not set
func (settings Settings) Frequency() int {
	if settings.AudioFrequency > 0 {
		return settings.AudioFrequency
	}
	return DefaultAudioFrequency
}

// Config is the machine configuration : its construction settings, and the
// model information and timings set by the machine
type Config struct {
	Settings               // Construction settings (audio frequency, options)
	Name     string        // Machine model name
	Model    int           // Machine model (internal)
	TStates  int           // TStates per frame
	Fps      int           // Frames per second
	Duration time.Duration // Duration of a frame
}

// SetTimings sets machine timings (TStates, Fps, Duration)
func (config *Config) SetTimings(tstates, fps int) {
	config.TStates = tstates
	config.Fps = fps
	config.Duration = time.Duration(1e9 / fps)
}

// State contains the serialized machine state
type State struct {
	Format string // Format extension
	Data   []byte // Serialized data
}
