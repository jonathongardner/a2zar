package iio

import "github.com/jonathongardner/a2zar/archive"

func NewNotOpenErrorReader() errReader {
	return NewErrorReader(archive.ErrReaderNotOpen)
}
func NewClosedErrorReader() errReader {
	return NewErrorReader(archive.ErrReaderClosed)
}

func NewErrorReader(err error) errReader {
	return errReader{}
}

// errReader implements io.Reader. Always returns ErrFileReaderNotOpen
type errReader struct {
	err error
}

func (er errReader) Read(p []byte) (int, error) {
	return 0, er.err
}

// Pad nops
func (er errReader) Pad() error {
	return nil
}
