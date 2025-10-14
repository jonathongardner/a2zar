package rno

import "github.com/jonathongardner/a2zar/archive"

func NewReader() errReader {
	return errReader{}
}

// errReader implements io.Reader. Always returns ErrFileReaderNotOpen
type errReader struct{}

func (er errReader) Read(p []byte) (int, error) {
	return 0, archive.ErrReaderNotOpen
}
