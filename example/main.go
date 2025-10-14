package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jonathongardner/a2zar/archive"
	"github.com/jonathongardner/a2zar/xar"
)

func main() {
	if len(os.Args) != 3 {
		panic("requires two arguments example type file")
	}

	file, err := os.Open(os.Args[2])
	if err != nil {
		panic(fmt.Errorf("failed to open file: %v", err))
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file: %v", err))
		}
	}()

	var ar archive.Reader
	typ := os.Args[1]
	switch typ {
	case "xar":
		ar, err = xar.NewReader(file)
	default:
		panic(fmt.Errorf("unknown type %v", typ))
	}
	if err != nil {
		panic(fmt.Errorf("failed to create %v reader: %v", typ, err))
	}

	for {
		fi, err := ar.NextHeader()
		if err == io.EOF {
			break
		}
		if errors.Is(err, archive.ErrUnknownFileType) {
			fmt.Printf("name: %s, err: %v\n", fi.Path(), err)
			continue
		}
		if errors.Is(err, archive.ErrInvalidFileInfo) {
			fmt.Printf("name: %s, err: %v\n", fi.Path(), err)
			continue
		}
		if err != nil {
			panic(fmt.Errorf("failed to read file: %v", err))
		}

		fmt.Printf("name: %s, mode: %v, size: %d\n", fi.Path(), fi.Mode(), fi.Size())
		// if symlink, print the target
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			fmt.Printf("  symlink: %v\n", fi.Symlink())
			continue
		}
		if fi.Mode().IsDir() {
			continue
		}

		if fi.Mode().IsRegular() {
			// print the sha1 of the xr reader
			h := sha1.New()
			if _, err := io.Copy(h, ar); err != nil {
				panic(fmt.Errorf("failed to compute sha1: %v", err))
			}
			fmt.Printf("  sha1: %x\n", h.Sum(nil))
			continue
		}
	}
}
