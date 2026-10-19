# emu8 - Makefile

# go commands
GO = go
GO_BUILD = $(GO) build -v
GO_CLEAN = $(GO) clean -v
GO_TEST = $(GO) test -v
GO_GET = $(GO) get -v

# app variables
APP = ./cmd/emu8
OUTPUT = emu8

# build & clean

.DEFAULT_GOAL = build

.PHONY: all
all: clean deps build cross

.PHONY: build
build: emu8

.PHONY: clean
clean:
	$(GO_CLEAN) $(APP)
	rm -f $(OUTPUT) $(OUTPUT)-*

.PHONY: emu8
emu8:
	$(GO_BUILD) -o $(OUTPUT) $(APP)

# cross compilation (windows)

GO_CGO_OPTS = CGO_ENABLED="1" CGO_LDFLAGS="-lmingw32 -lSDL2" CGO_CFLAGS="-D_REENTRANT"
GO_WIN_386 = GOOS="windows" GOARCH="386" GO386="sse2" CC="/usr/bin/i686-w64-mingw32-gcc" $(GO_CGO_OPTS)
GO_WIN_AMD64 = GOOS="windows" GOARCH="amd64" CC="/usr/bin/x86_64-w64-mingw32-gcc" $(GO_CGO_OPTS)

.PHONY: cross emu8-win-386 emu8-win-amd64
cross: emu8-win-386 emu8-win-amd64

emu8-win-386:
	$(GO_WIN_386) $(GO_BUILD) -o $(OUTPUT)-win-386.exe $(APP)

emu8-win-amd64:
	$(GO_WIN_AMD64) $(GO_BUILD) -o $(OUTPUT)-win-amd64.exe $(APP)

# dependencies

deps:
	$(GO_GET) ./...

# tests

.PHONY: test test-short test-sdl bench-sdl

# all tests (including the slow Z80 zexall test)
test:
	$(GO_TEST) ./...

# fast tests (slow tests are skipped with testing.Short)
test-short:
	$(GO_TEST) -short ./...

# SDL video output tests and benchmarks (build tag sdltest, see ui/sdl/README.md).
# -count=1 : SDL environment variables are not part of the go test cache key
SDL_TEST = $(GO) test -count=1 -tags sdltest

# correctness test, headless (dummy video driver, software renderer)
test-sdl:
	SDL_VIDEODRIVER=dummy $(SDL_TEST) -v ./ui/sdl/

# benchmarks on the display (renderer : SDL_RENDER_DRIVER=software|opengl...)
bench-sdl:
	$(SDL_TEST) -run '^$$' -bench . ./ui/sdl/
