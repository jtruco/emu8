package sdl

import (
	"log"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/veandco/go-sdl2/sdl"
)

// App is the SDL application
type App struct {
	config   *config.Config
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
	app.video = NewVideo(app.config)
	app.audio = NewAudio(app.config)
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
	control := emu.Controller()
	control.Video().SetDisplay(app.video)
	control.Audio().SetPlayer(app.audio)
	app.control = control
	// init SDL video output
	if !app.video.Init(control.Video().Device()) {
		app.End()
		return false
	}
	// init SDL audio
	if !app.audio.Init() {
		app.End()
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
		app.pollEvents()
		app.emulator.Emulate()
		app.emulator.Sync()
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
			if app.control.Tape().HasDrive() {
				if app.control.Tape().Drive().HasTape() {
					if app.control.Tape().Drive().IsPlaying() {
						app.control.Tape().Drive().Stop()
					} else {
						app.control.Tape().Drive().Play()
					}
				} else {
					log.Println("App : There is no tape loaded !")
				}
			} else {
				log.Println("App : Machine has no tape drive !")
			}
		case sdl.K_F8:
			if app.control.Tape().HasDrive() {
				if app.control.Tape().Drive().HasTape() {
					app.control.Tape().Drive().Rewind()
				} else {
					log.Println("App : There is no tape loaded !")
				}
			} else {
				log.Println("App : Machine has no tape drive !")
			}
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
