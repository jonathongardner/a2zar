package archive

import "errors"

// ErrUnknownFileType returned when an entry has an unknown/supported file type (i.e. regular, directory, symlink)
var ErrUnknownFileType = errors.New("unknown file type")

// ErrInvalidFileInfo returned if the parse file info is invalid (i.e. bad mod, modtime, etc)
var ErrInvalidFileInfo = errors.New("invalid fileinfo")

// ErrReaderNotOpen returned if the archive reader is not open yet
var ErrReaderNotOpen = errors.New("reader not open")
