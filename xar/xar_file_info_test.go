package xar

import (
	"os"
	"testing"
	"time"

	"github.com/jonathongardner/a2zar/internal/test"
)

type expXar struct {
	name    string
	path    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	symlink string
	errs    []error
}

type testRes struct {
	xar        *xarFile
	passedPath string
	exp        expXar
}

func TestFileInfo(t *testing.T) {
	toTest := []testRes{
		{
			xar: &xarFile{
				Name:  "file",
				Mode:  "0644",
				MTime: "2025-07-14T09:26:50Z",
				Type:  typeFile,
				Link:  "shouldnt/use",
				Data: xarFileData{
					Size: 12345,
				},
			},
			passedPath: "some/dir/",
			exp: expXar{
				name:    "file",
				path:    "some/dir/file",
				mode:    0644,
				size:    12345,
				modTime: time.Date(2025, 7, 14, 9, 26, 50, 0, time.UTC),
			},
		},
		{
			xar: &xarFile{
				Name:  "dir",
				Mode:  "0611",
				MTime: "2025-07-14T09:26:50Z",
				Type:  typeDirectory,
				Link:  "shouldnt/use",
				Data: xarFileData{
					Size: 12345, // shouldnt use
				},
			},
			passedPath: "/cool/dir",
			exp: expXar{
				name:    "dir",
				path:    "/cool/dir/dir",
				size:    0,
				mode:    0611 | os.ModeDir,
				modTime: time.Date(2025, 7, 14, 9, 26, 50, 0, time.UTC),
				errs:    []error{ErrInvalidSize},
			},
		},
		{
			xar: &xarFile{
				Name:  "mode-already-dir",
				Mode:  "20000000644",
				MTime: "2025-07-13T09:26:50Z",
				Type:  typeDirectory,
				Link:  "shouldnt/use",
				Data:  xarFileData{Size: 0},
			},
			exp: expXar{
				name:    "mode-already-dir",
				path:    "mode-already-dir",
				size:    0,
				mode:    0644 | os.ModeDir,
				modTime: time.Date(2025, 7, 13, 9, 26, 50, 0, time.UTC),
			},
		},
		{
			xar: &xarFile{
				Name:  "sym",
				Mode:  "0344",
				MTime: "2025-07-14T09:26:50Z",
				Type:  typeSymlink,
				Link:  "should/use",
				Data: xarFileData{
					Size: 12345, // shouldnt use
				},
			},
			exp: expXar{
				name:    "sym",
				path:    "sym",
				size:    0,
				mode:    0344 | os.ModeSymlink,
				modTime: time.Date(2025, 7, 14, 9, 26, 50, 0, time.UTC),
				symlink: "should/use",
				errs:    []error{ErrInvalidSize},
			},
		},
		{
			xar: &xarFile{
				Name:  "fix-mode",
				Mode:  "6644",
				MTime: "2025-07-14T09:26:50Z",
				Type:  typeFile,
				Link:  "shouldnt/use",
				Data:  xarFileData{Size: 1},
			},
			exp: expXar{
				name:    "fix-mode",
				path:    "fix-mode",
				size:    1,
				mode:    0644,
				modTime: time.Date(2025, 7, 14, 9, 26, 50, 0, time.UTC),
			},
		},
		{
			xar: &xarFile{
				Name:  "bad-time",
				Mode:  "0644",
				MTime: "abcd",
				Type:  typeFile,
				Link:  "shouldnt/use",
				Data:  xarFileData{Size: 2},
			},
			exp: expXar{
				name:    "bad-time",
				path:    "bad-time",
				size:    2,
				mode:    0644,
				modTime: BadTime,
				errs:    []error{ErrInvalidTime},
			},
		},
		{
			xar: &xarFile{
				Name:  "bad-mode",
				Mode:  "abcd",
				MTime: "2025-07-14T09:26:50Z",
				Type:  typeFile,
				Link:  "shouldnt/use",
				Data:  xarFileData{Size: 3},
			},
			exp: expXar{
				name:    "bad-mode",
				path:    "bad-mode",
				size:    3,
				mode:    0,
				modTime: time.Date(2025, 7, 14, 9, 26, 50, 0, time.UTC),
				errs:    []error{ErrInvalidMode},
			},
		},
	}

	for _, tt := range toTest {
		t.Run(tt.exp.name, func(t *testing.T) {
			nfi := newFileInfo(tt.passedPath, tt.xar)
			assertFileInfo(t, nfi, tt.exp, "", nfi.headerErrs())
		})
	}
}

func assertFileInfo(t *testing.T, fi *XarFileInfo, exp expXar, msg string, err error) {
	test.AssertEqual(t, exp.name, fi.Name(), msg+"FileInfo should have correct name")
	test.AssertEqual(t, exp.path, fi.Path(), msg+"FileInfo should have correct path")
	test.AssertEqual(t, exp.size, fi.Size(), msg+"FileInfo should have correct size")
	test.AssertEqual(t, exp.mode, fi.Mode(), msg+"FileInfo should have correct mode")
	test.AssertEqual(t, exp.modTime, fi.ModTime(), msg+"FileInfo should have correct modTime")
	test.AssertEqual(t, exp.symlink, fi.Symlink(), msg+"FileInfo should have correct symlink")
	test.AssertErrors(t, exp.errs, err, msg+"FileInfo should have correct error")
}
