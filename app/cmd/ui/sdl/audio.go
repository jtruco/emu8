package sdl

import (
	"log"

	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/veandco/go-sdl2/sdl"
)

// Audio the SDL audio engine
type Audio struct {
	config *config.AudioConfig // Audio configuration
}

// NewAudio the SDL audio
func NewAudio(config *config.Config) *Audio {
	audio := new(Audio)
	audio.config = &config.Audio
	return audio
}

// Init the SDL audio
func (audio *Audio) Init() bool {
	var want, spec sdl.AudioSpec
	want.Freq = int32(audio.config.Frequency)
	want.Format = sdl.AUDIO_S16LSB
	want.Channels = 2 // stereo
	want.Samples = 1024
	err := sdl.OpenAudio(&want, &spec)
	if err != nil {
		log.Println("Error initializing SDL audio : " + err.Error())
		return false

	}
	sdl.PauseAudio(false)
	return true
}

// Close closes audio resources
func (audio *Audio) Close() {
	sdl.CloseAudio()
}

// Play plays the audio buffer
func (audio *Audio) Play(buffer *audio.Buffer) {
	if !audio.config.Mute {
		sdl.QueueAudio(1, buffer.Data())
	}
}
