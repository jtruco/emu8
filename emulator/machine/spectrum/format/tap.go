package format

import (
	"log"

	"github.com/jtruco/emu8/emulator/device/io/tape"
)

// -----------------------------------------------------------------------------
// TAP tape format
// -----------------------------------------------------------------------------

// TAP format extension
const TAP = "tap"

// TAP blocks types
const (
	tapBlockHeader = 0x00
	tapBLockData   = 0xff
)

// TapBlockHeader information
type TapBlockHeader struct {
	tapType  byte
	filename string
	length   uint16
	par1     uint16
	par2     uint16
}

// TapBlock is a tape block
type TapBlock struct {
	tape.BlockInfo
	data   []byte
	header TapBlockHeader
}

// Info gets block information
func (block *TapBlock) Info() *tape.BlockInfo {
	return &block.BlockInfo
}

// Data gets block data bytes
func (block *TapBlock) Data() []byte {
	return block.data
}

// Tap implements the a tape format .TAP
type Tap struct {
	info        tape.Info    // Tape information
	blocks      []tape.Block // Block array
	pilotPulses int          // Pilot pulses
	bitMask     byte         // Current bit mask
	bitTime     int          // Current bit time
}

// NewTap creates a new tape
func NewTap() tape.Tape {
	tap := new(Tap)
	tap.blocks = make([]tape.Block, 0, 2)
	return tap
}

// Info gets tape information
func (tap *Tap) Info() *tape.Info {
	return &tap.info
}

// Blocks gets the tape blocks
func (tap *Tap) Blocks() []tape.Block {
	return tap.blocks
}

// Load loads the tape file data
func (tap *Tap) Load(data []byte) bool {
	tapeLength := len(data)
	if tapeLength == 0 {
		log.Print("Tape (TAP) : Invalid format: 0-length data")
		return false
	}
	index := 0
	for offset := 0; offset < tapeLength; {
		block := new(TapBlock)
		length := int(readWord(data, offset))
		offset += 2
		block.Type = data[offset]
		block.Index = index
		block.Offset = offset
		block.Length = length
		block.data = data[offset : offset+length]
		if block.Type == tapBlockHeader {
			block.header.tapType = block.data[1]
			block.header.filename = readString(block.data, 2, 10)
			block.header.length = readWord(block.data, 12)
			block.header.par1 = readWord(block.data, 14)
			block.header.par2 = readWord(block.data, 16)
		}
		tap.blocks = append(tap.blocks, block)
		offset += length
		index++
	}
	return true
}

// Play tap
func (tap *Tap) Play(control *tape.Control) {
	switch control.State {

	case tapeStateStart:
		control.Block = tap.blocks[control.BlockIndex]
		control.BlockPos = 0
		if control.Block.Info().Type == tapBlockHeader {
			tap.pilotPulses = tapeHeaderPulses
		} else {
			tap.pilotPulses = tapeDataPulses
		}
		control.Ear = tape.LevelLow
		control.Timeout = 0
		control.State = tapeStatePilot
		// log tape block
		tapblock := control.Block.(*TapBlock)
		if tapblock.Type == tapBlockHeader {
			log.Println("Tape (TAP) : Header block:", tapblock.header.filename)
		} else {
			log.Println("Tape (TAP) : Data block:", tapblock.Length, "bytes")
		}

	case tapeStatePilot:
		control.Ear ^= tape.LevelMask
		tap.pilotPulses--
		if tap.pilotPulses > 0 {
			control.Timeout = tapeTimingPilot
		} else {
			control.Timeout = tapeTimingSync1
			control.State = tapeStateSync
		}

	case tapeStateSync:
		control.Ear ^= tape.LevelMask
		control.Timeout = tapeTimingSync2
		control.State = tapeStateByte

	case tapeStateByte:
		tap.bitMask = 0x80
		control.State = tapeStateBit1

	case tapeStateBit1:
		control.Ear ^= tape.LevelMask
		if (control.DataAtPos() & tap.bitMask) == 0 {
			tap.bitTime = tapeTimingZero
		} else {
			tap.bitTime = tapeTimingOne
		}
		control.Timeout = tap.bitTime
		control.State = tapeStateBit2

	case tapeStateBit2:
		control.Ear ^= tape.LevelMask
		control.Timeout = tap.bitTime
		tap.bitMask >>= 1
		if tap.bitMask == 0 {
			control.BlockPos++
			if !control.EndOfBlock() {
				control.State = tapeStateByte
			} else {
				control.State = tapeStatePause
			}
		} else {
			control.State = tapeStateBit1
		}

	case tapeStatePause:
		control.Ear ^= tape.LevelMask
		control.Timeout = tapeTimingEoB
		control.State = tapeStatePauseStop

	case tapeStatePauseStop:
		control.BlockIndex++
		if control.EndOfTape() {
			control.State = tapeStateStop
		} else {
			control.State = tapeStateStart
		}

	case tapeStateStop:
		control.Playing = false // Stop

	default:
		control.State = tapeStateStop
	}
}
