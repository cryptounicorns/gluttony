package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatReader(t *testing.T) {
	var (
		tearsInRain = `I've seen things you people wouldn't believe. Attack ships on fire off the shoulder of Orion. I watched C-beams glitter in the dark near the Tannhäuser Gate. All those moments will be lost in time, like tears in rain. Time to die.`

		samples = []struct {
			name   string
			seed   []byte
			buf    []byte
			result []byte
		}{
			{
				name:   "empty seed",
				seed:   []byte{},
				buf:    make([]byte, 10),
				result: make([]byte, 10),
			},
			{
				name:   "seed same as buf",
				seed:   []byte(tearsInRain),
				buf:    make([]byte, len(tearsInRain)),
				result: []byte(tearsInRain),
			},
			{
				name:   "seed larger than buf",
				seed:   []byte(tearsInRain),
				buf:    make([]byte, 45),
				result: []byte(tearsInRain[:45]),
			},
			{
				name:   "seed repeats in buf",
				seed:   []byte(tearsInRain[:45]),
				buf:    make([]byte, 90),
				result: []byte(tearsInRain[:45] + tearsInRain[:45]),
			},
			{
				name:   "seed part repeats in buf",
				seed:   []byte(tearsInRain[:75]),
				buf:    make([]byte, 80),
				result: []byte(tearsInRain[:75] + tearsInRain[:5]),
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				_, err := NewRepeatReader(sample.seed).
					Read(sample.buf)
				if err != nil {
					t.Error(err)
					return
				}

				assert.Equal(t, sample.result, sample.buf)
			},
		)
	}
}
