package symbolizer

import (
	"strings"
)

type FrameFunc func(*Frame) bool

func GoFrame(f *Frame) bool {
	i := strings.Index(f.Func, "/")
	if i < 0 {
		return true
	}
	switch f.Func[:i] {
	case
		"archive", "bufio", "builtin", "bytes", "compress", "container",
		"crypto", "database", "debug", "encoding", "errors", "expvar", "flag",
		"fmt", "go", "hash", "html", "image", "index", "io", "log", "math",
		"mime", "net", "os", "path", "reflect", "regexp", "runtime", "sort",
		"strconv", "strings", "sync", "syscall", "testing", "text", "time",
		"unicode", "unsafe":
		return true
	}
	return false
}

func Simplify(t *Trace, ff FrameFunc) *Trace {
	for _, stk := range t.Goroutine {
		var topn int
		for i, f := range stk.Frame {
			if ff(f) {
				topn = i
				break
			}
		}
		stk.Frame = stk.Frame[:topn]
	}
	return t
}
