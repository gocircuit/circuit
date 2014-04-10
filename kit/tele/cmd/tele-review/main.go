// tele-review is a debugging tool for internal use only.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
}

const (
	logTimeLayout = "2006/01/02 15:04:05"
)

var (
	addrToFrame map[string]string
)

type frame struct {
	Time  time.Time
	Addr  string
	Frame string
}

func parse(r io.Reader) <-chan *frame {
	ch := make(chan *frame)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			var err error
			f := &frame{}
			orig := scanner.Text()
			l := orig
			// Time
			f.Time, err = time.Parse(logTimeLayout, l[:len(logTimeLayout)])
			if err != nil {
				// Not a log line. Skip.
				continue
			}
			l = skipWhite(l[len(logTimeLayout):])
			// Addr
			if len(l) < 2 {
				continue
			}
			if l[0] != '(' {
				continue
			}
			l = l[1:]
			j := strings.Index(l, ")")
			if j < 0 {
				continue
			}
			f.Addr = l[:j]
			l = skipWhite(l[j+1:])
			// Frame
			j = strings.IndexAny(l, " \t")
			if j < 0 {
				f.Frame = l
				fmt.Fprintf(os.Stderr, "Log line has no message:\n%s\n", orig)
			} else {
				f.Frame = l[:j]
			}
			ch <- f
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()
	return ch
}

func skipWhite(l string) string {
	for i, c := range l {
		switch c {
		case ' ', '\t':
		default:
			return l[i:]
		}
	}
	return ""
}
