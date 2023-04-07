package tool

import (
	"os"
	"path/filepath"
)

func GetDynamicPath(path string) string {
	test := os.Getenv("base_path")
	return filepath.Join(test, path)
}
