package test

import "path/filepath"

func LargeFile(path string) string {
	return filepath.Join("..", "testdata", "lfs", path)
}
