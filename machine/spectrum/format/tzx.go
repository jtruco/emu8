package format

import (
	"log"

	"github.com/jtruco/emu8/device/tape"
)

// -----------------------------------------------------------------------------
// TZX tape format
// -----------------------------------------------------------------------------

// TZX constants
const (
	tzxHeaderSignature = "ZXTape!"
	tzxStartEar        = TapeEarOff
)

// TZX states
const (
	_ = iota + tapeStateStop
	tapeStateTzxHeader
	tapeStateLeaderNochg
	tapeStateNewByteNochg
	tapeStateLastPulse
	tapeStatePureTone
	tapeStatePureToneNochg
	tapeStatePulseSequence
	tapeStatePulseSequenceNochg
)

// TzxBlock is a tape block
type TzxBlock struct {
	tape.BlockInfo
	data []byte
}

// Info gets block information
func (block *TzxBlock) Info() *tape.BlockInfo {
	return &block.BlockInfo
}

// Data gets block data bytes
func (block *TzxBlock) Data() []byte {
	return block.data
}

// Tzx implements the a tape format .TAP
type Tzx struct {
	info          tape.Info    // Tape information
	data          []byte       // Data buffer
	blocks        []tape.Block // block array
	blockLen      int
	leaderPulses  int
	leaderLenght  int
	sync1Lenght   int
	sync2Lenght   int
	zeroLenght    int
	oneLenght     int
	bitsLastByte  byte
	endBlockPause int
	mask          byte
	bitTime       int
	loopCount     int
	loopStart     int
}

// NewTzx creates a new tape
func NewTzx() tape.Tape {
	tzx := &Tzx{}
	tzx.blocks = make([]tape.Block, 0, 2)
	return tzx
}

// Info gets tape information
func (tzx *Tzx) Info() *tape.Info {
	return &tzx.info
}

// Blocks gets the tape blocks
func (tzx *Tzx) Blocks() []tape.Block {
	return tzx.blocks
}

// Load loads the tape file data
func (tzx *Tzx) Load(data []byte) bool {
	tapeLength := len(data)
	if tapeLength == 0 {
		log.Print("TZX : Invalid format. 0-length data.")
		return false
	}
	if string(data[0:7]) != tzxHeaderSignature {
		log.Print("TZX : Invalid TZX header signature.")
		return false
	}
	index := 0
	for offset := 0; offset < tapeLength; {
		block := &TapBlock{}
		block.Type = data[offset]
		block.Index = index
		block.Offset = offset
		length := 0
		switch block.Type {
		case 0x10: // Standard speed data
			length = 5 + readInt(data, offset+3)
		case 0x11: // Turbo speed data
			length = 19 + readIntN(data, offset+16, 3)
		case 0x12: // Pure tone
			length = 5
		case 0x13: // Pulse sequence
			length = 2 + 2*int(data[offset+1])
		case 0x14: // Pure data
			length = 11 + readIntN(data, offset+8, 3)
		case 0x15: // Direct data
			length = 9 + readIntN(data, offset+6, 3)
		case 0x18: // CSW recording
			length = 5 + readIntN(data, offset+1, 4)
		case 0x19: // Generalized data
			length = 5 + readIntN(data, offset+1, 4)
		case 0x20: // Pause (silence) or 'Stop the Tape' command
			length = 3
		case 0x21: // Group Start
			length = 2 + int(data[offset+1])
		case 0x22: // Group End
			length = 1
		case 0x23: // Jump to Block
			length = 3
		case 0x24: // Loop Start
			length = 3
		case 0x25: // Loop End
			length = 1
		case 0x26: // Call Sequence
			length = 3 + 2*readInt(data, offset+1)
		case 0x27: // Return from Sequence
			length = 1
		case 0x28: // Select Block
			length = 3 + readInt(data, offset+1)
		case 0x2A: // Stop the tape if in 48K mode
			length = 5
		case 0x2B: // Set Signal Level
			length = 6
		case 0x30: // Text Description
			length = 2 + int(data[offset+1])
		case 0x31: // Message
			length = 3 + int(data[offset+2])
		case 0x32: // Archive Info
			length = 3 + readInt(data, offset+1)
		case 0x33: // Hardware Type
			length = 2 + 3*int(data[offset+1])
		case 0x35: // Custom Info
			length = 21 + 2*readIntN(data, offset+17, 4)
		case 'Z': // Glue Block
			length = 10
		default:
			log.Printf("TZX : Unknown ID block %x \n", block.Type)
			return false
		}
		block.Length = length
		block.data = data[offset : offset+length]
		tzx.blocks = append(tzx.blocks, block)
		offset += length
		index++
	}
	return true
}

