package format

import (
	"log"

	"github.com/jtruco/emu8/emulator/device/io/tape"
)

// -----------------------------------------------------------------------------
// TZX tape format
// -----------------------------------------------------------------------------

// TZX format extension
const TZX = "tzx"

// TZX constants
const (
	tzxHeaderSignature = "ZXTape!"
	tzxHeaderLength    = 10 // signature, 0x1A, major and minor version
	tzxStartEar        = tape.LevelLow
	tzxLogAllBlocks    = false
)

// TZX states
const (
	_ = iota + tapeStateStop
	tapeStateTzxHeader
	tapeStatePilotNc
	tapeStateByteNc
	tapeStateLastPulse
	tapeStatePureTone
	tapeStatePureToneNc
	tapeStatePulseSeq
	tapeStatePulseSeqNc
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

// Tzx implements the a tape format .TZX
type Tzx struct {
	info          tape.Info    // Tape information
	blocks        []tape.Block // block array
	blockLength   int          // Block length
	pilotPulses   int          // Pilot pulses
	pilotTiming   int          // Pilot timing
	sync1Timing   int          // Sync1 timing
	sync2timing   int          // Sync2 timing
	zeroTiming    int          // Timing of 0 bit
	oneTiming     int          // Timing of 1 bit
	bitsLastByte  byte         // Number of bits of last byte
	endBlockPause int          // Pause at end of block
	bitMask       byte         // Current bit mask
	bitTime       int          // Current bit time
	lastBit       byte         // Last bit of current byte
	loopCount     int          // Control loop count
	loopStart     int          // Control loop start
}

// NewTzx creates a new tape
func NewTzx() tape.Tape {
	tzx := new(Tzx)
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
	if tapeLength < tzxHeaderLength || string(data[0:7]) != tzxHeaderSignature {
		log.Print("Tape (TZX) : Invalid TZX header signature")
		return false
	}
	index := 0
	for offset := 0; offset < tapeLength; {
		block := new(TapBlock)
		block.Type = data[offset]
		block.Index = index
		block.Offset = offset
		length, ok := tzxBlockLength(data[offset:])
		if !ok {
			log.Printf("Tape (TZX) : Invalid block #%d ID 0x%x at offset %d", index, block.Type, offset)
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

// tzxBlockLength returns the length of the block at the start of data. It is
// not valid if the block ID is unknown or the block is truncated.
func tzxBlockLength(data []byte) (int, bool) {
	// fixed part of the block (ID included) and its length field
	size, lenPos, lenBytes, lenFactor := 0, 0, 0, 1
	switch data[0] {
	case 0x10: // Standard speed data
		size, lenPos, lenBytes = 5, 3, 2
	case 0x11: // Turbo speed data
		size, lenPos, lenBytes = 19, 16, 3
	case 0x12: // Pure tone
		size = 5
	case 0x13: // Pulse sequence
		size, lenPos, lenBytes, lenFactor = 2, 1, 1, 2
	case 0x14: // Pure data
		size, lenPos, lenBytes = 11, 8, 3
	case 0x15: // Direct data
		size, lenPos, lenBytes = 9, 6, 3
	case 0x18: // CSW recording
		size, lenPos, lenBytes = 5, 1, 4
	case 0x19: // Generalized data
		size, lenPos, lenBytes = 5, 1, 4
	case 0x20: // Pause (silence) or 'Stop the Tape' command
		size = 3
	case 0x21: // Group Start
		size, lenPos, lenBytes = 2, 1, 1
	case 0x22: // Group End
		size = 1
	case 0x23: // Jump to Block
		size = 3
	case 0x24: // Loop Start
		size = 3
	case 0x25: // Loop End
		size = 1
	case 0x26: // Call Sequence
		size, lenPos, lenBytes, lenFactor = 3, 1, 2, 2
	case 0x27: // Return from Sequence
		size = 1
	case 0x28: // Select Block
		size, lenPos, lenBytes = 3, 1, 2
	case 0x2A: // Stop the tape if in 48K mode
		size = 5
	case 0x2B: // Set Signal Level
		size = 6
	case 0x30: // Text Description
		size, lenPos, lenBytes = 2, 1, 1
	case 0x31: // Message
		size, lenPos, lenBytes = 3, 2, 1
	case 0x32: // Archive Info
		size, lenPos, lenBytes = 3, 1, 2
	case 0x33: // Hardware Type
		size, lenPos, lenBytes, lenFactor = 2, 1, 1, 3
	case 0x35: // Custom Info
		size, lenPos, lenBytes = 21, 17, 4
	case 'Z': // Glue Block
		size = tzxHeaderLength
	default:
		return 0, false
	}
	if size > len(data) {
		return 0, false
	}
	length := size
	if lenBytes > 0 {
		length += lenFactor * readIntN(data, lenPos, lenBytes)
	}
	if length > len(data) {
		return 0, false
	}
	return length, true
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
			if tzxLogAllBlocks {
				log.Printf("Tape (TZX) : Playing block #%d - 0x%x", control.BlockIndex, control.Block.Info().Type)
			}
		}

	case tapeStatePilot:
		control.Ear ^= tape.LevelMask
		control.State = tapeStatePilotNc

	case tapeStatePilotNc:
		tzx.pilotPulses--
		if tzx.pilotPulses > 0 {
			control.Timeout = tzx.pilotTiming
			control.State = tapeStatePilot
		} else {
			control.Timeout = tzx.sync1Timing
			control.State = tapeStateSync
		}

	case tapeStateSync:
		control.Ear ^= tape.LevelMask
		control.Timeout = tzx.sync2timing
		control.State = tapeStateByte

	case tapeStateByteNc:
		control.Ear ^= tape.LevelMask
		control.State = tapeStateByte

	case tapeStateByte:
		if tzx.blockLength <= 0 { // no data
			control.State = tapeStateLastPulse
			break
		}
		tzx.bitMask = 0x80
		tzx.lastBit = 0x00
		if tzx.blockLength == 1 {
			tzx.lastBit = 0x80 >> tzx.bitsLastByte
		}
		control.State = tapeStateBit1

	case tapeStateBit1:
		control.Ear ^= tape.LevelMask
		if (control.DataAtPos() & tzx.bitMask) == 0 {
			tzx.bitTime = tzx.zeroTiming
		} else {
			tzx.bitTime = tzx.oneTiming
		}
		control.Timeout = tzx.bitTime
		control.State = tapeStateBit2

	case tapeStateBit2:
		control.Ear ^= tape.LevelMask
		control.Timeout = tzx.bitTime
		tzx.bitMask >>= 1
		if tzx.bitMask == tzx.lastBit {
			control.BlockPos++
			tzx.blockLength--
			if tzx.blockLength > 0 {
				control.State = tapeStateByte
			} else {
				control.State = tapeStateLastPulse
			}
		} else {
			control.State = tapeStateBit1
		}

	case tapeStateLastPulse:
		control.Ear ^= tape.LevelMask
		if tzx.endBlockPause > 0 {
			control.Timeout = tapeTimingEoB
			control.State = tapeStatePause
		} else {
			control.State = tapeStateTzxHeader
		}

	case tapeStatePureTone:
		control.Ear ^= tape.LevelMask
		control.State = tapeStatePureToneNc

	case tapeStatePureToneNc:
		if tzx.pilotPulses > 0 {
			tzx.pilotPulses--
			control.Timeout = tzx.pilotTiming
			control.State = tapeStatePureTone
		} else {
			control.State = tapeStateTzxHeader
		}

	case tapeStatePulseSeq:
		control.Ear ^= tape.LevelMask
		control.State = tapeStatePulseSeqNc

	case tapeStatePulseSeqNc:
		if tzx.pilotPulses > 0 {
			tzx.pilotPulses--
			control.Timeout = tzxReadPulse(control.Block.Data(), control.BlockPos)
			control.BlockPos += 2
			control.State = tapeStatePulseSeq
		} else {
			control.State = tapeStateTzxHeader
		}

	case tapeStatePause:
		control.Ear = tzxStartEar
		if !control.EndOfTape() {
			control.Timeout = tzx.endBlockPause * tapeTimingEoB
			control.State = tapeStateTzxHeader
		} else {
			control.State = tapeStateStop
		}

	case tapeStatePauseStop:
		if tzx.endBlockPause > 0 {
			control.Ear = tzxStartEar
			control.Timeout = tzx.endBlockPause * tapeTimingEoB
			control.State = tapeStateTzxHeader
		} else {
			tzx.stopTape(control)
		}

	case tapeStateStop:
		control.Playing = false // Stop

	default:
		control.State = tapeStateStop
	}
}

// tzxReadPulse reads a pulse length in Ts (at least 1 Ts)
func tzxReadPulse(data []byte, pos int) int {
	if pulse := readInt(data, pos); pulse > 0 {
		return pulse
	}
	return 1
}

// tzxLastByteBits returns the valid used bits of the last byte (1..8)
func tzxLastByteBits(bits byte) byte {
	if bits == 0 || bits > 8 {
		return 8
	}
	return bits
}

// stopTape stops the playback (stop the tape blocks). Playing again continues
// with the next block.
func (tzx *Tzx) stopTape(control *tape.Control) {
	control.Playing = false
	control.State = tapeStateTzxHeader
}

func (tzx *Tzx) parseHeader(control *tape.Control) {
	data := control.Block.Data()
	id := control.Block.Info().Type
	switch id {

	case 0x10:
		tzx.pilotTiming = tapeTimingPilot
		tzx.sync1Timing = tapeTimingSync1
		tzx.sync2timing = tapeTimingSync2
		tzx.zeroTiming = tapeTimingZero
		tzx.oneTiming = tapeTimingOne
		tzx.bitsLastByte = 8
		tzx.endBlockPause = readInt(data, control.BlockPos+1)
		tzx.blockLength = readInt(data, control.BlockPos+3)
		control.BlockPos += 5
		if tzx.blockLength > 0 && control.DataAtPos() >= 0x80 {
			tzx.pilotPulses = tapeDataPulses
		} else {
			tzx.pilotPulses = tapeHeaderPulses
		}
		control.State = tapeStatePilotNc
		control.BlockIndex++
		// log : standard data block
		if tzx.blockLength >= tapHeaderLength && data[5] == 0 {
			switch data[6] {
			case 0:
				log.Println("Tape (TZX) : Program header block:", readString(data, 7, 10))
			case 3:
				log.Println("Tape (TZX) : Bytes header block:", readString(data, 7, 10))
			default:
				log.Println("Tape (TZX) : Standard header block")
			}
		} else {
			log.Println("Tape (TZX) : Standard data block:", tzx.blockLength, "bytes")
		}

	case 0x11:
		tzx.pilotTiming = tzxReadPulse(data, control.BlockPos+1)
		tzx.sync1Timing = tzxReadPulse(data, control.BlockPos+3)
		tzx.sync2timing = tzxReadPulse(data, control.BlockPos+5)
		tzx.zeroTiming = tzxReadPulse(data, control.BlockPos+7)
		tzx.oneTiming = tzxReadPulse(data, control.BlockPos+9)
		tzx.pilotPulses = readInt(data, control.BlockPos+11)
		tzx.bitsLastByte = tzxLastByteBits(data[control.BlockPos+13])
		tzx.endBlockPause = readInt(data, control.BlockPos+14)
		tzx.blockLength = readIntN(data, control.BlockPos+16, 3)
		control.BlockPos += 19
		control.State = tapeStatePilotNc
		control.BlockIndex++
		// log : turbo data block
		log.Println("Tape (TZX) : Turbo data block:", tzx.blockLength, "bytes")

	case 0x12: // Pure Tone Block
		tzx.pilotTiming = tzxReadPulse(data, control.BlockPos+1)
		tzx.pilotPulses = readInt(data, control.BlockPos+3)
		control.BlockPos += 5
		control.State = tapeStatePureToneNc
		control.BlockIndex++

	case 0x13: // Pulse Sequence Block
		tzx.pilotPulses = int(data[control.BlockPos+1])
		control.BlockPos += 2
		control.State = tapeStatePulseSeqNc
		control.BlockIndex++

	case 0x14: // Pure Data Block
		tzx.zeroTiming = tzxReadPulse(data, control.BlockPos+1)
		tzx.oneTiming = tzxReadPulse(data, control.BlockPos+3)
		tzx.bitsLastByte = tzxLastByteBits(data[control.BlockPos+5])
		tzx.endBlockPause = readInt(data, control.BlockPos+6)
		tzx.blockLength = readIntN(data, control.BlockPos+8, 3)
		control.BlockPos += 11
		control.State = tapeStateByteNc
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
		// relative jump : signed offset
		target := int(int16(readWord(data, control.BlockPos+1)))
		if target == 0 { // invalid jump to itself : next block
			target = 1
		}
		control.BlockIndex += target
		if control.BlockIndex < 0 { // invalid jump before the first block
			control.BlockIndex = 0
		}

	case 0x24: // Loop Start
		tzx.loopCount = readInt(data, control.BlockPos+1)
		control.BlockIndex++
		tzx.loopStart = control.BlockIndex

	case 0x25: // Loop End
		tzx.loopCount--
		if tzx.loopCount <= 0 { // loop done (or no loop start)
			control.BlockIndex++
		} else {
			control.BlockIndex = tzx.loopStart
		}

	case 0x28: // Select Block
		control.BlockIndex++

	case 0x2A: // Stop the tape if in 48K mode
		control.BlockIndex++
		tzx.stopTape(control)

	case 0x2B: // Set Signal Level
		if data[control.BlockPos+5] == 0 {
			control.Ear = tape.LevelLow
		} else {
			control.Ear = tape.LevelHigh
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
		log.Printf("Tape (TZX) : Error at block #%d: Unsupported type: 0x%x", control.BlockIndex, id)
		control.BlockIndex++
	}
}
