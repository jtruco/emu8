package format

import (
	"log"

	"github.com/jtruco/emu8/device/tape"
)

// -----------------------------------------------------------------------------
// TAP tape format
// -----------------------------------------------------------------------------

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
	info         tape.Info    // Tape information
	data         []byte       // Data buffer
	blocks       []tape.Block // block array
	leaderPulses int
	mask         byte
	bitTime      int
}

// NewTap creates a new tape
func NewTap() tape.Tape {
	tap := &Tap{}
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
		log.Print("TAP : Invalid format. 0-length data.")
		return false
	}
	index := 0
	for offset := 0; offset < tapeLength; {
		block := &TapBlock{}
		length := int(readWord(data, offset))
		offset += 2
		block.Type = data[offset]
		block.Index = index
		block.Offset = offset
		block.Length = length
		block.data = data[offset : offset+length]
		if block.Type == tapBlockHeader {
			block.header.tapType = data[1]
			block.header.filename = string(data[2:12])
			block.header.length = readWord(data, 12)
			block.header.par1 = readWord(data, 14)
			block.header.par2 = readWord(data, 16)
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
			tap.leaderPulses = tapeHeaderPulses
		} else {
			tap.leaderPulses = tapeDataPulses
		}
		control.Ear = TapeEarOff
		control.State = tapeStateLeader
		control.Timeout = tapeLeaderLenght
		// log.Println("TAP : Load block ", control.Block.(*TapBlock).header.filename)

	case tapeStateLeader:
		control.Ear ^= TapeEarMask
		tap.leaderPulses--
		if tap.leaderPulses > 0 {
			control.Timeout = tapeLeaderLenght
		} else {
			control.State = tapeStateSync
			control.Timeout = tapeSync1Lenght
		}

	case tapeStateSync:
		control.Ear ^= TapeEarMask
		control.State = tapeStateNewByte
		control.Timeout = tapeSync2Lenght

	case tapeStateNewByte:
		tap.mask = 0x80
		control.State = tapeStateNewBit

	case tapeStateNewBit:
		control.Ear ^= TapeEarMask
		if (control.DataAtPos() & tap.mask) == 0 {
			tap.bitTime = tapeZeroLenght
		} else {
			tap.bitTime = tapeOneLenght
		}
		control.State = tapeStateHalf2
		control.Timeout = tap.bitTime

	case tapeStateHalf2:
		control.Ear ^= TapeEarMask
		control.Timeout = tap.bitTime
		tap.mask >>= 1
		if tap.mask == 0 {
			control.BlockPos++
			if !control.EndOfBlock() {
				control.State = tapeStateNewByte
			} else {
				control.State = tapeStatePause
			}
		} else {
			control.State = tapeStateNewBit
		}

	case tapeStatePause:
		control.Ear ^= TapeEarMask
		control.State = tapeStatePauseStop
		control.Timeout = 3500

	case tapeStatePauseStop:
		control.BlockIndex++
		if control.EndOfTape() {
			control.State = tapeStateStop
		} else {
			control.State = tapeStateStart // Next block
		}

	case tapeStateStop:
		control.Playing = false // Stop

	default:
		control.State = tapeStateStop
	}
}
