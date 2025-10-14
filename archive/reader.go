package archive

import (
	"io"
	"os"
)

// Header a fileinfo like interface with Path
type Header interface {
	os.FileInfo
	Path() string
	Symlink() string
}

// Reader a reader with `NextHeader` to iterate over entries in archives
type Reader interface {
	io.Reader
	// NextHeader returns the next Header interface
	NextHeader() (Header, error)
}
