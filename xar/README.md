# XAR

XAR is a Go package for reading XAR files. 

## Example
```golang
package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jonathongardner/a2zar/xar"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(fmt.Errorf("failed to open file: %v", err))
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file: %v", err))
		}
	}()

	xr, err := xar.NewReader(file)
	if err != nil {
		panic(fmt.Errorf("failed to create xar reader: %v", err))
	}

	for {
		fi, err := xr.Next()
		if err == io.EOF {
			break
		}
		if errors.Is(err, xar.ErrUnknownFileType) {
			fmt.Printf("name: %s, err: %v\n", fi.Path(), err)
			continue
		}
		if err != nil {
			panic(fmt.Errorf("failed to read file: %v", err))
		}

		fmt.Printf("name: %s, mode: %v, size: %d, errs: %v\n", fi.Path(), fi.Mode(), fi.Size(), fi.errs)
		// if symlink, print the target
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			fmt.Print("  symlink: %v\n", fi.Symlink())
		}

		// check if any parsing errors occurred building the header
		if perr := fi.ParsingError(); perr != nil {
			fmt.Printf("  parsing error: %v\n", perr)
		}

		if fi.Mode().IsRegular() {
			// print the sha1 of the xr reader
			h := sha1.New()
			if _, err := io.Copy(h, xr); err != nil {
				panic(fmt.Errorf("failed to compute sha1: %v", err))
			}
			fmt.Printf("  sha1: %x\n", h.Sum(nil))
		}
	}
}
```

Cloned from https://github.com/arelate/xargon
