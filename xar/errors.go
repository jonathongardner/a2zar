package xar

import "errors"

// ErrInvalidHeader returned if xar header is invalid
var ErrInvalidHeader = errors.New("invalid header")
var ErrInvalidToC = errors.New("invalid table of contents")

// ErrInvalidFileInfo returned if the parse file info is invalid (i.e. bad mod, modtime, etc)
var ErrInvalidFileInfo = errors.New("invalid fileinfo")
var ErrInvalidTime = errors.New("invalid time format")
var ErrInvalidMode = errors.New("invalid mode format")

// ErrInvalidSize set when non links or directories has size not 0
var ErrInvalidSize = errors.New("invalid size")

// ErrUnknownFileEncoding returned when an entry uses as unknown/supported compression (i.e. none, gz, bz)
var ErrUnknownFileEncoding = errors.New("unknown file encoding")

// ErrUnknownFileType returned when an entry has an unknown/supported file type (i.e. regular, directory, symlink)
var ErrUnknownFileType = errors.New("unknown file type")

// errReader implements io.Reader. Always returns ErrFileReaderNotOpen
type errReader struct{}

var ErrFileReaderNotOpen = errors.New("file reader not open")

func (er errReader) Read(p []byte) (int, error) {
	return 0, ErrFileReaderNotOpen
}
