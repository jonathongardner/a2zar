package archive

import "errors"

// ErrInvalidMagic returned when a header has an invalid magic
var ErrInvalidMagic = errors.New("invalid magic")

// ErrUnknownFileType returned when an entry has an unknown/supported file type (i.e. regular, directory, symlink)
var ErrUnknownFileType = errors.New("unknown file type")

// ErrInvalidFileInfo returned if the parse file info is invalid (i.e. bad mod, modtime, etc)
var ErrInvalidFileInfo = errors.New("invalid fileinfo")

// ErrReaderNotOpen returned if the archive reader is not open yet
var ErrReaderNotOpen = errors.New("reader not open")

// ErrReaderClosed returned if the archive reader is closed
var ErrReaderClosed = errors.New("reader closed")
