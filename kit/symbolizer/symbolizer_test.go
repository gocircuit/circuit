package symbolizer

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestSymbolize(t *testing.T) {
	f, err := os.Open("xxx")
	if err != nil {
		t.Fatalf("open (%s)", err)
	}
	defer f.Close()
	trace := Symbolize(f)
	fmt.Printf("%#v\n", trace)
}

const raw = `goroutine profile: total 18
1 @ 0x421819e 0x4217fa7 0x4214fc2 0x41d0017 0x4189aef 0x4188b4b 0x40f2d72 0x40f9f55 0x40fda75 0x404d3a0
#	0x421819e	runtime/pprof.writeRuntimeProfile+0x9e		/Users/petar/.datahero/build/go/src/pkg/runtime/pprof/pprof.go:540
#	0x4217fa7	runtime/pprof.writeGoroutine+0x87		/Users/petar/.datahero/build/go/src/pkg/runtime/pprof/pprof.go:502
#	0x4214fc2	runtime/pprof.(*Profile).WriteTo+0xb2		/Users/petar/.datahero/build/go/src/pkg/runtime/pprof/pprof.go:229
#	0x41d0017	circuit/sys/acid.(*Acid).RuntimeProfile+0x107	/Users/petar/.datahero/build/circuit/src/circuit/sys/acid/acid.go:53
#	0x4189aef	reflect.Value.call+0xe9f			/Users/petar/.datahero/build/go/src/pkg/reflect/value.go:474
#	0x4188b4b	reflect.Value.Call+0x9b				/Users/petar/.datahero/build/go/src/pkg/reflect/value.go:345
#	0x40f2d72	circuit/sys/lang.call+0x432			/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/call.go:50
#	0x40f9f55	circuit/sys/lang.(*Runtime).serveCall+0x495	/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/value.go:113
#	0x40fda75	circuit/sys/lang.func·009+0x435			/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/runtime.go:94

1 @ 0x404d1e4 0x403969b 0x4039c98 0x405a03d 0x405a095 0x40363aa 0x404b3c8 0x404d3a0
#	0x405a03d	circuit/load/worker.init·1+0x2d	/Users/petar/.datahero/build/circuit/src/circuit/load/worker/worker.go:26
#	0x405a095	circuit/load/worker.init+0x45	/Users/petar/.datahero/build/circuit/src/circuit/load/worker/worker.go:0
#	0x40363aa	main.init+0x4a			/Users/petar/.datahero/build/app/src/datahero/worker/init.go:0
#	0x404b3c8	runtime.main+0x88		/Users/petar/.datahero/build/go/src/pkg/runtime/proc.c:179

1 @ 0x404d3a0

1 @ 0x404d90e 0x4047b3e 0x404d3a0
#	0x404d90e	runtime.entersyscallblock+0x16e	/Users/petar/.datahero/build/go/src/pkg/runtime/proc.c:1333
#	0x4047b3e	runtime.MHeap_Scavenger+0xee	/Users/petar/.datahero/build/go/src/pkg/runtime/mheap.c:454

1 @ 0x4059320 0x4038c58 0x41de37f 0x41dccdc 0x404d3a0
#	0x4059320	return+0x0						/Users/petar/.datahero/build/go/src/pkg/runtime/asm_amd64.s:508
#	0x4038c58	runtime.cgocall+0x128					/Users/petar/.datahero/build/go/src/pkg/runtime/cgocall.c:162
#	0x41de37f	circuit/kit/zookeeper._Cfunc_wait_for_watch+0x2f	circuit/kit/zookeeper/_obj/_cgo_defun.c:475
#	0x41dccdc	circuit/kit/zookeeper._watchLoop+0x1c			circuit/kit/zookeeper/_obj/_cgo_gotypes.go:1337

1 @ 0x404d90e 0x4057887 0x41cf1fc 0x404d3a0
#	0x41cf1fc	os/signal.loop+0x1c	/Users/petar/.datahero/build/go/src/pkg/os/signal/signal_unix.go:21

1 @ 0x404d1e4 0x4039844 0x4039ce5 0x40e4bd5 0x404d3a0
#	0x40e4bd5	circuit/kit/debug.func·003+0xb5	/Users/petar/.datahero/build/circuit/src/circuit/kit/debug/sigpanic.go:58

1 @ 0x404d1e4 0x4039844 0x4039ce5 0x40e4ab5 0x404d3a0
#	0x40e4ab5	circuit/kit/debug.func·002+0xb5	/Users/petar/.datahero/build/circuit/src/circuit/kit/debug/sigpanic.go:45

3 @ 0x404d1e4 0x4039844 0x4039ce5 0x41dfaa2 0x404d3a0
#	0x41dfaa2	circuit/kit/zookeeper/zutil.func·001+0x42	/Users/petar/.datahero/build/circuit/src/circuit/kit/zookeeper/zutil/util.go:52

1 @ 0x404d1e4 0x4056def 0x40568b2 0x414a881 0x414d6e1 0x415e995 0x4101a79 0x404d3a0
#	0x4056def	netpollblock+0x9f				/Users/petar/.datahero/build/go/src/pkg/runtime/znetpoll_darwin_amd64.c:255
#	0x40568b2	net.runtime_pollWait+0x82			/Users/petar/.datahero/build/go/src/pkg/runtime/znetpoll_darwin_amd64.c:118
#	0x414a881	net.(*pollDesc).WaitRead+0x31			/Users/petar/.datahero/build/go/src/pkg/net/fd_poll_runtime.go:75
#	0x414d6e1	net.(*netFD).accept+0x2c1			/Users/petar/.datahero/build/go/src/pkg/net/fd_unix.go:385
#	0x415e995	net.(*TCPListener).AcceptTCP+0x45		/Users/petar/.datahero/build/go/src/pkg/net/tcpsock_posix.go:229
#	0x4101a79	circuit/sys/transport.(*Transport).loop+0x29	/Users/petar/.datahero/build/circuit/src/circuit/sys/transport/transport.go:118

1 @ 0x404d1e4 0x4039b76 0x4039c98 0x41019f2 0x40f8ac1 0x40fd63c 0x404d3a0
#	0x41019f2	circuit/sys/transport.(*Transport).Accept+0x32	/Users/petar/.datahero/build/circuit/src/circuit/sys/transport/transport.go:113
#	0x40f8ac1	circuit/sys/lang.(*Runtime).accept+0x61		/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/runtime.go:73
#	0x40fd63c	circuit/sys/lang.func·008+0x6c			/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/runtime.go:52

2 @ 0x404d1e4 0x4056def 0x40568b2 0x414a881 0x414ba53 0x41590a3 0x411d52c 0x411d930 0x413fe97 0x413ffef 0x41b34da 0x41ba19b 0x41ba697 0x41baab1 0x41ba983 0x4100e71 0x410035b 0x404d3a0
#	0x4056def	netpollblock+0x9f				/Users/petar/.datahero/build/go/src/pkg/runtime/znetpoll_darwin_amd64.c:255
#	0x40568b2	net.runtime_pollWait+0x82			/Users/petar/.datahero/build/go/src/pkg/runtime/znetpoll_darwin_amd64.c:118
#	0x414a881	net.(*pollDesc).WaitRead+0x31			/Users/petar/.datahero/build/go/src/pkg/net/fd_poll_runtime.go:75
#	0x414ba53	net.(*netFD).Read+0x2b3				/Users/petar/.datahero/build/go/src/pkg/net/fd_unix.go:195
#	0x41590a3	net.(*conn).Read+0xc3				/Users/petar/.datahero/build/go/src/pkg/net/net.go:123
#	0x411d52c	bufio.(*Reader).fill+0x10c			/Users/petar/.datahero/build/go/src/pkg/bufio/bufio.go:79
#	0x411d930	bufio.(*Reader).Read+0x1b0			/Users/petar/.datahero/build/go/src/pkg/bufio/bufio.go:147
#	0x413fe97	io.ReadAtLeast+0xf7				/Users/petar/.datahero/build/go/src/pkg/io/io.go:284
#	0x413ffef	io.ReadFull+0x6f				/Users/petar/.datahero/build/go/src/pkg/io/io.go:302
#	0x41b34da	encoding/gob.decodeUintReader+0xaa		/Users/petar/.datahero/build/go/src/pkg/encoding/gob/decode.go:65
#	0x41ba19b	encoding/gob.(*Decoder).recvMessage+0x4b	/Users/petar/.datahero/build/go/src/pkg/encoding/gob/decoder.go:73
#	0x41ba697	encoding/gob.(*Decoder).decodeTypeSequence+0x47	/Users/petar/.datahero/build/go/src/pkg/encoding/gob/decoder.go:159
#	0x41baab1	encoding/gob.(*Decoder).DecodeValue+0x101	/Users/petar/.datahero/build/go/src/pkg/encoding/gob/decoder.go:223
#	0x41ba983	encoding/gob.(*Decoder).Decode+0x1c3		/Users/petar/.datahero/build/go/src/pkg/encoding/gob/decoder.go:202
#	0x4100e71	circuit/sys/transport.(*swapConn).Read+0xf1	/Users/petar/.datahero/build/circuit/src/circuit/sys/transport/swap.go:83
#	0x410035b	circuit/sys/transport.(*link).readLoop+0x7b	/Users/petar/.datahero/build/circuit/src/circuit/sys/transport/link.go:126

1 @ 0x404d1e4 0x4057426 0x405762e 0x40b9bc2 0x40f54a2 0x40fd88e 0x404d3a0
#	0x4057426	semacquireimpl+0x116				/Users/petar/.datahero/build/go/src/pkg/runtime/zsema_darwin_amd64.c:113
#	0x405762e	sync.runtime_Semacquire+0x2e			/Users/petar/.datahero/build/go/src/pkg/runtime/zsema_darwin_amd64.c:165
#	0x40b9bc2	sync.(*WaitGroup).Wait+0xf2			/Users/petar/.datahero/build/go/src/pkg/sync/waitgroup.go:109
#	0x40f54a2	circuit/sys/lang.(*Runtime).serveGo+0x6d2	/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/func.go:87
#	0x40fd88e	circuit/sys/lang.func·009+0x24e			/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/runtime.go:90

1 @ 0x404d1e4 0x403969b 0x4039c98 0x40f8b4d 0x40cdc7f 0x40fd284 0x404d3a0
#	0x40f8b4d	circuit/sys/lang.(*Runtime).Hang+0x2d	/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/runtime.go:109
#	0x40cdc7f	circuit/use/circuit.Hang+0x2f		/Users/petar/.datahero/build/circuit/src/circuit/use/circuit/bind.go:148
#	0x40fd284	circuit/sys/lang.func·002+0x54		/Users/petar/.datahero/build/circuit/src/circuit/sys/lang/func.go:38

1 @ 0x404d1e4 0x404616d 0x404d3a0
#	0x404616d	runfinq+0x6d	/Users/petar/.datahero/build/go/src/pkg/runtime/mgc0.c:2182

`

func TestSimplify(t *testing.T) {
	trace := Simplify(Symbolize(bytes.NewBufferString(raw)), GoFrame)
	fmt.Printf("%#v\n", trace)
}
