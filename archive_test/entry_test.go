package archive_test

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jonathongardner/a2zar/archive"
	"github.com/jonathongardner/a2zar/internal/test"
)

var knownPaths = map[string]entry{
	"readme.md": {path: "README.md", size: 196, mode: 0644, mtime: time.Unix(1760403731, 0), sha1: "4b6dc437de95771df609a38bb7d102d79ad388d9"},
	"bar":       {path: "bar", mode: 0755 | os.ModeDir, mtime: time.Unix(1760403731, 0)},
	"baz":       {path: "bar/baz", size: 5, mode: 0644, mtime: time.Unix(1779560749, 0), sha1: "886f90cb542138934de905357d0fdbf35c6bff33"},
	"chew":      {path: "chew", size: 2048, mode: 0644, mtime: time.Unix(1779560734, 0), sha1: "c12306b17f72062188d6bbfe7a76f15945a8e1a6"},
	"foo":       {path: "foo", size: 512, mode: 0644, mtime: time.Unix(1779560751, 0), sha1: "aec2d949eea7b34ee3e91baf40a03879e59b2935"},
	"symlink":   {path: "bar/symlink", mode: 0777 | os.ModeSymlink, mtime: time.Unix(1760403731, 0), symlink: "../foo"},
}

// entry an expected archive entry
type entry struct {
	path    string
	size    int64
	mode    os.FileMode
	mtime   time.Time
	sha1    string
	symlink string
}

func (e entry) WithPath(p string) entry {
	e.path = p
	return e
}

// entries get all entries in an archive Reader.
// Test Fatal if non EOF error returned by reader.
// returns ordered by name
func entries(t *testing.T, arch archive.Reader) (toReturn []entry) {
	for {
		fi, err := arch.NextHeader()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("error getting next header (%v) %v", len(toReturn), err)
		}

		toAdd := entry{path: fi.Path(), size: fi.Size(), mode: fi.Mode(), mtime: fi.ModTime(), symlink: fi.Symlink()}
		if toAdd.mode.IsRegular() {
			sha1 := sha1.New()
			n, err := io.Copy(sha1, arch)
			if err != nil {
				t.Fatalf("error reading archive: (%v) %v", len(toReturn), err)
			}
			if n != fi.Size() {
				t.Errorf("unexpected copy exp: %d, act: %d", fi.Size(), n)
			}
			toAdd.sha1 = hex.EncodeToString(sha1.Sum(nil))
		}
		toReturn = append(toReturn, toAdd)
	}

	// sort by name from reader
	slices.SortFunc(toReturn, func(a, b entry) int {
		return strings.Compare(a.path, b.path)
	})

	return
}

// assertEntries checks if all entries wanted are in got
func assertEntries(t *testing.T, want []entry, got []entry) {
	gotCount := len(got)
	test.AssertEqual(t, len(want), gotCount, "entries length dont match")
	for i, expRes := range want {
		if i >= gotCount {
			t.Errorf("missing result %v", expRes)
			continue
		}

		actRes := got[i]
		test.AssertEqualF(t, expRes.path, actRes.path, "path doesnt match (%d)", i)
		test.AssertEqualF(t, expRes.symlink, actRes.symlink, "symlink doesnt match (%d)", i)
		test.AssertEqualF(t, expRes.size, actRes.size, "size doesnt match (%d)", i)
		test.AssertEqualF(t, expRes.mode, actRes.mode, "mode doesnt match (%d)", i)
		test.AssertEqualF(t, expRes.mtime.UTC().UnixNano(), actRes.mtime.UTC().UnixNano(), "mtime doesnt match (%d)", i)
		test.AssertEqualF(t, expRes.sha1, actRes.sha1, "sha1 doesnt match (%d)", i)
	}
}
