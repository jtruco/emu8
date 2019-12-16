package sdl

import (
	"log"
	"time"

	"github.com/jtruco/emu8/config"
	"github.com/jtruco/emu8/emulator"
	"github.com/veandco/go-sdl2/sdl"
)

// App is the SDL application
type App struct {
	title    string
	config   *config.Config
	video    *Video
	audio    *Audio
	emulator *emulator.Emulator
}

// NewApp creates a new application
func NewApp() *App {
	app := &App{}
	app.config = config.Get()
	app.video = NewVideo(app)
	app.audio = NewAudio(app)
	return app
}

// Init the SDL app
func (app *App) Init() bool {
	// init sdl
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		log.Println("Error initializing SDL : " + err.Error())
		return false
	}
	// init emulator
	if app.emulator = emulator.FromMachine(app.config.MachineModel); app.emulator == nil {
		log.Println("Error initializing emulator : machine model is not valid.")
		return false
	}
	app.emulator.Controller().Video().SetRenderer(app.video)
	app.emulator.Controller().Audio().SetPlayer(app.audio)
	app.emulator.Init()
	if app.config.Snapfile != "" {
		app.emulator.Machine().LoadFile(app.config.Snapfile)
	}
	// init SDL video output
	if !app.video.Init() {
		return false
	}
	// init SDL audio
	if !app.audio.Init() {
		return false
	}
	return true
}

// Run the SDL app
func (app *App) Run() {
	const sleep = 20 * time.Millisecond
	app.emulator.Start()
	running := true
	for running {
		count := 0
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			count++
			switch e := event.(type) {

			case *sdl.QuitEvent:
				running = false

			case *sdl.KeyboardEvent:
				// check function keys
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_F5:
						app.emulator.Reset()
					}
				}
				// send 1erst key event to emulator
				if e.Repeat == 0 {
					if e.Type == sdl.KEYDOWN {
						app.emulator.Controller().Keyboard().KeyDown(int(e.Keysym.Scancode))
					} else {
						app.emulator.Controller().Keyboard().KeyUp(int(e.Keysym.Scancode))
					}
				}
			}
		}
		// sleep
		if count == 0 {
			time.Sleep(sleep)
		}
	}
	// stop emulator
	app.emulator.Stop()
}

// End the SDL app
func (app *App) End() {
	sdl.Quit()
}
