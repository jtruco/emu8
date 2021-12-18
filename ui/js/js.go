// Package js is the Js/Wasm user interface
package js

const (
	emu8_default_ui = "emu8"
)

var myJsApp = NewApp(emu8_default_ui)

// Get returns a JavaScript user application
func Get() *App { return myJsApp }
