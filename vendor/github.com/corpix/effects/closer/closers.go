package closer

import (
	"io"
)

type Closers []io.Closer

func (c Closers) Close() error {
	var (
		n   = len(c) - 1
		err error
	)

	for n >= 0 {
		err = c[n].Close()
		if err != nil {
			return err
		}
		n--
	}

	return nil
}
