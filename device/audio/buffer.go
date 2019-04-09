package audio

// -----------------------------------------------------------------------------
// Buffer & Player
// -----------------------------------------------------------------------------

// Buffer is a 16bit audio doble buffer : samples and audio data
type Buffer struct {
	frequency int      // Frequency of audio
	fps       int      // Frames per second
	size      int      // buffer size in samples
	samples   []uint16 // sample data u16 format
	data      []byte   // data buffer. Format : SDL AUDIO_U16LSB
}

// NewBuffer creates a new buffer of Freq and FPS
func NewBuffer(frequency, fps int) *Buffer {
	buffer := &Buffer{}
	buffer.frequency = frequency
	buffer.fps = fps
	buffer.size = frequency / fps
	buffer.samples = make([]uint16, buffer.size)
	buffer.data = make([]byte, buffer.size*2) // 2 bps
	return buffer
}

// AddSample adds a sample at index
func (buffer *Buffer) AddSample(index int, sample uint16) {
	buffer.samples[index] += sample
}

// BuildData builds audio buffer data
func (buffer *Buffer) BuildData() {
	for i, j := 0, 0; i < buffer.size; i++ {
		sample := buffer.samples[i]
		high, low := uint8(sample>>8), uint8(sample&0xff)
		buffer.data[j] = low
		j++
		buffer.data[j] = high
		j++
	}
}

// Data gets the audio data buffer. SDL AUDIO_U16LSB format.
func (buffer *Buffer) Data() []byte {
	return buffer.data
}

// FPS of the audio buffer
func (buffer *Buffer) FPS() int {
	return buffer.fps
}

// Frequency of the audio buffer
func (buffer *Buffer) Frequency() int {
	return buffer.frequency
}

// GetSample gets sample at index
func (buffer *Buffer) GetSample(index int) uint16 {
	return buffer.samples[index]
}

// Reset the samples buffer
func (buffer *Buffer) Reset() {
	for i := range buffer.samples {
		buffer.samples[i] = 0
	}
}

// Samples gets the audio samples
func (buffer *Buffer) Samples() []uint16 {
	return buffer.samples
}

// SetSample sets sample at index
func (buffer *Buffer) SetSample(index int, sample uint16) {
	buffer.samples[index] = sample
}

// Size number of samples of the buffer
func (buffer *Buffer) Size() int {
	return buffer.size
}
