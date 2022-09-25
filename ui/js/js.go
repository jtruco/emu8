// Package js is the Js/Wasm user interface
package js

const emu8_default_ui = "emu8"

// Get returns a JavaScript user application
func GetApp() *App { return NewApp(emu8_default_ui) }
