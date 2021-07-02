package sdl

import (
	"log"

	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/veandco/go-sdl2/sdl"
)

// Audio the SDL audio engine
type Audio struct {
	app    *App                // SDL Application
	config *config.AudioConfig // Audio configuration
	device audio.Audio         // Machine audio device
}

// NewAudio the SDL audio
func NewAudio(app *App) *Audio {
	audio := new(Audio)
	audio.app = app
	audio.config = &app.config.Audio
	return audio
}

// Init the SDL audio
func (audio *Audio) Init(device audio.Audio) bool {
	audio.device = device
	var want, spec sdl.AudioSpec
	want.Freq = int32(audio.config.Frequency)
	want.Format = sdl.AUDIO_S16LSB
	want.Channels = 1 // mono
	want.Samples = 1024
	err := sdl.OpenAudio(&want, &spec)
	if err != nil {
		log.Println("SDL : Error initializing SDL audio:", err.Error())
		return false

	}
	sdl.PauseAudio(false)
	log.Println("SDL : Audio initialized:", want.Freq, "Hz")
	return true
}

// Close closes audio resources
func (audio *Audio) Close() {
	log.Println("SDL : Closing audio resources")
	sdl.CloseAudio()
}

// Play plays the audio buffer
func (audio *Audio) Play(buffer *audio.Buffer) {
	if !audio.config.Mute {
		if audio.app.async {
			go audio.onPlay(buffer)
		} else {
			audio.onPlay(buffer)
		}
	}
}

func (audio *Audio) onPlay(buffer *audio.Buffer) {
	sdl.QueueAudio(1, buffer.Data())
}
