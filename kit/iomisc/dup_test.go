package iomisc

import (
	"io"
	"os"
	"testing"
)

var lines = []string{"hello ", "world\n", "foo", "bar\n1+", "2=3"}

func TestDup(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		for _, line := range lines {
			w.Write([]byte(line))
		}
		w.Close()
	}()
	r1, r2 := Dup(r)
	l1, l2 := PrefixReader("1:", r1), PrefixReader("2:", r2)
	io.Copy(os.Stderr, l1)
	io.Copy(os.Stderr, l2)
}
