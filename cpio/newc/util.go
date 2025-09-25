package newc

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jonathongardner/a2zar/internal/iio"
)

func readHex(r io.Reader, size int) ([]int64, error) {
	data := make([]byte, 8*size)
	if _, err := iio.ReadFull(r, data); err != nil {
		return nil, err
	}
	toReturn := make([]int64, 0, size)
	for len(data) > 0 {
		i, err := strconv.ParseInt(string(data[:8]), 16, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse hex field %d: %w", len(toReturn), err)
		}
		toReturn = append(toReturn, i)
		data = data[8:]
	}

	return toReturn, nil
}
