package bin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/jonathongardner/a2zar/internal/iio"
)

// endian returns the byte order of the header based on magic
func endian(data []byte) (binary.ByteOrder, error) {
	switch {
	case bytes.Equal(data, magicBig[:]):
		return binary.BigEndian, nil
	case bytes.Equal(data, magicLittle[:]):
		return binary.LittleEndian, nil
	default:
		return nil, errors.New("unknown magic")
	}
}

func readShorts(r io.Reader, size int) ([]uint16, error) {
	header := make([]byte, size*2)
	// if no bytes are it will return io.ErrUnexpectedEOF
	if _, err := iio.ReadFull(r, header); err != nil {
		return nil, err
	}
	b, err := endian(header[0:2])
	if err != nil {
		return nil, err
	}

	toReturn := make([]uint16, 0, size)
	for len(header) > 0 {
		toReturn = append(toReturn, b.Uint16(header[:2]))
		header = header[2:]
	}

	return toReturn, nil
}

// ints converts the shorts into uint32
// following cpio binary format
func toUint32(nums []uint16) uint32 {
	return (uint32(nums[0]) << 16) | uint32(nums[1])
}
