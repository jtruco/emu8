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
	running  bool
}

// NewApp creates a new application
func NewApp() *App {
	app := &App{}
	app.config = config.Get()
	app.video = NewVideo(app)
	app.audio = NewAudio(app)
	return app
}

// Init the SDL App
func (app *App) Init() bool {
	// init sdl
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		log.Println("Error initializing SDL : " + err.Error())
		return false
	}
	// init emulator
	emulator := emulator.FromModel(app.config.MachineModel)
	if emulator == nil {
		log.Println("Error initializing emulator : machine model is not valid.")
		return false
	}
	emulator.Controller().Video().SetRenderer(app.video)
	emulator.Controller().Audio().SetPlayer(app.audio)
	emulator.Init()
	if app.config.FileName != "" {
		emulator.Machine().LoadFile(app.config.FileName)
	}
	app.emulator = emulator
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

// Run the SDL App
func (app *App) Run() {
	const sleep = 20 * time.Millisecond // 50 Hz
	app.emulator.Start()
	app.running = true
	for app.running {
		count := 0
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			count++
			switch e := event.(type) {
			case *sdl.QuitEvent:
				app.running = false
			case *sdl.KeyboardEvent:
				app.processKeyboard(e)
			}
		}
		if count == 0 {
			time.Sleep(sleep)
		}
	}
	app.emulator.Stop()
}

// End the SDL App
func (app *App) End() {
	sdl.Quit()
}

func (app *App) processKeyboard(e *sdl.KeyboardEvent) {
	if e.Repeat > 0 {
		return
	}
	captured := false
	// check function keys
	if e.Type == sdl.KEYDOWN {
		captured = true
		switch e.Keysym.Sym {
		// Emulator
		case sdl.K_ESCAPE:
			app.running = false
		case sdl.K_F5:
			app.emulator.Reset()
		// Tape Drive
		case sdl.K_F6:
			app.emulator.Controller().Tape().Drive().Play()
		case sdl.K_F7:
			app.emulator.Controller().Tape().Drive().Stop()
		case sdl.K_F8:
			app.emulator.Controller().Tape().Drive().Rewind()
		// Emulator
		case sdl.K_F9:
			app.emulator.Start()
		case sdl.K_F10:
			app.emulator.Stop()
		// UI
		case sdl.K_F11:
			app.video.ToggleFullscreen()
		default:
			captured = false
		}
	}
	if !captured {
		// send key event to emulator
		if e.Type == sdl.KEYDOWN {
			app.emulator.Controller().Keyboard().KeyDown(int(e.Keysym.Scancode))
		} else {
			app.emulator.Controller().Keyboard().KeyUp(int(e.Keysym.Scancode))
		}
	}
}
