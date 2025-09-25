package cpio

import (
	"errors"
	"io"

	"github.com/jonathongardner/a2zar/archive"
)

type cpioType int

const (
	NewC cpioType = iota
)

func (c cpioType) String() string {
	switch c {
	case NewC:
		return "NewC"
	default:
		panic("unknown cpio type")
	}
}

type Reader interface {
	archive.Reader
	// Type returns the cpio type
	Type() cpioType
}

// NewReader returns a cpio Reader
func NewReader(r io.Reader) (Reader, error) {
	return nil, errors.New("not supported")
}
