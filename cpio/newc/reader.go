package newc

import (
	"fmt"
	"io"
	"time"

	"github.com/jonathongardner/a2zar/archive"
	"github.com/jonathongardner/a2zar/internal/iio"
	"github.com/jonathongardner/a2zar/internal/utils"
)

const (
	trailer = "TRAILER!!!"
	block   = 4
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
	if err := r.loadFields(h); err != nil {
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

// loadFields loads fields based on magic
// return io.EOF if header is trailer
func (r *Reader) loadFields(nch *Header) error {
	// if no bytes are it will return io.ErrUnexpectedEOF
	if _, err := iio.ReadFull(r.r, nch.magic[:]); err != nil {
		return err
	}

	if nch.isNormal() {
		return r.loadNormalFields(nch)
	}
	return fmt.Errorf("not supported yet")
}

// loadNormalFields loads the normal fields
// return io.EOF it header is trailer
func (ar *Reader) loadNormalFields(nch *Header) error {
	fields, err := readHex(ar.r, 13)
	if err != nil {
		return err
	}
	// 	0: ino, // 0-8
	// 	1: mode, // 8-16
	// 	2: uid, // 16-24
	// 	3: gid, // 24-32
	// 	4: nlink, // 32-40
	// 	5: mtime, // 40-48
	// 	6: filesize, //48-56
	// 	7: devmajor, // 56-64
	// 	8: devminor, // 64-72
	// 	9: rdevmajor, // 72-80
	// 	10: rdevminor, // 80-88
	// 	11: namesize, // 88-96
	// 	12: check, // 96-104

	nch.mode = utils.UnixToMode(fields[1])
	nch.mtime = time.Unix(fields[5], 0)
	nch.filesize = fields[6]

	namesize := fields[11]
	// offset 110 (6 + 104) for headers
	pr := iio.NewLimitPadReaderWithOffset(ar.r, namesize, 110, block)
	buf, err := io.ReadAll(pr)
	// if we get error no need to pad
	if err != nil {
		return err
	}

	if err := pr.Pad(); err != nil {
		return err
	}

	nch.filename = string(buf[:namesize-1])

	if nch.isTrailer() {
		ar.pr = iio.NewClosedErrorReader()
		return io.EOF
	}

	ar.pr = iio.NewLimitPadReader(ar.r, nch.filesize, block)
	if nch.isSymlink() {
		if nch.filesize > 1024 {
			return ErrSymlinkToLarge
		}

		// pr will return io.UnexpectedEOF if EOF is returned to early
		symlink, err := io.ReadAll(ar.pr)
		if err != nil {
			return err
		}
		if err := ar.pr.Pad(); err != nil {
			return err
		}
		nch.symlink = string(symlink)
		// TODO: think about returning error reader
		ar.pr = iio.NewEmptyReader()
		nch.filesize = 0
	}

	return nil
}
