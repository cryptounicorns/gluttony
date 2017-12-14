package reader

import (
	"io"
)

type RepeatReader struct {
	buf []byte
}

func (r *RepeatReader) Read(buf []byte) (int, error) {
	var (
		n           int
		copiedBytes int
	)

	for {
		n = copy(buf[copiedBytes:], r.buf)
		copiedBytes += n
		if n == 0 {
			break
		}
	}

	return copiedBytes, nil
}

func NewRepeatReader(buf []byte) io.Reader {
	return &RepeatReader{
		buf: buf,
	}
}
