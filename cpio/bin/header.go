package bin

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jonathongardner/a2zar/cpio/helpers"
	"github.com/jonathongardner/a2zar/internal/iio"
	"github.com/jonathongardner/a2zar/internal/utils"
)

var (
	magicBig    = [2]byte{113, 199} // 0x71 0xC7 uint16(070707)
	magicLittle = [2]byte{199, 113} // 0xC7 0x71 uint16(070707)
)

const (
	trailer = "TRAILER!!!"
	block   = 2
)

type Header struct {
	magic    uint16
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

// loadFields loads fields based on magic
// return io.EOF if header is trailer
func (nch *Header) loadFields(ar *Reader) error {
	shorts, err := readShorts(ar.r, 13)
	if err != nil {
		return err
	}
	// 0: Magic    0-2
	// 1: Dev      2-4
	// 2: Ino      4-6
	// 3: Mode     6-8
	// 4: Uid      8-10
	// 5: Gid      10-12
	// 6: Nlink    12-14
	// 7: Rdev     14-16
	// 8-9: Mtime    16-20
	// 10: NameSize 20-22
	// 11-12: FileSize 22-26

	nch.magic = shorts[0]
	nch.mode = utils.UnixToMode(int64(shorts[3]))
	nch.mtime = time.Unix(int64(toUint32(shorts[8:10])), 0)
	nch.filesize = int64(toUint32(shorts[11:13]))

	namesize := int64(shorts[10])
	// offset 26 for headers
	pr := iio.NewLimitPadReaderWithOffset(ar.r, namesize, 26, block)
	buf, err := io.ReadAll(pr)
	// if we get error no need to pad
	if err != nil {
		return err
	}
	if err := pr.Pad(); err != nil {
		return err
	}
	// subtract one for null byte
	nch.filename = string(buf[:namesize-1])

	if nch.isTrailer() {
		ar.pr = iio.NewClosedErrorReader()
		return io.EOF
	}

	ar.pr = iio.NewLimitPadReader(ar.r, nch.filesize, block)
	if nch.isSymlink() {
		if nch.filesize > 1024 {
			return helpers.ErrSymlinkToLarge
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
