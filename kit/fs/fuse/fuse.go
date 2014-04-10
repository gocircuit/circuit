// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

// BUG(rsc): The mount code for FreeBSD has not been written yet.

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	KernelVersionMajor = 7
	KernelVersionMinor = 8
)

// A Request represents a single FUSE request received from the kernel.
// Use a type switch to determine the specific kind.
// A request of unrecognized type will have concrete type *Header.
type Request interface {
	// Hdr returns the Header associated with this request.
	Hdr() *Header

	// RespondError responds to the request with the given error.
	RespondError(Error)

	String() string
}

// A RequestID identifies an active FUSE request.
type RequestID uint64

// A NodeID is a number identifying a directory or file.
// It must be unique among IDs returned in LookupResponses
// that have not yet been forgotten by ForgetRequests.
type NodeID uint64

// A HandleID is a number identifying an open directory or file.
// It only needs to be unique while the directory or file is open.
type HandleID uint64

// The RootID identifies the root directory of a FUSE file system.
const RootID NodeID = rootID

// A Header describes the basic information sent in every request.
type Header struct {
	Conn *Conn     // connection this request was received on
	ID   RequestID // unique ID for request
	Node NodeID    // file or directory the request is about
	Uid  uint32    // user ID of process making request
	Gid  uint32    // group ID of process making request
	Pid  uint32    // process ID of process making request
}

func (h *Header) String() string {
	return fmt.Sprintf("ID=%#x Node=%#x Uid=%d Gid=%d Pid=%d", h.ID, h.Node, h.Uid, h.Gid, h.Pid)
}

func (h *Header) Hdr() *Header {
	return h
}

func (h *Header) RespondError(err Error) {
	// FUSE uses negative errors!
	// TODO: File bug report against OSXFUSE: positive error causes kernel panic.
	out := &outHeader{Error: -err.errno(), Unique: uint64(h.ID)}
	h.Conn.respond(out, unsafe.Sizeof(*out))
}

// An InitRequest is the first request sent on a FUSE file system.
type InitRequest struct {
	Header
	Major        uint32
	Minor        uint32
	MaxReadahead uint32
	Flags        InitFlags
}

func (r *InitRequest) String() string {
	return fmt.Sprintf("Init [%s] %d.%d ra=%d fl=%v", &r.Header, r.Major, r.Minor, r.MaxReadahead, r.Flags)
}

// An InitResponse is the response to an InitRequest.
type InitResponse struct {
	MaxReadahead uint32
	Flags        InitFlags
	MaxWrite     uint32
}

func (r *InitResponse) String() string {
	return fmt.Sprintf("Init %+v", *r)
}

