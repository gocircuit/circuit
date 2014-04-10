package symbolizer

import (
	"bufio"
	"io"
	"strings"
)

/*
_start
	goroutine profile: total 18
_header
	1 @ 0x421819e 0x4217fa7 0x4214fc2 0x41d0017 0x4189aef 0x4188b4b 0x40f2d72 0x40f9f55 0x40fda75 0x404d3a0
_stack
#   0x421819e	runtime/pprof.writeRuntimeProfile+0x9e	/Users/petar/.datahero/build/go/src/pkg/runtime/pprof/pprof.go:540
_frame

_empty
*/

type Trace struct {
	Header    string
	Goroutine []*Stack
}

type Stack struct {
	Frame []*Frame
}

type Frame struct {
	Func   string
	Source string
}

type lineType int

const (
	_start = lineType(iota)
	_header
	_stack
	_frame
	_empty
)

func Symbolize(r io.Reader) (trace *Trace) {
	trace = &Trace{}

	scanner := bufio.NewScanner(r)
	var prev lineType = _start
	for scanner.Scan() {
		text := scanner.Text()
		//println(fmt.Sprintf("---> %s", text))
		switch prev {

		case _start:
			trace.Header = text
			prev = _header

		case _stack, _frame:
			if strings.TrimSpace(text) == "" {
				prev = _empty
				break
			}
			parts := strings.Split(text, "\t")
			//fmt.Printf("#goroutine=%d ••• #parts=%d ••• %#v\n", len(trace.Goroutine), len(parts), parts)
			stack := trace.Goroutine[len(trace.Goroutine)-1]
			stack.Frame = append(stack.Frame, &Frame{Func: parts[2], Source: parts[3]})
			prev = _frame

		case _empty:
			if strings.TrimSpace(text) == "" {
				prev = _empty
				break
			}
			fallthrough

		case _header:
			trace.Goroutine = append(trace.Goroutine, &Stack{})
			prev = _stack

		default:
			panic("urgh")
		}
	}
	return
}
