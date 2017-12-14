package writer

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiWriterWrite(t *testing.T) {
	var (
		samples = []struct {
			name    string
			data    []byte
			buffers []io.Writer
		}{
			{
				name:    "0",
				data:    []byte(tearsInRain),
				buffers: []io.Writer{},
			},
			{
				name: "1",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
				},
			},
			{
				name: "4",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				w := NewMultiWriter(sample.buffers...)
				n, err := w.Write(sample.data)
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.data), n)

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						sample.data,
						buf.(*bytes.Buffer).Bytes(),
					)
				}
			},
		)
	}
}

func TestMultiWriterAddWrite(t *testing.T) {
	var (
		samples = []struct {
			name    string
			data    []byte
			buffers []io.Writer
		}{
			{
				name:    "0",
				data:    []byte(tearsInRain),
				buffers: []io.Writer{},
			},
			{
				name: "1",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
				},
			},
			{
				name: "4",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				w := NewMultiWriter()

				for _, buf := range sample.buffers {
					w.Add(buf)
				}

				n, err := w.Write(sample.data)
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.data), n)

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						sample.data,
						buf.(*bytes.Buffer).Bytes(),
					)
				}
			},
		)
	}
}

func TestMultiWriterWriteRemoveWrite(t *testing.T) {
	var (
		samples = []struct {
			name    string
			data    []byte
			buffers []io.Writer
		}{
			{
				name:    "0",
				data:    []byte(tearsInRain),
				buffers: []io.Writer{},
			},
			{
				name: "1",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
				},
			},
			{
				name: "4",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				w := NewMultiWriter(sample.buffers...)
				n, err := w.Write(sample.data)
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.data), n)

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						sample.data,
						buf.(*bytes.Buffer).Bytes(),
					)
				}

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						true,
						w.Remove(buf),
					)
				}

				n, err = w.Write(sample.data)
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.data), n)

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						false,
						w.Has(buf),
					)
				}
			},
		)
	}
}

func TestMultiWriterAddWriteRemoveWrite(t *testing.T) {
	var (
		samples = []struct {
			name    string
			data    []byte
			buffers []io.Writer
		}{
			{
				name:    "0",
				data:    []byte(tearsInRain),
				buffers: []io.Writer{},
			},
			{
				name: "1",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
				},
			},
			{
				name: "4",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
					bytes.NewBuffer(nil),
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				w := NewMultiWriter()

				for _, buf := range sample.buffers {
					w.Add(buf)
				}

				n, err := w.Write(sample.data)
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.data), n)

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						sample.data,
						buf.(*bytes.Buffer).Bytes(),
					)
				}

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						true,
						w.Remove(buf),
					)
				}

				n, err = w.Write(sample.data)
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.data), n)

				for _, buf := range sample.buffers {
					assert.Equal(
						t,
						false,
						w.Has(buf),
					)
				}
			},
		)
	}
}
