package audio

import "testing"

// AY master clocks per envelope step with envelope period 1
const ayEnvelopeStepClocks = 16

func newTestAY() *AY38910 {
	ay := NewAY38910(NewConfig(48000, 50, 2496))
	ay.Init()
	return ay
}

// refEnvelope is the datasheet envelope: CONTINUE, ATTACK, ALTERNATE, HOLD
// bits, 16 steps per ramp.
func refEnvelope(shape byte, steps int) []uint8 {
	cont, att, alt, hold := shape&8 != 0, shape&4 != 0, shape&2 != 0, shape&1 != 0
	out := make([]uint8, steps)
	for i := range out {
		ramp, k := i/16, uint8(i%16)
		if ramp > 0 && (!cont || hold) {
			if cont && att != alt {
				out[i] = 15
			} else {
				out[i] = 0
			}
			continue
		}
		up := att
		if alt && ramp%2 == 1 {
			up = !up
		}
		if up {
			out[i] = k
		} else {
			out[i] = 15 - k
		}
	}
	return out
}

func TestAYEnvelopeShapes(t *testing.T) {
	const steps = 64
	for shape := byte(0); shape < 16; shape++ {
		ay := newTestAY()
		ay.WriteRegister(AY38910EnvelopeFrequencyLow, 1)
		ay.WriteRegister(AY38910EnvelopeShape, shape)
		want := refEnvelope(shape, steps)
		for i := 0; i < steps; i++ {
			if got := ay.envelope.volume; got != want[i] {
				t.Errorf("shape %04b step %d: volume = %d, want %d", shape, i, got, want[i])
				break
			}
			ay.Emulate(ayEnvelopeStepClocks)
		}
	}
}

func TestAYEnvelopePeriod(t *testing.T) {
	cases := []struct {
		name          string
		period        uint16
		clocksPerStep int
	}{
		{"period 0 acts as 1", 0, 16},
		{"period 1", 1, 16},
		{"period 3", 3, 48},
		{"period 256", 256, 4096},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ay := newTestAY()
			ay.WriteRegister(AY38910EnvelopeFrequencyLow, byte(c.period))
			ay.WriteRegister(AY38910EnvelopeFrequencyHigh, byte(c.period>>8))
			ay.WriteRegister(AY38910EnvelopeShape, 0x0d) // attack and hold
			ay.Emulate(c.clocksPerStep - 1)
			if v := ay.envelope.volume; v != 0 {
				t.Errorf("volume before first step = %d, want 0", v)
			}
			ay.Emulate(1)
			if v := ay.envelope.volume; v != 1 {
				t.Errorf("volume after %d clocks = %d, want 1", c.clocksPerStep, v)
			}
		})
	}
}

func TestAYEnvelopeShapeWriteRestarts(t *testing.T) {
	ay := newTestAY()
	ay.WriteRegister(AY38910EnvelopeFrequencyLow, 1)
	ay.WriteRegister(AY38910EnvelopeShape, 0x0c) // attack, continuous
	ay.Emulate(ayEnvelopeStepClocks * 5)
	if v := ay.envelope.volume; v != 5 {
		t.Fatalf("volume after 5 steps = %d, want 5", v)
	}
	ay.WriteRegister(AY38910EnvelopeShape, 0x08) // decay : restarts at 15
	if v := ay.envelope.volume; v != 15 {
		t.Errorf("volume after shape write = %d, want 15", v)
	}
	ay.Emulate(ayEnvelopeStepClocks)
	if v := ay.envelope.volume; v != 14 {
		t.Errorf("volume one step after shape write = %d, want 14", v)
	}
}

func TestAYResetEnvelopeIdle(t *testing.T) {
	ay := newTestAY()
	ay.WriteRegister(AY38910EnvelopeShape, 0x0c)
	ay.Emulate(ayEnvelopeStepClocks * 7)
	ay.Reset()
	ay.Emulate(ayEnvelopeStepClocks * 40)
	if v := ay.envelope.volume; v != 0 {
		t.Errorf("envelope volume after reset = %d, want idle 0", v)
	}
}

// A channel in envelope mode follows the envelope volume
func TestAYChannelEnvelopeVolume(t *testing.T) {
	ay := newTestAY()
	ay.WriteRegister(AY38910MixerControl, 0x3e)      // tone A only
	ay.WriteRegister(AY38910ChannelAFrequencyLow, 1) // toggle every 8 clocks
	ay.WriteRegister(AY38910ChannelAVolume, 0x10)    // envelope mode
	ay.WriteRegister(AY38910EnvelopeFrequencyLow, 1)
	ay.WriteRegister(AY38910EnvelopeShape, 0x0d) // attack and hold
	levels := map[uint16]bool{}
	for i := 0; i < 16*ayEnvelopeStepClocks/8; i++ {
		ay.Emulate(8)
		levels[ay.channelA.level] = true
	}
	if len(levels) < 8 {
		t.Errorf("channel A took %d distinct levels during an attack ramp, want the envelope steps", len(levels))
	}
	if ay.channelA.level != 0 && ay.channelA.level != ay.channelA.levels[15] {
		t.Errorf("channel A level after attack = %d, want 0 or max %d", ay.channelA.level, ay.channelA.levels[15])
	}
}
