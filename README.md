# emu8

[![License](https://img.shields.io/github/license/jtruco/emu8.svg?style=flat)](https://github.com/jtruco/emu8/blob/master/LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/jtruco/emu8)](https://github.com/jtruco/emu8/releases/latest)
[![Build Status](https://travis-ci.com/jtruco/emu8.svg?branch=master)](https://travis-ci.com/jtruco/emu8)
[![Go Report Card](https://goreportcard.com/badge/github.com/jtruco/emu8)](https://goreportcard.com/report/github.com/jtruco/emu8)

**The 8-bit machine emulator**

---

**emu8** is an open source 8-bit machine emulator written in the Go programming language.

Currently these machine models are supported :
- Sinclair ZX Spectrum 16K and 48K
- Amstrad CPC 464

There are plans to implement more 8-bit machines and models like : ZX80, ZX81, Commodore 64, BBC Micro A/B, MSX1 ...

## Installation

To download and build the command line application you can use :
```
go get github.com/jtruco/emu8
```

See the [requirements](#requirements) section for additional build and dependency instructions.

## Requirements

**emu8** is mostly programmed using the Go standard library. The core emulator engine can compile in almost all platforms supported by the Go compiler.

The current frontend requires [SDL2](http://libsdl.org/) (Simple Directmedia Layer library, version 2), and its Go binding package *go-sdl2*.
In order to build this dependency is necessary to install the essential C/C++ building tools and the SDL2 development library.

On Ubuntu / Debian like distributions you can install :
```
sudo apt-get install build-essential
sudo apt-get install libsdl2-dev
```

Then you can download and build the binding :
```
go get github.com/veandco/go-sdl2
```

For other platforms like *Windows* or *macOS*, please refer to the [go-sdl2](https://github.com/veandco/go-sdl2) compiling instructions.

## Usage

**emu8** is a command line application that uses SDL2 as its user interface.

To run the emulator and load a snapshot file into the machine just type :
```
./emu8 manicminer.sna
```

**emu8** looks for files in the current working directory. If it fails, then it tries in the following subdirectories by type :
- ./rom : ROM files (*.rom)
- ./snap : Snapshot image files (.sna, .z80, ...)
- ./tape : Tape container files (.tap, .tzx, cdt, ...)

The default machine model is the classic *Speccy* or *ZX Spectrum 48k*.
To select another machine model use :
```
./emu8 -m cpc464 blagger.sna
```

*Current supported models are : zx16k, zx48k, and cpc464.*

### Keyboard accelerators
Once the emulator is running you can control it with the following keys :
- Esc : Exits the application.
- F2 : Takes a snapshot of the machine state and saves it.
- F5 : Resets the machine to its initial state.
- F6 : Pauses and Resumes the machine emulation.
- F7 : Plays and Stops the tape.
- F8 : Rewinds the tape.
- F10 : Exits the application.
- F11 : Toggle full-screen video mode.

## Features

General status and main features :
- Written in pure Go.
- Multi-platform desktop support (Linux, Windows, macOS).
- Multi-machine architecture.
- User interface : video, audio and user input.
- Joystick support (only one port by now).
- Video scale2x and fullscreen (beta) support.
- Zip compressed files support.

### Sinclair ZX Spectrum ( Status : Release )
The emulation is stable and accurate for the current supported models :
- ZX Spectrum 16k and 48k models supported.
- Zilog Z80 CPU emulation.
- Contended video memory emulation.
- Accurate border and scanline video effects.
- Beeper emulation.
- Snapshot formats supported : SNA, Z80.
- Tape formats supported (read only) : TAP, TZX.
- Kempston joystick support.

### Amstrad CPC ( Status : Beta )
The emulation is good, but needs some fixes and audio quality improvements :
- Amstrad CPC 464 model supported.
- Zilog Z80 CPU emulation.
- MC6845 CRTC device emulation.
- Accurate scanline and video timings emulation.
- AY-3-8912 audio device emulation (alpha).
- Snapshot formats supported : SNA.
- Tape formats supported (read only) : CDT.
- Joystick support.

## Roadmap
These are the main goals and features for the next versions :
- Web application using [WebAssembly](https://github.com/golang/go/wiki/WebAssembly) and HTML5/Js.
- Support more machines and models.
- [libretro](https://github.com/libretro) core implementation.
- Mobile : native apps for Android & iOS.
- Gui : A cross-platform desktop user interface.
- OpenGL / OpenAL desktop app frontend.

## Contributing

**emu8** is currently in alpha development stage, and its structure is subject to change frequently.

This is a spare time project. I hope to have a more stable code base in next releases.

## License

**emu8** is licensed under GNU General Public License v3.0. See LICENSE file for more details.
