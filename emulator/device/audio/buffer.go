package audio

// -----------------------------------------------------------------------------
// Buffer - Stereo samples buffer
// -----------------------------------------------------------------------------

// Sample is a mono 16bit audio sample
type Sample = uint16

// Buffer is a 16bit audio doble buffer : samples and audio data
type Buffer struct {
	samples []Sample // sample data u16 format
	data    []byte   // data buffer. Format : SDL AUDIO_U16LSB
	filter  Filter   // audio sample filter
}

// NewBuffer creates a new buffer of Freq and FPS
func NewBuffer(size int) *Buffer {
	buffer := new(Buffer)
	buffer.samples = make([]Sample, size)
	buffer.data = make([]byte, size*2) // 2bps mono
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
		buffer.samples[i] = 0
	}
}

// Filter returns the curren audio filter
func (buffer *Buffer) Filter() Filter { return buffer.filter }

// SetFilter sets the audio filter
func (buffer *Buffer) SetFilter(filter Filter) { buffer.filter = filter }

// Sample operations

// GetSample gets audio sample at index
func (buffer *Buffer) GetSample(index int) Sample {
	return buffer.samples[index]
}

// SetSample sets the audio sample at index
func (buffer *Buffer) SetSample(index int, sample Sample) {
	buffer.samples[index] = sample
}

// AddSample adds (and apply filter) an audio sample at buffer
func (buffer *Buffer) AddSample(index int, sample Sample) {
	if buffer.filter != nil {
		sample = buffer.filter.Add(sample)
	}
	buffer.samples[index] += sample
}

// Audio data buffer

// Data gets the audio data buffer. SDL AUDIO_U16LSB format.
func (buffer *Buffer) Data() []byte {
	return buffer.data
}

// BuildData builds the output audio buffer
func (buffer *Buffer) BuildData() {
	// 16bit mono audio buffer
	for i, j := 0, 0; i < buffer.Size(); i++ {
		sample := buffer.samples[i]
		high, low := uint8(sample>>8), uint8(sample&0xff)
		buffer.data[j] = low
		j++
		buffer.data[j] = high
		j++
	}
}
