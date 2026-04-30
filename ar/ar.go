// Package ar provides both a reader for Unix ar archives.
package ar

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// GlobalHeader contains the magic bytes that are written at the beginning
	// of an ar file (also known as armag).
	GlobalHeader = "!<arch>\n"
	// ThinArchiveGlobalHeader contains the magic bytes that are written at the
	// beginning of a thin ar archive, which is currently not supported.
	ThinArchiveGlobalHeader = "!<thin>\n"
	// HeaderTerminator contains the byte sequence with which file headers are
	// terminated (also known as ar_fmag).
	HeaderTerminator = "`\n"
	// HeaderSize holds the size of an ar header in bytes.
	HeaderSize = 60
)

var (
	// ErrInvalidGlobalHeader is returned by NewReader if the provided data does
	// not start with the correct ar global header.
	ErrInvalidGlobalHeader = fmt.Errorf("invalid global header")
)

const (
	nameFieldSize    = 16
	modTimeFieldSize = 12
	uidFieldSize     = 6
	gidFieldSize     = 6
	modeFiledSize    = 8
	sizeFieldSize    = 10

	padding                    = '\n'
	bsdExtendedFormatPrefix    = "#1/"
	gnuExtendedFormatNameTable = "//"
)

// Type specifies the ar archive variant with the corresponding extensions.
type Type int

var (
	// TypeBasic is a basic ar archive without any extensions that supports 16
	// character file names.
	TypeBasic Type = 0 //nolint:revive
	// TypeBSD is the ar archive type written by BSD's ar tool with support
	// for large file names.
	TypeBSD Type = 1
	// TypeGNU is the ar archive type written by GNU's or SYSTEM V's ar tool
	// with support for large file names.
	TypeGNU Type = 2
)

// Header holds an ar entry's metadata.
type Header struct {
	// Name holds the file name. If the name has more than 16 characters, it is
	// considered an extended name which will be written in BSD style.
	name string
	// ModTime holds the time stamp of the last file modification.
	modTime time.Time
	// UID holds the file owner's UID.
	uid int64
	// GID holds the fole owner's GID.
	gid int64
	// Mode holds the file's mode.
	mode os.FileMode
	// Size holds the file size of up to 9999999999 bytes. If Size is
	// UnknownSize, the file size will be automatically determined by the actual
	// bytes written for the entry.
	size int64
}

// Name returns the name field in the ar without any `/`
func (h *Header) Name() string {
	return filepath.Base(h.name)
}

// Path returns the full string in archive name field
func (h *Header) Path() string {
	return h.name
}

// Size returns size of file
func (h *Header) Size() int64 {
	return h.size
}

// Mode returns mod time of file
func (h *Header) Mode() os.FileMode {
	return h.mode
}

// ModTime returns mode of file
func (h *Header) ModTime() time.Time {
	return h.modTime
}

// Below methods are just needed to implement archive.Header interface

// Symlink always returns empty string.
// Needed to implement archive.Header interface
func (h *Header) Symlink() string {
	return ""
}

// IsDir always returns false b/c AR files dont have directories
// Needed to implement archive.Header interface
func (h *Header) IsDir() bool {
	return false
}

func (h *Header) Sys() any {
	return nil
}
