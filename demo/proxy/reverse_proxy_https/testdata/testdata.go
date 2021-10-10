package testdata

import (
	"path/filepath"
	"runtime"
)

var basepath string

func init() {
	_, currentPath, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentPath)
}

func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(basepath, rel)
}
