package xar

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/jonathongardner/a2zar/archive"
)

const (
	typeFile      = "file"
	typeDirectory = "directory"
	typeSymlink   = "symlink"
)

type XarFileInfo struct {
	file        *xarFile
	path        string
	symlink     string
	modTime     time.Time
	mode        os.FileMode
	size        int64
	unknownType bool
	errs        []error
}

// return 1970-01-01T00:00:00Z as a bad time
var BadTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// newFileInfo creates a new XarFileInfo
// parses the file header and set the correct mode, modTime, and symlink
func newFileInfo(dir string, xf *xarFile) *XarFileInfo {
	toReturn := &XarFileInfo{
		file: xf, path: path.Join(dir, xf.Name), errs: make([]error, 0),
		mode: 0, modTime: BadTime, size: 0,
	}

	toReturn.parseTime(xf.MTime)

	// parse mode depending on the type of the file
	switch xf.Type {
	case typeDirectory:
		toReturn.parseMod(xf.Mode, os.ModeDir)
		toReturn.checkSize0(xf.Data.Size)
	case typeSymlink:
		toReturn.parseMod(xf.Mode, os.ModeSymlink)
		toReturn.symlink = xf.Link
		toReturn.checkSize0(xf.Data.Size)
	case typeFile:
		toReturn.parseMod(xf.Mode, 0)
		toReturn.size = xf.Data.Size // only set size
	default:
		toReturn.unknownType = true
		toReturn.parseMod(xf.Mode, 0)
		toReturn.size = xf.Data.Size
	}
	return toReturn
}

// parseMod parses the mode of the file and sets the mode field
// add errors to the errs field if the mode is invalid. Fixes
// the mode to be a valid file type
func (fi *XarFileInfo) parseMod(mode string, typeMode os.FileMode) {
	// default to the mode so that if there is an error the mod is at least correct
	fi.mode = typeMode
	if mode == "" {
		fi.errs = append(fi.errs, fmt.Errorf("%w: empty", ErrInvalidMode))
		return
	}
	parsedMode, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		fi.errs = append(fi.errs, fmt.Errorf("%w: %w", ErrInvalidMode, err))
		return
	}
	fi.mode = os.FileMode(parsedMode)

	typ := fi.mode & typeMode
	if typ != 0 && typ != typeMode {
		fi.errs = append(fi.errs, fmt.Errorf("%w: invalid type (exp: %s, act: %s)", ErrInvalidMode, typeMode, typ))
		return
	}
	// reset mode to correct the type
	fi.mode = (fi.mode & fs.ModePerm) | typeMode
}

// parseTime parses the modification time of the file and sets the modTime field
// add errors to the errs field if the time is invalid
func (fi *XarFileInfo) parseTime(mtime string) {
	if mtime == "" {
		fi.errs = append(fi.errs, fmt.Errorf("%w: empty", ErrInvalidTime))
		return
	}

	modTime, err := time.Parse(time.RFC3339, mtime)
	if err != nil {
		fi.errs = append(fi.errs, fmt.Errorf("%w: %w", ErrInvalidTime, err))
		return
	}
	fi.modTime = modTime
}

func (fi *XarFileInfo) checkSize0(size int64) {
	if size != 0 {
		fi.errs = append(fi.errs, fmt.Errorf("%w: expected size 0, got %d", ErrInvalidSize, size))
	}
}

// headerErrs returns ErrUnknownFileType if the file type is unknown
// or ErrInvalidHeader if there was an error paring the ToC
func (fi *XarFileInfo) headerErrs() error {
	if fi.unknownType {
		return fmt.Errorf("%w: %v", archive.ErrUnknownFileType, fi.file.Type)
	}
	if len(fi.errs) > 0 {
		return fmt.Errorf("%w: %w", ErrInvalidHeader, errors.Join(fi.errs...))
	}
	return nil
}

func (fi *XarFileInfo) Raw() xarFile {
	return *fi.file
}

// --------------FileInfo--------------
// Name returns the name of the file
func (fi *XarFileInfo) Name() string {
	return fi.file.Name
}

// Size returns the size of the file
func (fi *XarFileInfo) Size() int64 {
	return fi.size
}

// Mode returns the mode of the file
func (fi *XarFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime returns the modification time of the file
func (fi *XarFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir returns true if the file is a directory
func (fi *XarFileInfo) IsDir() bool {
	return fi.mode.IsDir()
}

// Sys returns nothing right now
func (fi *XarFileInfo) Sys() any {
	return nil
}

//--------------FileInfo--------------

// Path returns the full path of the file
func (fi *XarFileInfo) Path() string {
	return fi.path
}

// Symlink returns the target of the symlink or an empty string if the file is not a symlink
func (fi *XarFileInfo) Symlink() string {
	return fi.symlink
}
