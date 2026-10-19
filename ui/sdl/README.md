# SDL user interface : tests

The SDL video output tests and benchmarks (`video_test.go`) open SDL windows, so they are behind the build tag
`sdltest` : `go test ./...`, `make test` and the CI don't compile or run them (the CI only vets them, to keep them
compiling).

## Commands

From the repository root :
```
make test-sdl                               # correctness test, headless (SDL_VIDEODRIVER=dummy, software renderer)
go test -count=1 -tags sdltest -v ./ui/sdl/ # correctness test on the display (accelerated renderer)
make bench-sdl                              # benchmarks, on the display
SDL_RENDER_DRIVER=software make bench-sdl   # benchmarks with another renderer (software, opengl, opengles2...)
```

## Tests

- `TestVideoOutputMatchesScreen` draws partial, full and forced (refresh) screen updates into an offscreen target of
  the window size and reads it back : it must be equal to the screen pixels. It checks the dirty band upload of the
  streaming texture and the copy. Window pixels can't be used : they may be undefined when read back (double
  buffering, OpenGL with a window not mapped or covered).
- `BenchmarkVideoOutput/{zx,cpc}/{full,partial}` measures one screen update (texture upload, copy and present) with
  the Spectrum (320x240) and CPC (768x264, x0.5) screen geometries at scale 2, for full screen changes and for a
  moving 32x32 sprite (partial changes, dirty regions). The renderer has no vsync; if the driver forces it, the times
  are the display refresh period.

The tests use synthetic screens : they need no ROMs or game files. They skip when SDL can't be initialized (no
display).

## Writing SDL tests

- Put them behind the build tag : `//go:build sdltest`.
- Lock the test goroutine to its OS thread while SDL objects exist (`runtime.LockOSThread`, see `newTestVideo`) :
  OpenGL contexts and X11 must be used from one thread, and Go moves goroutines between threads (random black or
  garbage readbacks otherwise). The app is not affected : go-sdl2 locks the main goroutine to the main thread.
- Run them with `-count=1`, as the make targets do : SDL reads its environment variables from C, so the Go test cache
  doesn't see them.
