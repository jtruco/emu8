// Package sdl is the libSDLv2 user interface implementation
package sdl

import (
	"log"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/veandco/go-sdl2/sdl"
)

// app constants
const (
	loopSleepMillis = 10 // SDL poll interval
)

// App is the SDL application
type App struct {
	config   *config.Config
	async    bool
	video    *Video
	audio    *Audio
	emulator *emulator.Emulator
	control  *controller.Controller
	running  bool
}

// NewApp creates a new application
func NewApp() *App {
	app := new(App)
	app.config = config.Get()
	app.async = app.config.EmulatorAsync
	app.video = NewVideo(app)
	app.audio = NewAudio(app)
	return app
}

// Init the SDL App
func (app *App) Init(emu *emulator.Emulator) bool {
	// init sdl
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_JOYSTICK); err != nil {
		log.Println("SDL : Error initializing SDL:", err.Error())
		return false
	}
	log.Print("SDL : Initialized")
	// open sdl joystick (only one supported)
	if sdl.NumJoysticks() > 0 {
		joystick := sdl.JoystickOpen(0)
		if joystick != nil {
			log.Println("SDL : Joystick found:", sdl.JoystickNameForIndex(0))
		}
	}
	// init emulator
	app.emulator = emu
	app.control = emu.Control()
	app.control.Video().SetDisplay(app.video)
	app.control.Audio().SetPlayer(app.audio)
	// init SDL video output
	if !app.video.Init(app.control.Video().Device()) {
		app.End()
		return false
	}
	// init SDL audio
	if !app.audio.Init(app.control.Audio().Device()) {
		app.End()
		return false
	}
	return true
}

// Run the SDL App
func (app *App) Run() {

	// init emulator
	app.emulator.Init()
	app.emulator.SetAsync(app.async)
	if app.config.FileName != "" {
		app.emulator.LoadFile(app.config.FileName)
	}
	app.emulator.Start()

	// event loop
	app.running = true
	for app.running {
		// Poll SDL events
		app.pollEvents()

		// Sync emulation
		if !app.async {
			app.emulator.Emulate()
			app.emulator.Sync()
			continue
		}

		// Async emulation
		select {
		case <-app.video.updateUi:
			app.video.onUpdate(false)
		default:
			sdl.Delay(loopSleepMillis)
		}
	}

	app.emulator.Stop()
}

// End the SDL App
func (app *App) End() {
	if app.audio != nil {
		app.audio.Close()
	}
	if app.video != nil {
		app.video.Destroy()
	}
	sdl.Quit()
	log.Print("App : Terminated !")
}

// poll SDL event queue
func (app *App) pollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
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
			log.Print("App : Exiting app")
		// Tape Drive
		case sdl.K_F7:
			app.control.Tape().TogglePlay()
		case sdl.K_F8:
			app.control.Tape().Rewind()
		// UI
		case sdl.K_F4:
			app.audio.config.Mute = !app.audio.config.Mute
			if app.audio.config.Mute {
				log.Println("App : Audio is muted")
			} else {
				log.Println("App : Audio is enabled")
			}
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
