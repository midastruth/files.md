//go:build wasm

package fs

import (
	"os"
)

var Ctime = func(fi os.FileInfo) int64 {
	return fi.ModTime().Unix()
}
