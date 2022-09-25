package sdl

import (
	"log"

	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/veandco/go-sdl2/sdl"
)

// Audio the SDL audio engine
type Audio struct {
	Frequency int32 // Audio frequency
	Channels  int32 // Audio channels
	Mute      bool  // Mute audio
	Async     bool  // Asynchronous mode
}

// NewAudio the SDL audio
func NewAudio() *Audio {
	audio := new(Audio)
	audio.Frequency = 48000
	audio.Channels = 1 // mono
	return audio
}

// Init the SDL audio
func (audio *Audio) Init() bool {
	var want, spec sdl.AudioSpec
	want.Freq = audio.Frequency
	want.Format = sdl.AUDIO_S16LSB
	want.Channels = uint8(audio.Channels)
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
	if !audio.Mute {
		if audio.Async {
			go audio.onPlay(buffer)
		} else {
			audio.onPlay(buffer)
		}
	}
}

func (audio *Audio) onPlay(buffer *audio.Buffer) {
	sdl.QueueAudio(1, buffer.Data())
}
