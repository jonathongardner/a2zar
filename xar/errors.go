package xar

import "errors"

var ErrInvalidHeader = errors.New("invalid header")
var ErrInvalidToC = errors.New("invalid table of contents")
var ErrInvalidTime = errors.New("invalid time format")
var ErrInvalidMode = errors.New("invalid mode format")
var ErrInvalidSize = errors.New("invalid size")
var ErrFileReaderNotOpen = errors.New("file reader not open")
var ErrUnknownFileEncoding = errors.New("unknown file encoding")
var ErrUnknownFileType = errors.New("unknown file type")

type errReader struct{}

func (er errReader) Read(p []byte) (int, error) {
	return 0, ErrFileReaderNotOpen
}
