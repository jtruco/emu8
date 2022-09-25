package sdl

import (
	"errors"
	"log"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
	"github.com/veandco/go-sdl2/sdl"
)

// app constants
const (
	loopSleepMillis = 10 // SDL poll interval
)

// App is the SDL UI App
type App struct {
	emulator *emulator.Emulator
	video    *Video
	audio    *Audio
	filename string
	async    bool
	running  bool
}

// NewApp creates a new application
func NewApp() *App {
	app := new(App)
	app.video = NewVideo()
	app.audio = NewAudio()
	app.configure(config.Get())
	return app
}

func (app *App) configure(conf *config.Config) {
	app.filename = conf.App.File
	app.async = conf.Emulator.Async
	app.video.Title = conf.App.Title
	app.video.FullScreen = conf.Video.FullScreen
	app.video.Scale = float32(conf.Video.Scale)
	app.video.Async = conf.Emulator.Async
	app.audio.Frequency = int32(conf.Audio.Frequency)
	app.audio.Mute = conf.Audio.Mute
	app.audio.Async = conf.Emulator.Async
}

// Connect connects the Emulator
func (app *App) Connect(emulator *emulator.Emulator) {
	app.emulator = emulator
	emulator.Control().Video().SetDisplay(app.video)
	app.video.Screen = emulator.Control().Video().Device().Screen()
	emulator.Control().Audio().SetPlayer(app.audio)
}

// Init the SDL App
func (app *App) Init() error {
	// init sdl
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_JOYSTICK); err != nil {
		log.Println("SDL : Error initializing SDL:", err.Error())
		return errors.New("could not initialize SDL library")
	}
	// init SDL video output
	if !app.video.Init() {
		app.Quit()
		return errors.New("could not initialize video subsystem")
	}
	// init SDL audio
	if !app.audio.Init() {
		app.Quit()
		return errors.New("could not initialize audio subsystem")
	}
	// open sdl joystick (only one supported)
	if sdl.NumJoysticks() > 0 {
		joystick := sdl.JoystickOpen(0)
		if joystick != nil {
			log.Println("SDL : Joystick found:", sdl.JoystickNameForIndex(0))
		}
	}
	log.Print("SDL : Initialized")
	return nil
}

// Run the SDL App
func (app *App) Run() {
	// init emulator
	app.emulator.SetAsync(app.async)
	app.emulator.Init()
	if app.filename != "" {
		app.emulator.LoadFile(app.filename)
	}
	app.emulator.Start()
	// event loop
	app.running = true
	for app.running {
		app.pollEvents() // Poll SDL events
		if !app.async {  // Sync emulation
			app.emulator.Emulate()
			app.emulator.Sync()
			continue
		}
		select { // Async emulation
		case <-app.video.UpdateUi:
			app.video.OnUpdate(false)
		default:
			sdl.Delay(loopSleepMillis)
		}
	}
	// end emulation
	app.emulator.Stop()
}

// Quit closes the SDL resources
func (app *App) Quit() {
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
		case *sdl.WindowEvent:
			app.processWindowEvent(e)
		case *sdl.KeyboardEvent:
			app.processKeyboard(e)
		case *sdl.JoyAxisEvent:
			app.processJoyAxis(e)
		case *sdl.JoyButtonEvent:
			app.processJoyButton(e)
		}
	}
}

func (app *App) processWindowEvent(e *sdl.WindowEvent) {
	switch e.Event {
	case sdl.WINDOWEVENT_SHOWN, sdl.WINDOWEVENT_RESIZED:
		app.video.OnUpdate(true) // Refresh screen
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
			app.emulator.Control().Tape().TogglePlay()
		case sdl.K_F8:
			app.emulator.Control().Tape().Rewind()
		// UI
		case sdl.K_F4:
			app.audio.Mute = !app.audio.Mute
			if app.audio.Mute {
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
			app.emulator.Control().Keyboard().KeyDown(int(e.Keysym.Scancode))
		} else {
			app.emulator.Control().Keyboard().KeyUp(int(e.Keysym.Scancode))
		}
	}
}

func (app *App) processJoyAxis(e *sdl.JoyAxisEvent) {
	app.emulator.Control().Joystick().AxisEvent(
		byte(e.Which), e.Axis, byte(e.Value>>8))
}

func (app *App) processJoyButton(e *sdl.JoyButtonEvent) {
	app.emulator.Control().Joystick().ButtonEvent(
		byte(e.Which), e.Button, e.State)
}
