package twitter

import (
	"fmt"
	"testing"
)

const filename = "/Users/petar/0/src/github.com/gocircuit/circuit/_xdata/twitter_raw_10000.txt"

func TestSweep(t *testing.T) {
	s, err := NewFileScanner(filename)
	if err != nil {
		t.Fatalf("open (%s)", err)
	}
	var n int
	for s.Scan() && n < 10 {
		r := s.Record()
		if r == nil {
			println("corrupt record")
			continue
		}
		fmt.Printf("R=%s\n", r.Text)
		n++
	}
}