// Respond replies to the request with the given response.
func (r *InitRequest) Respond(resp *InitResponse) {
	out := &initOut{
		outHeader:    outHeader{Unique: uint64(r.ID)},
		Major:        kernelVersion,
		Minor:        kernelMinorVersion,
		MaxReadahead: resp.MaxReadahead,
		Flags:        uint32(resp.Flags),
		MaxWrite:     resp.MaxWrite,
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A StatfsRequest requests information about the mounted file system.
type StatfsRequest struct {
	Header
}

func (r *StatfsRequest) String() string {
	return fmt.Sprintf("Statfs [%s]\n", &r.Header)
}

// Respond replies to the request with the given response.
func (r *StatfsRequest) Respond(resp *StatfsResponse) {
	out := &statfsOut{
		outHeader: outHeader{Unique: uint64(r.ID)},
		St: kstatfs{
			Blocks:  resp.Blocks,
			Bfree:   resp.Bfree,
			Bavail:  resp.Bavail,
			Files:   resp.Files,
			Bsize:   resp.Bsize,
			Namelen: resp.Namelen,
			Frsize:  resp.Frsize,
		},
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A StatfsResponse is the response to a StatfsRequest.
type StatfsResponse struct {
	Blocks  uint64 // Total data blocks in file system.
	Bfree   uint64 // Free blocks in file system.
	Bavail  uint64 // Free blocks in file system if you're not root.
	Files   uint64 // Total files in file system.
	Ffree   uint64 // Free files in file system.
	Bsize   uint32 // Block size
	Namelen uint32 // Maximum file name length?
	Frsize  uint32 // ?
}

func (r *StatfsResponse) String() string {
	return fmt.Sprintf("Statfs %+v", *r)
}

// An AccessRequest asks whether the file can be accessed
// for the purpose specified by the mask.
type AccessRequest struct {
	Header
	Mask uint32
}

func (r *AccessRequest) String() string {
	return fmt.Sprintf("Access [%s] mask=%#x", &r.Header, r.Mask)
}

// Respond replies to the request indicating that access is allowed.
// To deny access, use RespondError.
func (r *AccessRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// An Attr is the metadata for a single file or directory.
type Attr struct {
	Inode  uint64      // inode number
	Size   uint64      // size in bytes
	Blocks uint64      // size in blocks
	Atime  time.Time   // time of last access
	Mtime  time.Time   // time of last modification
	Ctime  time.Time   // time of last inode change
	Crtime time.Time   // time of creation (OS X only)
	Mode   os.FileMode // file mode
	Nlink  uint32      // number of links
	Uid    uint32      // owner uid
	Gid    uint32      // group gid
	Rdev   uint32      // device numbers
	Flags  uint32      // chflags(2) flags (OS X only)
}

func unix(t time.Time) (sec uint64, nsec uint32) {
	nano := t.UnixNano()
	sec = uint64(nano / 1e9)
	nsec = uint32(nano % 1e9)
	return
}

func (a *Attr) attr() (out attr) {
	out.Ino = a.Inode
	out.Size = a.Size
	out.Blocks = a.Blocks
	out.Atime, out.AtimeNsec = unix(a.Atime)
	out.Mtime, out.MtimeNsec = unix(a.Mtime)
	out.Ctime, out.CtimeNsec = unix(a.Ctime)
	out.SetCrtime(unix(a.Crtime))
	out.Mode = uint32(a.Mode) & 0777
	switch {
	default:
		out.Mode |= syscall.S_IFREG
	case a.Mode&os.ModeDir != 0:
		out.Mode |= syscall.S_IFDIR
	case a.Mode&os.ModeDevice != 0:
		if a.Mode&os.ModeCharDevice != 0 {
			out.Mode |= syscall.S_IFCHR
		} else {
			out.Mode |= syscall.S_IFBLK
		}
	case a.Mode&os.ModeNamedPipe != 0:
		out.Mode |= syscall.S_IFIFO
	case a.Mode&os.ModeSymlink != 0:
		out.Mode |= syscall.S_IFLNK
	case a.Mode&os.ModeSocket != 0:
		out.Mode |= syscall.S_IFSOCK
	}
	if a.Mode&os.ModeSetuid != 0 {
		out.Mode |= syscall.S_ISUID
	}
	if a.Mode&os.ModeSetgid != 0 {
		out.Mode |= syscall.S_ISGID
	}
	out.Nlink = a.Nlink
	if out.Nlink < 1 {
		out.Nlink = 1
	}
	out.Uid = a.Uid
	out.Gid = a.Gid
	out.Rdev = a.Rdev
	out.SetFlags(a.Flags)

	return
}

// A GetattrRequest asks for the metadata for the file denoted by r.Node.
type GetattrRequest struct {
	Header
}

func (r *GetattrRequest) String() string {
	return fmt.Sprintf("Getattr [%s]", &r.Header)
}

// Respond replies to the request with the given response.
func (r *GetattrRequest) Respond(resp *GetattrResponse) {
	out := &attrOut{
		outHeader:     outHeader{Unique: uint64(r.ID)},
		AttrValid:     uint64(resp.AttrValid / time.Second),
		AttrValidNsec: uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:          resp.Attr.attr(),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A GetattrResponse is the response to a GetattrRequest.
type GetattrResponse struct {
	AttrValid time.Duration // how long Attr can be cached
	Attr      Attr          // file attributes
}

func (r *GetattrResponse) String() string {
	return fmt.Sprintf("Getattr %+v", *r)
}

// A GetxattrRequest asks for the extended attributes associated with r.Node.
type GetxattrRequest struct {
	Header
	Size     uint32 // maximum size to return
	Position uint32 // offset within extended attributes
}

func (r *GetxattrRequest) String() string {
	return fmt.Sprintf("Getxattr [%s] %d @%d", &r.Header, r.Size, r.Position)
}

// Respond replies to the request with the given response.
func (r *GetxattrRequest) Respond(resp *GetxattrResponse) {
	out := &getxattrOut{
		outHeader: outHeader{Unique: uint64(r.ID)},
		Size:      uint32(len(resp.Xattr)),
	}
	r.Conn.respondData(&out.outHeader, unsafe.Sizeof(*out), resp.Xattr)
}

// A GetxattrResponse is the response to a GetxattrRequest.
type GetxattrResponse struct {
	Xattr []byte
}

func (r *GetxattrResponse) String() string {
	return fmt.Sprintf("Getxattr %x", r.Xattr)
}

// A ListxattrRequest asks to list the extended attributes associated with r.Node.
type ListxattrRequest struct {
	Header
	Size     uint32 // maximum size to return
	Position uint32 // offset within attribute list
}

func (r *ListxattrRequest) String() string {
	return fmt.Sprintf("Listxattr [%s] %d @%d", &r.Header, r.Size, r.Position)
}

// Respond replies to the request with the given response.
func (r *ListxattrRequest) Respond(resp *ListxattrResponse) {
	out := &getxattrOut{
		outHeader: outHeader{Unique: uint64(r.ID)},
		Size:      uint32(len(resp.Xattr)),
	}
	r.Conn.respondData(&out.outHeader, unsafe.Sizeof(*out), resp.Xattr)
}

// A ListxattrResponse is the response to a ListxattrRequest.
type ListxattrResponse struct {
	Xattr []byte
}

func (r *ListxattrResponse) String() string {
	return fmt.Sprintf("Listxattr %x", r.Xattr)
}

// A RemovexattrRequest asks to remove an extended attribute associated with r.Node.
type RemovexattrRequest struct {
	Header
	Name string // name of extended attribute
}

func (r *RemovexattrRequest) String() string {
	return fmt.Sprintf("Removexattr [%s] %q", &r.Header, r.Name)
}

// Respond replies to the request, indicating that the attribute was removed.
func (r *RemovexattrRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// A SetxattrRequest asks to set an extended attribute associated with a file.
type SetxattrRequest struct {
	Header
	Flags    uint32
	Position uint32 // OS X only
	Name     string
	Xattr    []byte
}

func (r *SetxattrRequest) String() string {
	return fmt.Sprintf("Setxattr [%s] %q %x fl=%v @%#x", &r.Header, r.Name, r.Xattr, r.Flags, r.Position)
}

// Respond replies to the request, indicating that the extended attribute was set.
func (r *SetxattrRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// A LookupRequest asks to look up the given name in the directory named by r.Node.
type LookupRequest struct {
	Header
	Name string
}

func (r *LookupRequest) String() string {
	return fmt.Sprintf("Lookup [%s] %q", &r.Header, r.Name)
}

// Respond replies to the request with the given response.
func (r *LookupRequest) Respond(resp *LookupResponse) {
	out := &entryOut{
		outHeader:      outHeader{Unique: uint64(r.ID)},
		Nodeid:         uint64(resp.Node),
		Generation:     resp.Generation,
		EntryValid:     uint64(resp.EntryValid / time.Second),
		EntryValidNsec: uint32(resp.EntryValid % time.Second / time.Nanosecond),
		AttrValid:      uint64(resp.AttrValid / time.Second),
		AttrValidNsec:  uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:           resp.Attr.attr(),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A LookupResponse is the response to a LookupRequest.
type LookupResponse struct {
	Node       NodeID
	Generation uint64
	EntryValid time.Duration
	AttrValid  time.Duration
	Attr       Attr
}

func (r *LookupResponse) String() string {
	return fmt.Sprintf("Lookup %+v", *r)
}

// An OpenRequest asks to open a file or directory
type OpenRequest struct {
	Header
	Dir   bool // is this Opendir?
	Flags uint32
	Mode  os.FileMode
}

func (r *OpenRequest) String() string {
	return fmt.Sprintf("Open [%s] dir=%v fl=%v mode=%v", &r.Header, r.Dir, r.Flags, r.Mode)
}

// Respond replies to the request with the given response.
func (r *OpenRequest) Respond(resp *OpenResponse) {
	out := &openOut{
		outHeader: outHeader{Unique: uint64(r.ID)},
		Fh:        uint64(resp.Handle),
		OpenFlags: uint32(resp.Flags),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A OpenResponse is the response to a OpenRequest.
type OpenResponse struct {
	Handle HandleID
	Flags  OpenFlags
}

func (r *OpenResponse) String() string {
	return fmt.Sprintf("Open %+v", *r)
}

// A CreateRequest asks to create and open a file (not a directory).
type CreateRequest struct {
	Header
	Name  string
	Flags uint32
	Mode  os.FileMode
}

func (r *CreateRequest) String() string {
	return fmt.Sprintf("Create [%s] %q fl=%v mode=%v", &r.Header, r.Name, r.Flags, r.Mode)
}

// Respond replies to the request with the given response.
func (r *CreateRequest) Respond(resp *CreateResponse) {
	out := &createOut{
		outHeader: outHeader{Unique: uint64(r.ID)},

		Nodeid:         uint64(resp.Node),
		Generation:     resp.Generation,
		EntryValid:     uint64(resp.EntryValid / time.Second),
		EntryValidNsec: uint32(resp.EntryValid % time.Second / time.Nanosecond),
		AttrValid:      uint64(resp.AttrValid / time.Second),
		AttrValidNsec:  uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:           resp.Attr.attr(),

		Fh:        uint64(resp.Handle),
		OpenFlags: uint32(resp.Flags),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A CreateResponse is the response to a CreateRequest.
// It describes the created node and opened handle.
type CreateResponse struct {
	LookupResponse
	OpenResponse
}

func (r *CreateResponse) String() string {
	return fmt.Sprintf("Create %+v", *r)
}

// A MkdirRequest asks to create (but not open) a directory.
type MkdirRequest struct {
	Header
	Name string
	Mode os.FileMode
}

func (r *MkdirRequest) String() string {
	return fmt.Sprintf("Mkdir [%s] %q mode=%v", &r.Header, r.Name, r.Mode)
}

// Respond replies to the request with the given response.
func (r *MkdirRequest) Respond(resp *MkdirResponse) {
	out := &entryOut{
		outHeader:      outHeader{Unique: uint64(r.ID)},
		Nodeid:         uint64(resp.Node),
		Generation:     resp.Generation,
		EntryValid:     uint64(resp.EntryValid / time.Second),
		EntryValidNsec: uint32(resp.EntryValid % time.Second / time.Nanosecond),
		AttrValid:      uint64(resp.AttrValid / time.Second),
		AttrValidNsec:  uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:           resp.Attr.attr(),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A MkdirResponse is the response to a MkdirRequest.
type MkdirResponse struct {
	LookupResponse
}

func (r *MkdirResponse) String() string {
	return fmt.Sprintf("Mkdir %+v", *r)
}

// A ReadRequest asks to read from an open file.
type ReadRequest struct {
	Header
	Dir    bool // is this Readdir?
	Handle HandleID
	Offset int64
	Size   int
}

func (r *ReadRequest) String() string {
	return fmt.Sprintf("Read [%s] %#x %d @%#x dir=%v", &r.Header, r.Handle, r.Size, r.Offset, r.Dir)
}

// Respond replies to the request with the given response.
func (r *ReadRequest) Respond(resp *ReadResponse) {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respondData(out, unsafe.Sizeof(*out), resp.Data)
}

// A ReadResponse is the response to a ReadRequest.
type ReadResponse struct {
	Data []byte
}

func (r *ReadResponse) String() string {
	return fmt.Sprintf("Read #%d", len(r.Data))
}

// A ReleaseRequest asks to release (close) an open file handle.
type ReleaseRequest struct {
	Header
	Dir          bool // is this Releasedir?
	Handle       HandleID
	Flags        uint32 // flags from OpenRequest
	ReleaseFlags ReleaseFlags
	LockOwner    uint32
}

func (r *ReleaseRequest) String() string {
	return fmt.Sprintf("Release [%s] %#x fl=%v rfl=%v owner=%#x", &r.Header, r.Handle, r.Flags, r.ReleaseFlags, r.LockOwner)
}

// Respond replies to the request, indicating that the handle has been released.
func (r *ReleaseRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// A DestroyRequest is sent by the kernel when unmounting the file system.
// No more requests will be received after this one, but it should still be
// responded to.
type DestroyRequest struct {
	Header
}

func (r *DestroyRequest) String() string {
	return fmt.Sprintf("Destroy [%s]", &r.Header)
}

// Respond replies to the request.
func (r *DestroyRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// A ForgetRequest is sent by the kernel when forgetting about r.Node
// as returned by r.N lookup requests.
type ForgetRequest struct {
	Header
	N uint64
}

func (r *ForgetRequest) String() string {
	return fmt.Sprintf("Forget [%s] %d", &r.Header, r.N)
}

// Respond replies to the request, indicating that the forgetfulness has been recorded.
func (r *ForgetRequest) Respond() {
	// Don't reply to forget messages.
}

// A Dirent represents a single directory entry.
type Dirent struct {
	// Inode this entry names.
	Inode uint64

	// Type of the entry, for example DT_File.
	//
	// Setting this is optional. The zero value (DT_Unknown) means
	// callers will just need to do a Getattr when the type is
	// needed. Providing a type can speed up operations
	// significantly.
	Type DirentType

	// Name of the entry
	Name string
}

// Type of an entry in a directory listing.
type DirentType uint32

const (
	// These don't quite match os.FileMode; especially there's an
	// explicit unknown, instead of zero value meaning file. They
	// are also not quite syscall.DT_*; nothing says the FUSE
	// protocol follows those, and even if they were, we don't
	// want each fs to fiddle with syscall.

	// The shift by 12 is hardcoded in the FUSE userspace
	// low-level C library, so it's safe here.

	DT_Unknown DirentType = 0
	DT_Socket  DirentType = syscall.S_IFSOCK >> 12
	DT_Link    DirentType = syscall.S_IFLNK >> 12
	DT_File    DirentType = syscall.S_IFREG >> 12
	DT_Block   DirentType = syscall.S_IFBLK >> 12
	DT_Dir     DirentType = syscall.S_IFDIR >> 12
	DT_Char    DirentType = syscall.S_IFCHR >> 12
	DT_FIFO    DirentType = syscall.S_IFIFO >> 12
)

func (t DirentType) String() string {
	switch t {
	case DT_Unknown:
		return "unknown"
	case DT_Socket:
		return "socket"
	case DT_Link:
		return "link"
	case DT_File:
		return "file"
	case DT_Block:
		return "block"
	case DT_Dir:
		return "dir"
	case DT_Char:
		return "char"
	case DT_FIFO:
		return "fifo"
	}
	return "invalid"
}

// AppendDirent appends the encoded form of a directory entry to data
// and returns the resulting slice.
func AppendDirent(data []byte, dir Dirent) []byte {
	de := dirent{
		Ino:     dir.Inode,
		Namelen: uint32(len(dir.Name)),
		Type:    uint32(dir.Type),
	}
	de.Off = uint64(len(data) + direntSize + (len(dir.Name)+7)&^7)
	data = append(data, (*[direntSize]byte)(unsafe.Pointer(&de))[:]...)
	data = append(data, dir.Name...)
	n := direntSize + uintptr(len(dir.Name))
	if n%8 != 0 {
		var pad [8]byte
		data = append(data, pad[:8-n%8]...)
	}
	return data
}

// A WriteRequest asks to write to an open file.
type WriteRequest struct {
	Header
	Handle HandleID
	Offset int64
	Data   []byte
	Flags  WriteFlags
}

func (r *WriteRequest) String() string {
	return fmt.Sprintf("Write [%s] %#x %d @%d fl=%v", &r.Header, r.Handle, len(r.Data), r.Offset, r.Flags)
}

// Respond replies to the request with the given response.
func (r *WriteRequest) Respond(resp *WriteResponse) {
	out := &writeOut{
		outHeader: outHeader{Unique: uint64(r.ID)},
		Size:      uint32(resp.Size),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A WriteResponse replies to a write indicating how many bytes were written.
type WriteResponse struct {
	Size int
}

func (r *WriteResponse) String() string {
	return fmt.Sprintf("Write %+v", *r)
}

// A SetattrRequest asks to change one or more attributes associated with a file,
// as indicated by Valid.
type SetattrRequest struct {
	Header
	Valid  SetattrValid
	Handle HandleID
	Size   uint64
	Atime  time.Time
	Mtime  time.Time
	Mode   os.FileMode
	Uid    uint32
	Gid    uint32

	// OS X only
	Bkuptime time.Time
	Chgtime  time.Time
	Crtime   time.Time
	Flags    uint32 // see chflags(2)
}

func (r *SetattrRequest) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Setattr [%s]", &r.Header)
	if r.Valid.Mode() {
		fmt.Fprintf(&buf, " mode=%v", r.Mode)
	}
	if r.Valid.Uid() {
		fmt.Fprintf(&buf, " uid=%d", r.Uid)
	}
	if r.Valid.Gid() {
		fmt.Fprintf(&buf, " gid=%d", r.Gid)
	}
	if r.Valid.Size() {
		fmt.Fprintf(&buf, " size=%d", r.Size)
	}
	if r.Valid.Atime() {
		fmt.Fprintf(&buf, " atime=%v", r.Atime)
	}
	if r.Valid.Mtime() {
		fmt.Fprintf(&buf, " mtime=%v", r.Mtime)
	}
	if r.Valid.Handle() {
		fmt.Fprintf(&buf, " handle=%#x", r.Handle)
	} else {
		fmt.Fprintf(&buf, " handle=INVALID-%#x", r.Handle)
	}
	if r.Valid.Crtime() {
		fmt.Fprintf(&buf, " crtime=%v", r.Crtime)
	}
	if r.Valid.Chgtime() {
		fmt.Fprintf(&buf, " chgtime=%v", r.Chgtime)
	}
	if r.Valid.Bkuptime() {
		fmt.Fprintf(&buf, " bkuptime=%v", r.Bkuptime)
	}
	if r.Valid.Flags() {
		fmt.Fprintf(&buf, " flags=%#x", r.Flags)
	}
	return buf.String()
}

// Respond replies to the request with the given response,
// giving the updated attributes.
func (r *SetattrRequest) Respond(resp *SetattrResponse) {
	out := &attrOut{
		outHeader:     outHeader{Unique: uint64(r.ID)},
		AttrValid:     uint64(resp.AttrValid / time.Second),
		AttrValidNsec: uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:          resp.Attr.attr(),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A SetattrResponse is the response to a SetattrRequest.
type SetattrResponse struct {
	AttrValid time.Duration // how long Attr can be cached
	Attr      Attr          // file attributes
}

func (r *SetattrResponse) String() string {
	return fmt.Sprintf("Setattr %+v", *r)
}

// A FlushRequest asks for the current state of an open file to be flushed
// to storage, as when a file descriptor is being closed.  A single opened Handle
// may receive multiple FlushRequests over its lifetime.
type FlushRequest struct {
	Header
	Handle    HandleID
	Flags     uint32
	LockOwner uint64
}

func (r *FlushRequest) String() string {
	return fmt.Sprintf("Flush [%s] %#x fl=%#x lk=%#x", &r.Header, r.Handle, r.Flags, r.LockOwner)
}

// Respond replies to the request, indicating that the flush succeeded.
func (r *FlushRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// A RemoveRequest asks to remove a file or directory.
type RemoveRequest struct {
	Header
	Name string // name of extended attribute
	Dir  bool   // is this rmdir?
}

func (r *RemoveRequest) String() string {
	return fmt.Sprintf("Remove [%s] %q dir=%v", &r.Header, r.Name, r.Dir)
}

// Respond replies to the request, indicating that the file was removed.
func (r *RemoveRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// A SymlinkRequest is a request to create a symlink making NewName point to Target.
type SymlinkRequest struct {
	Header
	NewName, Target string
}

func (r *SymlinkRequest) String() string {
	return fmt.Sprintf("Symlink [%s] from %q to target %q", &r.Header, r.NewName, r.Target)
}

// Respond replies to the request, indicating that the symlink was created.
func (r *SymlinkRequest) Respond(resp *SymlinkResponse) {
	out := &entryOut{
		outHeader:      outHeader{Unique: uint64(r.ID)},
		Nodeid:         uint64(resp.Node),
		Generation:     resp.Generation,
		EntryValid:     uint64(resp.EntryValid / time.Second),
		EntryValidNsec: uint32(resp.EntryValid % time.Second / time.Nanosecond),
		AttrValid:      uint64(resp.AttrValid / time.Second),
		AttrValidNsec:  uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:           resp.Attr.attr(),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A SymlinkResponse is the response to a SymlinkRequest.
type SymlinkResponse struct {
	LookupResponse
}

// A ReadlinkRequest is a request to read a symlink's target.
type ReadlinkRequest struct {
	Header
}

func (r *ReadlinkRequest) String() string {
	return fmt.Sprintf("Readlink [%s]", &r.Header)
}

func (r *ReadlinkRequest) Respond(target string) {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respondData(out, unsafe.Sizeof(*out), []byte(target))
}

// A LinkRequest is a request to create a hard link.
type LinkRequest struct {
	Header
	OldNode NodeID
	NewName string
}

func (r *LinkRequest) Respond(resp *LookupResponse) {
	out := &entryOut{
		outHeader:      outHeader{Unique: uint64(r.ID)},
		Nodeid:         uint64(resp.Node),
		Generation:     resp.Generation,
		EntryValid:     uint64(resp.EntryValid / time.Second),
		EntryValidNsec: uint32(resp.EntryValid % time.Second / time.Nanosecond),
		AttrValid:      uint64(resp.AttrValid / time.Second),
		AttrValidNsec:  uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:           resp.Attr.attr(),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A RenameRequest is a request to rename a file.
type RenameRequest struct {
	Header
	NewDir           NodeID
	OldName, NewName string
}

func (r *RenameRequest) String() string {
	return fmt.Sprintf("Rename [%s] from %q to dirnode %d %q", &r.Header, r.OldName, r.NewDir, r.NewName)
}

func (r *RenameRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

type MknodRequest struct {
	Header
	Name string
	Mode os.FileMode
	Rdev uint32
}

func (r *MknodRequest) String() string {
	return fmt.Sprintf("Mknod [%s] Name %q mode %v rdev %d", &r.Header, r.Name, r.Mode, r.Rdev)
}

func (r *MknodRequest) Respond(resp *LookupResponse) {
	out := &entryOut{
		outHeader:      outHeader{Unique: uint64(r.ID)},
		Nodeid:         uint64(resp.Node),
		Generation:     resp.Generation,
		EntryValid:     uint64(resp.EntryValid / time.Second),
		EntryValidNsec: uint32(resp.EntryValid % time.Second / time.Nanosecond),
		AttrValid:      uint64(resp.AttrValid / time.Second),
		AttrValidNsec:  uint32(resp.AttrValid % time.Second / time.Nanosecond),
		Attr:           resp.Attr.attr(),
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

type FsyncRequest struct {
	Header
	Handle HandleID
	Flags  uint32
}

func (r *FsyncRequest) String() string {
	return fmt.Sprintf("Fsync [%s] Handle %v Flags %v", &r.Header, r.Handle, r.Flags)
}

func (r *FsyncRequest) Respond() {
	out := &outHeader{Unique: uint64(r.ID)}
	r.Conn.respond(out, unsafe.Sizeof(*out))
}

// An InterruptRequest is a request to interrupt another pending request. The
// reponse to that request should return an error status of EINTR.
type InterruptRequest struct {
	Header
	IntrID RequestID // ID of the request to be interrupt.
}

func (r *InterruptRequest) Respond() {
	// nothing to do here
}

func (r *InterruptRequest) String() string {
	return fmt.Sprintf("Interrupt [%s] ID %v", &r.Header, r.IntrID)
}

/*{

// A XXXRequest xxx.
type XXXRequest struct {
	Header
	xxx
}

func (r *XXXRequest) String() string {
	return fmt.Sprintf("XXX [%s] xxx", &r.Header)
}

// Respond replies to the request with the given response.
func (r *XXXRequest) Respond(resp *XXXResponse) {
	out := &xxxOut{
		outHeader: outHeader{Unique: uint64(r.ID)},
		xxx,
	}
	r.Conn.respond(&out.outHeader, unsafe.Sizeof(*out))
}

// A XXXResponse is the response to a XXXRequest.
type XXXResponse struct {
	xxx
}

func (r *XXXResponse) String() string {
	return fmt.Sprintf("XXX %+v", *r)
}

 }
*/
