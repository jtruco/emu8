package js

import (
	"errors"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/device/io/keyboard"
)

// Emulator actions
const (
	actionStart      = "start"
	actionStop       = "stop"
	actionReset      = "reset"
	actionTapePlay   = "tape_play"
	actionTapeStop   = "tape_stop"
	actionTapeRewind = "tape_rewind"
)

var emuActions = []string{actionStart, actionStop, actionReset, actionTapePlay, actionTapeStop, actionTapeRewind}

// App is the Js/Wasm UI application
type App struct {
	jsUi     *JsUi
	jsLoop   JsFunc
	emulator *emulator.Emulator
	video    *Video
}

// NewApp returns a new user interface for canvas
func NewApp(cname string) *App {
	app := new(App)
	app.jsUi = NewJsUi(cname)
	app.jsUi.onKeyEvent = app.processKeyboard
	app.jsUi.onAction = app.processAction
	app.jsLoop = JsFuncOf(app.loop)
	app.video = NewVideo(&app.jsUi.canvas)
	app.configure()
	return app
}

// Connect connects the UI with the Emulator
func (app *App) configure() {}

// Connect connects the UI with the Emulator
func (app *App) Connect(emulator *emulator.Emulator) {
	app.emulator = emulator
	emulator.Control().Video().SetDisplay(app.video)
	app.video.screen = emulator.Control().Video().Device().Screen()
}

// Init initializes the UI
func (app *App) Init() error {
	// init UI
	if err := app.jsUi.Init(); err != nil {
		return err
	}
	if !app.video.Init() {
		return errors.New("UI: Could not initialize video subsystem")
	}
	return nil
}

// Run runs the UI loop
func (app *App) Run() {
	// init emulator
	app.emulator.Init()
	app.emulator.LoadFile(config.Get().App.File)

	// start emulation loop
	app.emulator.SetAsync(false)
	app.emulator.Start()
	RequestAnimationFrame(app.jsLoop)

	// FIXME : both approx seems to work (better ?)
	// app.emulator.SetAsync(true)
	// app.emulator.Start()
}

// loop
func (app *App) loop() {
	// emulate & sync
	if app.emulator.IsRunning() {
		app.emulator.Emulate()
		app.emulator.Sync()
	}
	RequestAnimationFrame(app.jsLoop) // Next
}

// Close closes all UI resources
func (app *App) Quit() {
	app.emulator.Stop()
	// FIXME : free resources (?)
}

// processKeyboard process keyboard input
func (app *App) processKeyboard(keyCode int, keyDown bool) {
	if keyDown {
		switch keyCode {
		case keyboard.KeyUnknown:
		case keyboard.KeyF5:
			app.emulator.Reset()
		case keyboard.KeyF6:
			if app.emulator.IsRunning() {
				app.emulator.Stop()
			} else {
				app.emulator.Start()
			}
			// Tape Drive
		case keyboard.KeyF7:
			app.emulator.Control().Tape().TogglePlay()
		case keyboard.KeyF8:
			app.emulator.Control().Tape().Rewind()
			// fullscreen
		case keyboard.KeyF11:
			app.jsUi.canvas.FullScreen()
		default:
			app.emulator.Control().Keyboard().KeyDown(keyCode)
		}
	} else if keyCode != keyboard.KeyUnknown {
		app.emulator.Control().Keyboard().KeyUp(keyCode)
	}
}

// processAction executes the given action
func (app *App) processAction(action string) {
	switch action {
	case actionStart:
		app.emulator.Start()
	case actionStop:
		app.emulator.Stop()
	case actionReset:
		app.emulator.Reset()
	case actionTapePlay:
		app.emulator.Control().Tape().TogglePlay()
	case actionTapeStop:
		app.emulator.Control().Tape().TogglePlay()
	case actionTapeRewind:
		app.emulator.Control().Tape().Rewind()
	}
}
