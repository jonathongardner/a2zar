package xar

import (
	"bytes"
	"compress/bzip2"
	"compress/zlib"
	"crypto/x509"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/jonathongardner/a2zar/archive"
)

const (
	applicationOctetStreamMimeType = "application/octet-stream"
	applicationGzipMimeType        = "application/x-gzip"
	applicationBzip2MimeType       = "application/x-bzip2"
)

type XarReader struct {
	reader     io.ReaderAt
	header     *xarHeader
	toc        *xarToc
	fileInfo   []*XarFileInfo
	fileReader io.Reader
	heapOffset int64
	position   int64
	certs      []*x509.Certificate
}

func NewReader(reader io.ReaderAt) (*XarReader, error) {
	xr := &XarReader{
		reader:     reader,
		position:   0,
		fileReader: errReader{},
	}

	// xar file = Header + TOC + Heap
	if err := xr.readHeader(); err != nil {
		return nil, err
	}

	if err := xr.readToc(); err != nil {
		return nil, err
	}

	xr.heapOffset = xarHeaderSize + int64(xr.header.tocLengthCompressed)

	// iterate over the toc files and build a list of file info
	// so we can iterate over them in order
	xr.fileInfo = make([]*XarFileInfo, 0, len(xr.toc.Toc.File))

	for _, file := range xr.toc.Toc.File {
		xr.fileInfo = append(xr.fileInfo, file.fileInfo("")...)
	}

	return xr, nil
}

func (xr *XarReader) readHeader() error {
	h := make([]byte, xarHeaderSize)
	if n, err := xr.reader.ReadAt(h, 0); err != nil {
		return err
	} else if n != xarHeaderSize {
		return fmt.Errorf("%w wrong size", ErrInvalidHeader)
	}

	xh := &xarHeader{
		magic:                 binary.BigEndian.Uint32(h[0:4]),
		size:                  binary.BigEndian.Uint16(h[4:6]),
		version:               binary.BigEndian.Uint16(h[6:8]),
		tocLengthCompressed:   binary.BigEndian.Uint64(h[8:16]),
		tocLengthUncompressed: binary.BigEndian.Uint64(h[16:24]),
		checksumAlgorithm:     binary.BigEndian.Uint32(h[24:28]),
	}

	// validate expected values

	if xh.magic != xarHeaderMagic {
		return fmt.Errorf("%w unexpected format", ErrInvalidHeader)
	}

	if xh.size != xarHeaderSize {
		return fmt.Errorf("%w unexpected size", ErrInvalidHeader)
	}

	if xh.version != xarHeaderVersion {
		return fmt.Errorf("%w unsupported version", ErrInvalidHeader)
	}

	if xh.tocLengthCompressed == 0 {
		return fmt.Errorf("%w unexpected Table of Contents compressed length", ErrInvalidHeader)
	}

	if xh.tocLengthUncompressed == 0 {
		return fmt.Errorf("%w unexpected Table of Contents uncompressed length", ErrInvalidHeader)
	}

	if xh.checksumAlgorithm == 3 {
		return fmt.Errorf("%w unsupported checksum algorithm", ErrInvalidHeader)
	}

	xr.header = xh

	return nil
}

func (xr *XarReader) readToc() error {
	toc := make([]byte, xr.header.tocLengthCompressed)
	if n, err := xr.reader.ReadAt(toc, xarHeaderSize); err != nil {
		return err
	} else if uint64(n) != xr.header.tocLengthCompressed {
		return fmt.Errorf("%w size mismatch", ErrInvalidToC)
	}

	br := bytes.NewReader(toc)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return err
	}

	return xml.NewDecoder(zr).Decode(&xr.toc)
}

// NextHeader reads the next file in the xar archive and returns its information.
// returns io.EOF when there are no more files to read.
// returns an error if the file cannot be read
// returns an error if the file type is unknown
func (xr *XarReader) NextHeader() (archive.Header, error) {
	return xr.Next()
}

// Next reads the next file in the xar archive and returns its information.
// returns io.EOF when there are no more files to read.
// returns an error if the file cannot be read
// returns an error if the file type is unknown
func (xr *XarReader) Next() (*XarFileInfo, error) {
	if xr.position >= int64(len(xr.fileInfo)) {
		xr.fileReader = errReader{}
		return nil, io.EOF
	}

	fileInfo := xr.fileInfo[xr.position]

	xr.position++

	var err error
	xr.fileReader, err = xr.openFile(fileInfo.file)
	if err != nil {
		return nil, err
	}

	return fileInfo, fileInfo.headerErrs()
}

func (xr *XarReader) openFile(xf *xarFile) (io.Reader, error) {
	sectionReader := io.NewSectionReader(xr.reader, xf.Data.Offset+xr.heapOffset, xf.Data.Length)

	enc := xf.Data.Encoding.Style
	switch enc {
	case applicationOctetStreamMimeType, "":
		return sectionReader, nil
	case applicationGzipMimeType:
		return zlib.NewReader(sectionReader)
	case applicationBzip2MimeType:
		return bzip2.NewReader(sectionReader), nil
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownFileEncoding, enc)
	}
}

// Read reads the file content and returns it as an io.Reader.
func (xr *XarReader) Read(b []byte) (int, error) {
	return xr.fileReader.Read(b)
}
