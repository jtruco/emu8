package controller

import (
	"github.com/jtruco/emu8/device/audio"
	"github.com/jtruco/emu8/device/io/keyboard"
	"github.com/jtruco/emu8/device/video"
)

// Controller is the emulator controller
type Controller interface {
	FileManager() *FileManager
	Keyboard() *keyboard.Controller
	Video() *video.Controller
	Audio() *audio.Controller
}
