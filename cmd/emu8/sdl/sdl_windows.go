package sdl

import "os"

// windows specific sdl initialization
func init() {
	os.Setenv("SDL_AUDIODRIVER", "directsound")
}
