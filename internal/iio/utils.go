package iio

import "io"

// ReadFull similar to io.ReadFull but always returns io.ErrUnexpectedEOF instead EOF
func ReadFull(r io.Reader, buf []byte) (int, error) {
	n, err := io.ReadFull(r, buf)
	if n != len(buf) && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return n, err
}
