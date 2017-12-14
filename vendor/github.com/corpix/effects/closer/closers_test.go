package closer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type closer struct {
	closed bool
}

func (c *closer) Close() error {
	if c.closed {
		panic("already closed!")
	}

	c.closed = true
	return nil
}

func TestClosers(t *testing.T) {
	var (
		samples = []struct {
			name    string
			closers Closers
		}{
			{
				name:    "zero",
				closers: Closers{},
			},
			{
				name: "single",
				closers: Closers{
					&closer{},
				},
			},
			{
				name: "multiple",
				closers: Closers{
					&closer{},
					&closer{},
					&closer{},
					&closer{},
					&closer{},
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				assert.Equal(t, nil, sample.closers.Close())
				for _, c := range sample.closers {
					assert.Equal(t, true, c.(*closer).closed)
				}
			},
		)
	}
}
