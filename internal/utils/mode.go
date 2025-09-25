package utils

import "os"

const (
	S_IFMT    = 0170000 // Bitmask for file type
	S_IFDIR   = 0040000 // Directory
	S_IFLNK   = 0120000 // Symbolic link
	S_IFREG   = 0100000 // Regular file
	S_PERMMSK = 0000777 // Permission bits mask
)

func UnixToMode(unixMode int64) os.FileMode {
	// Isolate the permission bits (lowest 9 bits)
	mode := os.FileMode(unixMode & S_PERMMSK)

	// Map the Unix type bits to Go's specific high-bit definitions
	switch unixMode & S_IFMT {
	case S_IFDIR:
		mode |= os.ModeDir
	case S_IFLNK:
		mode |= os.ModeSymlink
	}

	return mode
}
