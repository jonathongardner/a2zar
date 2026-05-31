package archive_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jonathongardner/a2zar/ar"
	"github.com/jonathongardner/a2zar/archive"
	"github.com/jonathongardner/a2zar/cpio/bin"
	"github.com/jonathongardner/a2zar/cpio/newc"
	"github.com/jonathongardner/a2zar/internal/test"
	"github.com/jonathongardner/a2zar/xar"
)

func TestArchiveReader(t *testing.T) {
	testCases := []struct {
		name       string
		readerFunc func(*os.File) (archive.Reader, error)
		exp        []entry
	}{
		{
			name: "xar",
			readerFunc: func(f *os.File) (archive.Reader, error) {
				return xar.NewReader(f)
			},
			exp: []entry{
				knownPaths["readme.md"],
				knownPaths["bar"],
				knownPaths["baz"],
				knownPaths["symlink"],
				knownPaths["chew"],
				knownPaths["foo"],
			},
		},
		{
			name: "ar",
			readerFunc: func(f *os.File) (archive.Reader, error) {
				return ar.NewReader(f)
			},
			exp: []entry{
				knownPaths["readme.md"],
				// ar doesnt have directories so its flat
				knownPaths["baz"].WithPath("baz"),
				knownPaths["chew"],
				knownPaths["foo"],
				// ar doesnt have symlink so it follows to the value
				knownPaths["foo"].WithPath("symlink"),
			},
		},
		{
			name: "newc.cpio",
			readerFunc: func(f *os.File) (archive.Reader, error) {
				return newc.NewReader(f)
			},
			exp: []entry{
				knownPaths["readme.md"],
				knownPaths["bar"],
				knownPaths["baz"],
				knownPaths["symlink"],
				knownPaths["chew"],
				knownPaths["foo"],
			},
		},
		{
			name: "bin.cpio",
			readerFunc: func(f *os.File) (archive.Reader, error) {
				return bin.NewReader(f)
			},
			exp: []entry{
				knownPaths["readme.md"],
				knownPaths["bar"],
				knownPaths["baz"],
				knownPaths["symlink"],
				knownPaths["chew"],
				knownPaths["foo"],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(test.LargeFile(fmt.Sprintf("test/golden-archive.%v", tc.name)))
			if err != nil {
				t.Fatalf("failed to open file: %v", err)
			}
			defer func() {
				err := file.Close()
				if err != nil {
					t.Fatalf("failed to close file: %v", err)
				}
			}()

			ar, err := tc.readerFunc(file)
			if err != nil {
				t.Fatal(fmt.Errorf("failed to create xar reader: %v", err))
			}

			assertEntries(t, tc.exp, entries(t, ar))
		})
	}
}
