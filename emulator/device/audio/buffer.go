package audio

// -----------------------------------------------------------------------------
// Buffer - Stereo samples buffer
// -----------------------------------------------------------------------------

const sampleSize = 0x04 // 2 bps * 2 channels

// Sample is a stereo sample
type Sample struct {
	Left  uint16
	Right uint16
}

// Buffer is a 16bit audio doble buffer : samples and audio data
type Buffer struct {
	samples []Sample // sample data u16 format
	data    []byte   // data buffer. Format : SDL AUDIO_U16LSB
}

// NewBuffer creates a new buffer of Freq and FPS
func NewBuffer(size int) *Buffer {
	buffer := new(Buffer)
	buffer.samples = make([]Sample, size)
	buffer.data = make([]byte, size*sampleSize)
	return buffer
}

// Samples gets the audio samples
func (buffer *Buffer) Samples() []Sample {
	return buffer.samples
}

// Size is the number of samples of the buffer
func (buffer *Buffer) Size() int { return len(buffer.samples) }

// Reset the samples buffer
func (buffer *Buffer) Reset() {
	for i := range buffer.samples {
		buffer.samples[i].Left = 0
		buffer.samples[i].Right = 0
	}
}

// Sample operations

// GetSample gets sample at index
func (buffer *Buffer) GetSample(index int) Sample {
	return buffer.samples[index]
}

// SetSample sets sample at index
func (buffer *Buffer) SetSample(index int, sample Sample) {
	buffer.samples[index].Left = sample.Left
	buffer.samples[index].Right = sample.Right
}

// AddSample adds a sample at index
func (buffer *Buffer) AddSample(index int, left uint16, right uint16) {
	buffer.samples[index].Left += left
	buffer.samples[index].Right += right
}

// Audio data buffer

// Data gets the audio data buffer. SDL AUDIO_U16LSB format.
func (buffer *Buffer) Data() []byte {
	return buffer.data
}

// BuildData builds the output audio buffer
func (buffer *Buffer) BuildData() {
	for i, j := 0, 0; i < buffer.Size(); i++ {
		sample := buffer.samples[i]
		high, low := uint8(sample.Left>>8), uint8(sample.Left&0xff)
		buffer.data[j] = low
		j++
		buffer.data[j] = high
		j++
		high, low = uint8(sample.Right>>8), uint8(sample.Right&0xff)
		buffer.data[j] = low
		j++
		buffer.data[j] = high
		j++
	}
}
