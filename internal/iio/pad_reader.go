// internal io
package iio

import (
	"errors"
	"io"
)

type PadReader interface {
	io.Reader
	Pad() error
}

type padReader struct {
	r             io.Reader
	pad           int64
	remainingSize int64
}

// NewEmptyReader returns a empty limit reader (return io.ErrUnexpectedEOF if EOF early)
func NewEmptyReader() *padReader {
	return NewLimitPadReader(nil, 0, 1)
}

// NewLimitReader returns a limit reader (return io.ErrUnexpectedEOF if EOF early)
func NewLimitReader(r io.Reader, limit int64) *padReader {
	return NewLimitPadReader(r, limit, 1)
}

// NewLimitPadReaderWithOffset returns a limit reader (return io.ErrUnexpectedEOF if EOF early)
func NewLimitPadReaderWithOffset(r io.Reader, limit, offset, block int64) *padReader {
	if 0 > limit {
		panic("limit must be greater than 0")
	}
	if 0 > offset {
		panic("offset must be greater than 0")
	}
	if 1 > block {
		panic("block must be greater than 1")
	}

	return &padReader{r: r, pad: (block - ((limit + offset) % block)) % block, remainingSize: limit}
}

// NewLimitPadReader returns a limit reader (return io.ErrUnexpectedEOF if EOF early)
func NewLimitPadReader(r io.Reader, limit, block int64) *padReader {
	return NewLimitPadReaderWithOffset(r, limit, 0, block)
}

func (r *padReader) Read(buffer []byte) (n int, err error) {
	if r.remainingSize == 0 {
		return 0, io.EOF
	}

	if int64(len(buffer)) > r.remainingSize {
		buffer = buffer[:r.remainingSize]
	}

	n, err = r.r.Read(buffer)
	r.remainingSize -= int64(n)
	if err != nil {
		if err == io.EOF && r.remainingSize != 0 {
			err = io.ErrUnexpectedEOF
		}
		return n, err
	}

	return n, nil
}

func (pr *padReader) Drain() error {
	if _, err := io.Copy(io.Discard, pr); err != nil {
		return err
	}
	return nil
}

// Par will drain the limit reader, will return ErrUnexpectedEOF
// if doesnt have enough bytes to drain. If no error will discard
// pad bytes from reader
func (pr *padReader) Pad() error {
	if err := pr.Drain(); err != nil {
		return err
	}

	n, err := io.CopyN(io.Discard, pr.r, pr.pad)
	pr.pad -= int64(n)
	// dont return EOF
	if errors.Is(err, io.EOF) {
		err = nil
	}
	// if we didnt get pad expected return unexpected eof
	if err == nil && pr.pad != 0 {
		err = io.ErrUnexpectedEOF
	}
	return err
}