// Play TZX tape
func (tzx *Tzx) Play(control *tape.Control) {
	switch control.State {

	case tapeStateStart:
		control.Ear = tzxStartEar
		control.State = tapeStateTzxHeader

	case tapeStateTzxHeader:
		if control.EndOfTape() {
			control.State = tapeStateStop
		} else {
			control.Block = tzx.blocks[control.BlockIndex]
			control.BlockPos = 0
			tzx.parseHeader(control)
			// log.Printf("TZX : Playing block #%d . Type : %x", control.BlockIndex, control.Block.Info().Type)
		}

	case tapeStateLeader:
		control.Ear ^= TapeEarMask
		control.State = tapeStateLeaderNochg

	case tapeStateLeaderNochg:
		tzx.leaderPulses--
		if tzx.leaderPulses > 0 {
			control.State = tapeStateLeader
			control.Timeout = tzx.leaderLenght
		} else {
			control.State = tapeStateSync
			control.Timeout = tzx.sync1Lenght
		}

	case tapeStateSync:
		control.Ear ^= TapeEarMask
		control.State = tapeStateNewByte
		control.Timeout = tzx.sync2Lenght

	case tapeStateNewByteNochg:
		control.Ear ^= TapeEarMask
		control.State = tapeStateNewByte

	case tapeStateNewByte:
		tzx.mask = 0x80
		control.State = tapeStateNewBit

	case tapeStateNewBit:
		control.Ear ^= TapeEarMask
		if (control.DataAtPos() & tzx.mask) == 0 {
			tzx.bitTime = tzx.zeroLenght
		} else {
			tzx.bitTime = tzx.oneLenght
		}
		control.State = tapeStateHalf2
		control.Timeout = tzx.bitTime

	case tapeStateHalf2:
		control.Ear ^= TapeEarMask
		control.Timeout = tzx.bitTime
		tzx.mask >>= 1
		lastBit := byte(0)
		if tzx.blockLen == 1 {
			lastBit = 0x80 >> tzx.bitsLastByte
		}
		if tzx.mask == lastBit {
			control.BlockPos++
			tzx.blockLen--
			if tzx.blockLen > 0 {
				control.State = tapeStateNewByte
			} else {
				control.State = tapeStateLastPulse
			}
		} else {
			control.State = tapeStateNewBit
		}

	case tapeStateLastPulse:
		control.Ear ^= TapeEarMask
		control.State = tapeStatePause
		control.Timeout = 3500 // TZX 1 ms

	case tapeStatePause:
		control.Ear = tzxStartEar
		control.State = tapeStateTzxHeader
		if !control.EndOfTape() {
			control.Timeout = tzx.endBlockPause * tapeTstatesMs
		}

	case tapeStatePureTone:
		control.Ear ^= TapeEarMask
		control.State = tapeStatePureToneNochg

	case tapeStatePureToneNochg:
		if tzx.leaderPulses > 0 {
			tzx.leaderPulses--
			control.Timeout = tzx.leaderLenght
			control.State = tapeStatePureTone
		} else {
			control.State = tapeStateTzxHeader
		}

	case tapeStatePulseSequence:
		control.Ear ^= TapeEarMask
		control.State = tapeStatePulseSequenceNochg

	case tapeStatePulseSequenceNochg:
		if tzx.leaderPulses > 0 {
			tzx.leaderPulses--
			control.Timeout = readInt(control.Block.Data(), control.BlockPos)
			control.BlockPos += 2
			control.State = tapeStatePulseSequence
		} else {
			control.State = tapeStateTzxHeader
		}

	case tapeStateStop:
		control.Playing = false // Stop

	default:
		control.State = tapeStateStop
	}
}

