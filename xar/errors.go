package xar

import "errors"

// ErrInvalidHeader returned if xar header is invalid
var ErrInvalidHeader = errors.New("invalid header")
var ErrInvalidToC = errors.New("invalid table of contents")

// ErrInvalidFileInfo returned if the parse file info is invalid (i.e. bad mod, modtime, etc)
var ErrInvalidTime = errors.New("invalid time format")
var ErrInvalidMode = errors.New("invalid mode format")

// ErrInvalidSize set when non links or directories has size not 0
var ErrInvalidSize = errors.New("invalid size")

// ErrUnknownFileEncoding returned when an entry uses as unknown/supported compression (i.e. none, gz, bz)
var ErrUnknownFileEncoding = errors.New("unknown file encoding")
