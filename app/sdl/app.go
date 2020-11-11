package main

import (
	"log"
	"time"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/veandco/go-sdl2/sdl"
)

const loopSleep = 20 * time.Millisecond // 50 Hz

// App is the SDL application
type App struct {
	config   *config.Config
	video    *Video
	audio    *Audio
	emulator *emulator.Emulator
	control  controller.Controller
	running  bool
}

// NewApp creates a new application
func NewApp() *App {
	app := new(App)
	app.config = config.Get()
	app.video = NewVideo(app)
	app.audio = NewAudio(app)
	return app
}

// Init the SDL App
func (app *App) Init(emu *emulator.Emulator) bool {
	// init sdl
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_JOYSTICK); err != nil {
		log.Println("Error initializing SDL : " + err.Error())
		return false
	}
	// open sdl joystick (only one supported)
	if sdl.NumJoysticks() > 0 {
		joystick := sdl.JoystickOpen(0)
		if joystick != nil {
			log.Println("Joystick found ...", sdl.JoystickNameForIndex(0))
		}
	}
	// init emulator
	app.emulator = emu
	control := emu.Controller()
	control.Video().SetRenderer(app.video)
	control.Audio().SetPlayer(app.audio)
	app.control = control
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

	// init emulator
	app.emulator.Init()
	if app.config.FileName != "" {
		app.emulator.LoadFile(app.config.FileName)
	}
	app.emulator.Start()

	// event loop
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
			case *sdl.JoyAxisEvent:
				app.processJoyAxis(e)
			case *sdl.JoyButtonEvent:
				app.processJoyButton(e)
			}
		}
		if count == 0 {
			time.Sleep(loopSleep)
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
		// Snaps
		case sdl.K_F2:
			app.emulator.TakeSnapshot()

		// Emulator
		case sdl.K_F5:
			app.emulator.Reset()
		case sdl.K_F6:
			if app.emulator.IsRunning() {
				app.emulator.Stop()
			} else {
				app.emulator.Start()
			}
		case sdl.K_F10:
			app.running = false // Exit app
		// Tape Drive
		case sdl.K_F7:
			if app.control.Tape().HasDrive() {
				if app.control.Tape().Drive().IsPlaying() {
					app.control.Tape().Drive().Stop()
				} else {
					app.control.Tape().Drive().Play()
				}
			}
		case sdl.K_F8:
			app.control.Tape().Drive().Rewind()
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
			app.control.Keyboard().KeyDown(int(e.Keysym.Scancode))
		} else {
			app.control.Keyboard().KeyUp(int(e.Keysym.Scancode))
		}
	}
}

func (app *App) processJoyAxis(e *sdl.JoyAxisEvent) {
	app.control.Joystick().AxisEvent(
		byte(e.Which), e.Axis, byte(e.Value>>8))
}

func (app *App) processJoyButton(e *sdl.JoyButtonEvent) {
	app.control.Joystick().ButtonEvent(
		byte(e.Which), e.Button, e.State)
}
