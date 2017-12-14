package writer

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/corpix/pool"
	"github.com/corpix/time"
	"github.com/stretchr/testify/assert"
)

type failingWriter struct {
	wrote int
	err   error
}

func (w failingWriter) Write([]byte) (int, error) {
	return w.wrote, w.err
}

func TestConcurrentMultiWriterWrite(t *testing.T) {
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
				w := NewConcurrentMultiWriter(
					ConcurrentMultiWriterConfig{
						ScheduleTimeout: 1 * time.Second,
						Pool: pool.Config{
							Workers:   10,
							QueueSize: 10,
						},
					},
					func(w *ConcurrentMultiWriter, v io.Writer, err error) {
						panic(err)
					},
					sample.buffers...,
				)
				defer w.Close()
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

func TestConcurrentMultiWriterAddWrite(t *testing.T) {
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
				w := NewConcurrentMultiWriter(
					ConcurrentMultiWriterConfig{
						ScheduleTimeout: 1 * time.Second,
						Pool: pool.Config{
							Workers:   10,
							QueueSize: 10,
						},
					},
					func(w *ConcurrentMultiWriter, v io.Writer, err error) {
						panic(err)
					},
				)
				defer w.Close()

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

func TestConcurrentMultiWriterWriteRemoveWrite(t *testing.T) {
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
				w := NewConcurrentMultiWriter(
					ConcurrentMultiWriterConfig{
						ScheduleTimeout: 1 * time.Second,
						Pool: pool.Config{
							Workers:   10,
							QueueSize: 10,
						},
					},
					func(w *ConcurrentMultiWriter, v io.Writer, err error) {
						panic(err)
					},
					sample.buffers...,
				)
				defer w.Close()

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

func TestConcurrentMultiWriterAddWriteRemoveWrite(t *testing.T) {
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
				w := NewConcurrentMultiWriter(
					ConcurrentMultiWriterConfig{
						ScheduleTimeout: 1 * time.Second,
						Pool: pool.Config{
							Workers:   10,
							QueueSize: 10,
						},
					},
					func(w *ConcurrentMultiWriter, v io.Writer, err error) {
						panic(err)
					},
				)
				defer w.Close()

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

func TestConcurrentMultiWriterErrorHandler(t *testing.T) {
	var (
		samples = []struct {
			name    string
			data    []byte
			buffers []io.Writer
			errors  []error
		}{
			{
				name:    "0",
				data:    []byte(tearsInRain),
				buffers: []io.Writer{},
				errors:  []error{},
			},
			{
				name: "1",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					failingWriter{0, NewErrScheduleTimeout(time.Duration(0).Duration())},
				},
				errors: []error{
					NewErrScheduleTimeout(time.Duration(0).Duration()),
				},
			},
			{
				name: "4",
				data: []byte(tearsInRain),
				buffers: []io.Writer{
					failingWriter{0, NewErrScheduleTimeout(time.Duration(0).Duration())},
					failingWriter{0, io.ErrShortWrite},
					failingWriter{0, nil}, // XXX: 0 bytes wrote != len(data) -> io.ErrShortWrite
					failingWriter{0, io.ErrUnexpectedEOF},
				},
				errors: []error{
					NewErrScheduleTimeout(time.Duration(0).Duration()),
					io.ErrShortWrite,
					io.ErrShortWrite,
					io.ErrUnexpectedEOF,
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					errors = []error{}
					errs   = make(chan error, len(sample.errors))
					wg     = &sync.WaitGroup{}
					w      = NewConcurrentMultiWriter(
						ConcurrentMultiWriterConfig{
							ScheduleTimeout: 1 * time.Second,
							Pool: pool.Config{
								Workers:   10,
								QueueSize: 10,
							},
						},
						func(w *ConcurrentMultiWriter, v io.Writer, err error) {
							// Could block control flow return from Write(...)
							// Feature, not a bug :)
							errs <- err
						},
						sample.buffers...,
					)
				)
				defer w.Close()
				n, err := w.Write(sample.data)
				assert.Equal(t, nil, err)
				assert.Equal(t, len(sample.data), n)

				wg.Add(len(sample.errors))
				go func() {
					for {
						select {
						case v, ok := <-errs:
							if ok {
								errors = append(
									errors,
									v,
								)
								wg.Done()
							} else {
								return
							}
						}
					}
				}()

				wg.Wait()

				assert.Equal(t, len(sample.errors), len(errors))
				for _, v := range errors {
					assert.Contains(t, sample.errors, v)
				}
			},
		)
	}
}
