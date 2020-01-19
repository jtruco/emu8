package sdl

import (
	"log"

	"github.com/jtruco/emu8/device/audio"
	"github.com/veandco/go-sdl2/sdl"
)

// Audio the SDL audio engine
type Audio struct {
	app *App // The SDL app
}

// NewAudio the SDL audio
func NewAudio(app *App) *Audio {
	audio := &Audio{}
	audio.app = app
	return audio
}

// Init the SDL audio
func (audio *Audio) Init() bool {
	var want, spec sdl.AudioSpec
	want.Freq = 48000
	want.Format = sdl.AUDIO_S16LSB
	want.Channels = 1
	want.Samples = 1024
	err := sdl.OpenAudio(&want, &spec)
	if err != nil {
		log.Println("Error initializing SDL audio : " + err.Error())
		return false

	}
	sdl.PauseAudio(false)
	return true
}

// Play plays the audio buffer
func (audio *Audio) Play(buffer *audio.Buffer) {
	sdl.QueueAudio(1, buffer.Data())
}