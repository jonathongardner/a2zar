package bin

import (
	"fmt"
	"io"

	"github.com/jonathongardner/a2zar/archive"
	"github.com/jonathongardner/a2zar/internal/iio"
)

type Reader struct {
	r  io.Reader
	pr iio.PadReader
}

func NewReader(r io.Reader) (*Reader, error) {
	return &Reader{r: r, pr: iio.NewNotOpenErrorReader()}, nil
}

// Next return the next header
func (r *Reader) Next() (*Header, error) {
	if err := r.pr.Pad(); err != nil {
		return nil, fmt.Errorf("pad to next header: %w", err)
	}
	h := &Header{}

	// return io.EOF if header is trailer
	if err := h.loadFields(r); err != nil {
		return nil, err
	}

	return h, nil
}

func (r *Reader) Read(data []byte) (int, error) {
	return r.pr.Read(data)
}

// NextHeader wrapper around Next
func (r *Reader) NextHeader() (archive.Header, error) {
	return r.Next()
}
