package newc

import (
	"os"
	"path/filepath"
	"time"
)

var (
	magic         = [6]byte{48, 55, 48, 55, 48, 49} // "070701"
	extendedMagic = [6]byte{48, 55, 48, 55, 48, 88} // "07070X"
)

type Header struct {
	magic    [6]byte
	mode     os.FileMode
	mtime    time.Time
	filesize int64

	filename string
	symlink  string
}

// ---------------FileInfo fields-------------
func (nch *Header) IsDir() bool {
	return nch.mode.IsDir()
}
func (nch *Header) Mode() os.FileMode {
	return nch.mode
}

func (nch *Header) ModTime() time.Time {
	return nch.mtime
}
func (nch *Header) Name() string {
	return filepath.Base(nch.filename)
}
func (nch *Header) Path() string {
	return nch.filename
}
func (nch *Header) Size() int64 {
	return nch.filesize
}
func (nch *Header) Symlink() string {
	return nch.symlink
}
func (nch *Header) Sys() any {
	return nil
}

//---------------FileInfo fields-------------

// isSymlink returns true if mode has symlink bit
func (nch *Header) isSymlink() bool {
	return nch.mode&os.ModeSymlink != 0
}

// isTrailer returns true if filename matches trailer name
func (nch *Header) isTrailer() bool {
	return nch.filename == trailer
}

// extended returns true if the magic matches the extended header (large than 4GB)
func (nch *Header) isExtended() bool {
	return nch.magic == extendedMagic
}

// normal returns true if the header matches the expected magic
func (nch *Header) isNormal() bool {
	return nch.magic == magic
}
