package archive

import (
	"io"
	"os"
)

// Header a fileinfo like interface with Path
type Header interface {
	os.FileInfo
	Path() string
	// returns symlink should be empty string if mode isnt symlink
	Symlink() string
}

// Reader a reader with `NextHeader` to iterate over entries in archives
type Reader interface {
	io.Reader
	// NextHeader returns the next Header interface. This is usually a wrapper around `Next()`
	NextHeader() (Header, error)
}