func (tzx *Tzx) parseHeader(control *tape.Control) {
	data := control.Block.Data()
	id := control.Block.Info().Type
	switch id {

	case 0x10:
		tzx.leaderLenght = tapeLeaderLenght
		tzx.sync1Lenght = tapeSync1Lenght
		tzx.sync2Lenght = tapeSync2Lenght
		tzx.zeroLenght = tapeZeroLenght
		tzx.oneLenght = tapeOneLenght
		tzx.bitsLastByte = 8
		tzx.endBlockPause = readInt(data, control.BlockPos+1)
		tzx.blockLen = readInt(data, control.BlockPos+3)
		control.BlockPos += 5
		if control.DataAtPos() < 0x80 {
			tzx.leaderPulses = tapeHeaderPulses
		} else {
			tzx.leaderPulses = tapeDataPulses
		}
		control.State = tapeStateLeaderNochg
		control.BlockIndex++

	case 0x11:
		tzx.leaderLenght = readInt(data, control.BlockPos+1)
		tzx.sync1Lenght = readInt(data, control.BlockPos+3)
		tzx.sync2Lenght = readInt(data, control.BlockPos+5)
		tzx.zeroLenght = readInt(data, control.BlockPos+7)
		tzx.oneLenght = readInt(data, control.BlockPos+9)
		tzx.leaderPulses = readInt(data, control.BlockPos+11)
		tzx.bitsLastByte = data[control.BlockPos+13]
		tzx.endBlockPause = readInt(data, control.BlockPos+14)
		tzx.blockLen = readIntN(data, control.BlockPos+16, 3)
		control.BlockPos += 19
		control.State = tapeStateLeaderNochg
		control.BlockIndex++

	case 0x12: // Pure Tone Block
		tzx.leaderLenght = readInt(data, control.BlockPos+1)
		tzx.leaderPulses = readInt(data, control.BlockPos+3)
		control.BlockPos += 5
		control.State = tapeStatePureToneNochg
		control.BlockIndex++

	case 0x13: // Pulse Sequence Block
		tzx.leaderPulses = int(data[control.BlockPos+1])
		control.BlockPos += 2
		control.State = tapeStatePulseSequenceNochg
		control.BlockIndex++

	case 0x14: // Pure Data Block
		tzx.zeroLenght = readInt(data, control.BlockPos+1)
		tzx.oneLenght = readInt(data, control.BlockPos+3)
		tzx.bitsLastByte = data[control.BlockPos+5]
		tzx.endBlockPause = readInt(data, control.BlockPos+6)
		tzx.blockLen = readIntN(data, control.BlockPos+8, 3)
		control.BlockPos += 11
		control.State = tapeStateNewByteNochg
		control.BlockIndex++

	case 0x20: // Pause (silence) or 'Stop the Tape' command
		tzx.endBlockPause = readInt(data, control.BlockPos+1)
		control.BlockPos += 3
		control.State = tapeStatePauseStop
		control.BlockIndex++

	case 0x21: // Group Start
		control.BlockIndex++

	case 0x22: // Group End
		control.BlockIndex++

	case 0x23: // Jump to Block
		target := readInt(data, control.BlockPos+1)
		control.BlockIndex += target

	case 0x24: // Loop Start
		tzx.loopCount = readInt(data, control.BlockPos+1)
		control.BlockIndex++
		tzx.loopStart = control.BlockIndex

	case 0x25: // Loop End
		tzx.loopCount--
		if tzx.loopCount == 0 {
			control.BlockIndex++
		} else {
			control.BlockIndex = tzx.loopStart
		}

	case 0x28: // Select Block
		control.BlockIndex++

	case 0x2A: // Stop the tape if in 48K mode
		control.State = tapeStateStop
		control.BlockIndex++

	case 0x2B: // Set Signal Level
		if data[control.BlockPos+5] == 0 {
			control.Ear = TapeEarOff
		} else {
			control.Ear = TapeEarOn
		}
		control.BlockIndex++

	case 0x30: // Text Description
		control.BlockIndex++

	case 0x31: // Message Block
		control.BlockIndex++

	case 0x32: // Archive Info
		control.BlockIndex++

	case 0x33: // Hardware Type
		control.BlockIndex++

	case 0x35: // Custom Info Block
		control.BlockIndex++

	case 'Z': // TZX Header && "Glue" Block
		control.BlockIndex++

	default:
		log.Printf("TZX : Playing block #%d . Unsupported type : %x", control.BlockIndex, id)
		control.BlockIndex++
	}
}
