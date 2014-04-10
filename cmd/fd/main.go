// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// fd is a command-line utility that opens a file for reading/writing and redirects its I/O to the standard descriptors.
//	Usage: fd /dev/file
package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/gocircuit/circuit/kit/iomisc"
)

var (
	flagTTY = flag.Bool("tty", false, "Attach a PTY/TTY pair between this console and the file")
	flagEve = flag.Bool("eve", false, "Eavesdrop on the standard descriptor")
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("usage: %s file\n", os.Args[0])
	}
	fw, err := os.OpenFile(flag.Arg(0), os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("open (%s)", err)
	}
	defer fw.Close()
	//
	// fr, err := os.OpenFile(flag.Arg(0), os.O_RDWR, 0666)
	// if err != nil {
	// 	log.Printf("open (%s)", err)
	// }
	// defer fr.Close()
	//
	// Go structure dup
	fd2, err := syscall.Dup(int(fw.Fd()))
	if err != nil {
		log.Fatalf("dup (%s)", err)
	}
	fr := os.NewFile(uintptr(fd2), "")
	defer fr.Close()
	//
	var z = newPinger()
	iomisc.SniffClose("stdout ⟽ file", os.Stdout, watchfile{fr}, z.ping)
	iomisc.SniffClose("stdin ⟾ file", watchfile{fw}, os.Stdin, z.ping)
	//
	z.pong()
	log.Println("pong")
	z.pong()
	log.Println("pong")
}

//
type watchfile struct {
	*os.File
}

func (w watchfile) Chdir() error {
	log.Printf("f⟫ Chdir")
	return w.File.Chdir()
}

func (w watchfile) Chmod(mode os.FileMode) error {
	log.Printf("f⟫ Chmod")
	return w.File.Chmod(mode)
}

func (w watchfile) Chown(uid, gid int) error {
	log.Printf("f⟫ Chown")
	return w.File.Chown(uid, gid)
}

func (w watchfile) Close() error {
	log.Printf("f⟫ Close")
	return w.File.Close()
}

func (w watchfile) Fd() uintptr {
	log.Printf("f⟫ Fd")
	return w.File.Fd()
}

func (w watchfile) Name() string {
	log.Printf("f⟫ Name")
	return w.File.Name()
}

func (w watchfile) Read(b []byte) (n int, err error) {
	log.Printf("f⟫ Read")
	return w.File.Read(b)
}

func (w watchfile) ReadAt(b []byte, off int64) (n int, err error) {
	log.Printf("f⟫ ReadAt")
	return w.File.ReadAt(b, off)
}

func (w watchfile) Readdir(n int) (fi []os.FileInfo, err error) {
	log.Printf("f⟫ Readdir")
	return w.File.Readdir(n)
}

func (w watchfile) Readdirnames(n int) (names []string, err error) {
	log.Printf("f⟫ Readdirnames")
	return w.File.Readdirnames(n)
}

func (w watchfile) Seek(offset int64, whence int) (ret int64, err error) {
	log.Printf("f⟫ Seek")
	return w.File.Seek(offset, whence)
}

func (w watchfile) Stat() (fi os.FileInfo, err error) {
	log.Printf("f⟫ Stat")
	return w.File.Stat()
}

func (w watchfile) Sync() (err error) {
	log.Printf("f⟫ Sync")
	return w.File.Sync()
}

func (w watchfile) Truncate(size int64) error {
	log.Printf("f⟫ Truncate")
	return w.File.Truncate(size)
}

func (w watchfile) Write(b []byte) (n int, err error) {
	defer w.Sync()
	log.Printf("f⟫ Write")
	return w.File.Write(b)
}

func (w watchfile) WriteAt(b []byte, off int64) (n int, err error) {
	defer w.Sync()
	log.Printf("f⟫ WriteAt")
	return w.File.WriteAt(b, off)
}

func (w watchfile) WriteString(s string) (ret int, err error) {
	defer w.Sync()
	log.Printf("f⟫ WriteString")
	return w.File.WriteString(s)
}

//
type pinger chan struct{}

func newPinger() pinger {
	ch := make(chan struct{})
	return pinger(ch)
}

func (p pinger) ping() {
	p <- struct{}{}
}

func (p pinger) pong() {
	<-p
}
