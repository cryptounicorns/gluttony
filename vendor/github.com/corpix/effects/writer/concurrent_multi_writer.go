package writer

import (
	"context"
	"io"
	"sync"

	"github.com/corpix/pool"
	"github.com/corpix/time"
)

// ConcurrentMultiWriterConfig is a configuration for ConcurrentMultiWriter.
// ScheduleTimeout is a timeout for worker pool to begin write data into
// the concrete writer.
// Pool is a configuration for github.com/corpix/pool
type ConcurrentMultiWriterConfig struct {
	ScheduleTimeout time.Duration
	Pool            pool.Config
}

// ConcurrentMultiWriter represents a dynamic list of Writer interfaces.
type ConcurrentMultiWriter struct {
	Config       ConcurrentMultiWriterConfig
	pool         *pool.Pool
	errorHandler func(*ConcurrentMultiWriter, io.Writer, error)
	locker       *sync.RWMutex
	writers      []io.Writer // FIXME: More effective solution?
}

// Add adds a Writer to the ConcurrentMultiWriter.
func (w *ConcurrentMultiWriter) Add(c io.Writer) {
	w.locker.Lock()
	defer w.locker.Unlock()

	w.writers = append(
		w.writers,
		c,
	)
}

// Remove removes Writer from the list if it exists
// and returns true in this case, otherwise it will be
// false.
func (w *ConcurrentMultiWriter) Remove(c io.Writer) bool {
	w.locker.Lock()
	defer w.locker.Unlock()

	for k, v := range w.writers {
		if v == c {
			if k < len(w.writers)-1 {
				w.writers = append(
					w.writers[0:k],
					w.writers[k+1:]...,
				)
			} else {
				w.writers = w.writers[0:k]
			}
			return true
		}
	}

	return false
}

// Has returns true if Writer exists in the list and false otherwise.
func (w *ConcurrentMultiWriter) Has(c io.Writer) bool {
	w.locker.RLock()
	defer w.locker.RUnlock()

	for _, v := range w.writers {
		if v == c {
			return true
		}
	}

	return false
}

func (w *ConcurrentMultiWriter) handleWorkError(wg *sync.WaitGroup, v io.Writer, err error) {
	defer wg.Done()

	w.errorHandler(w, v, err)
}

func (w *ConcurrentMultiWriter) newWork(buf []byte, v io.Writer, wg *sync.WaitGroup, cancel context.CancelFunc) pool.Executor {
	return func(ctx context.Context) {
		var (
			n   int
			err error
		)

		select {
		case <-ctx.Done():
			go w.handleWorkError(
				wg,
				v,
				NewErrScheduleTimeout(w.Config.ScheduleTimeout.Duration()),
			)
			return
		default:
			defer cancel()

			n, err = v.Write(buf)
			if err != nil {
				go w.handleWorkError(wg, v, err)
				return
			}

			if n != len(buf) {
				go w.handleWorkError(wg, v, io.ErrShortWrite)
				return
			}
		}

		// XXX: Done case should handled in `handleWorkError`
		// because we don't want to hold a slot in a worker pool
		// while somebody writing received error to logs or
		// doin other stuff with the error.
		wg.Done()
	}
}

// Write writes a buf to each concrete writer concurrently.
func (w *ConcurrentMultiWriter) Write(buf []byte) (int, error) {
	var (
		wg     = &sync.WaitGroup{}
		ctx    context.Context
		cancel context.CancelFunc
	)

	w.locker.RLock()
	wg.Add(len(w.writers))

	for _, v := range w.writers {
		ctx, cancel = context.WithTimeout(
			context.Background(),
			w.Config.ScheduleTimeout.Duration(),
		)

		w.pool.Feed <- pool.NewWork(
			ctx,
			w.newWork(buf, v, wg, cancel),
		)
	}

	w.locker.RUnlock()
	wg.Wait()

	return len(buf), nil
}

func (w *ConcurrentMultiWriter) Close() error {
	w.pool.Close()
	return nil
}

// NewConcurrentMultiWriter creates new Writer.
func NewConcurrentMultiWriter(c ConcurrentMultiWriterConfig, errorHandler func(*ConcurrentMultiWriter, io.Writer, error), ws ...io.Writer) *ConcurrentMultiWriter {
	var (
		writers = []io.Writer{}
	)

	return &ConcurrentMultiWriter{
		Config:       c,
		pool:         pool.NewFromConfig(c.Pool),
		errorHandler: errorHandler,
		locker:       &sync.RWMutex{},
		writers:      append(writers, ws...),
	}
}
