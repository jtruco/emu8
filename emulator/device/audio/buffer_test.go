package audio

import "testing"

func TestBufferAddSampleOutOfRange(t *testing.T) {
	buffer := NewBuffer(4)
	for _, index := range []int{-1, 4, 100} {
		buffer.AddSample(index, 1) // must not panic
	}
	buffer.AddSample(3, 7)
	if s := buffer.GetSample(3); s != 7 {
		t.Errorf("sample 3 = %d, want 7", s)
	}
}

// The last instruction of a frame can run a few Ts past the frame end
func TestAudioFrameOverrun(t *testing.T) {
	const tstates = 69888
	beeper := NewBeeper(NewConfig(48000, 50, tstates))
	beeper.Init()
	beeper.SetLevel(tstates+8, 1) // OUT at the end of the last instruction
	beeper.EndFrame()

	ay := NewAY38910(NewConfig(48000, 50, 2496))
	ay.Init()
	ay.Emulate(2496*8 + 64) // AY clocks of a frame and an overrun
	ay.EndFrame()
}
