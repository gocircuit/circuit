// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package hmac implements carrier transports over TCP using HMAC authentication and RC4 symmetric encryption.
package hmac

import (
	"bufio"
	"crypto/hmac"
	"crypto/rc4"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"

	//"runtime/debug"
	"strings"

	"github.com/hoijui/circuit/pkg/kit/tele/codec"
	"github.com/hoijui/circuit/pkg/kit/tele/trace"
)

func NewTransport(key []byte) codec.CarrierTransport {
	return &codecTransport{
		Frame: trace.NewFrame("hmac"),
		key:   key,
	}
}

type codecTransport struct {
	trace.Frame
	key []byte
}

func (ct *codecTransport) Listen(addr net.Addr) codec.CarrierListener {
	t := addr.String()
	if strings.Index(t, ":") < 0 {
		t = t + ":0"
	}
	l, err := net.Listen("tcp", t)
	if err != nil {
		return nil
	}
	return &codecListener{ct.key, l}
}

func (ct *codecTransport) Dial(addr net.Addr) (codec.CarrierConn, error) {
	c, err := net.Dial("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	return newCodecConn(trace.NewFrame("hmac", "dial"), c.(*net.TCPConn), ct.key)
}

type codecListener struct {
	key []byte
	net.Listener
}

func (l *codecListener) Addr() net.Addr {
	return l.Listener.Addr()
}

func (l *codecListener) Accept() codec.CarrierConn {
	for {
		c, err := l.Listener.Accept()
		if err != nil {
			log.Printf("error accepting tcp connection: %v", err)
			return nil
		}
		cc, err := newCodecConn(trace.NewFrame("hmac", "acpt"), c.(*net.TCPConn), l.key)
		if err != nil {
			continue
		}
		return cc
	}
}

type codecConn struct {
	trace.Frame
	tcp *net.TCPConn
	key []byte // shared private key for authentication
	r   *rc4Reader
	w   *rc4Writer
}

func newCodecConn(f trace.Frame, tcp *net.TCPConn, key []byte) (*codecConn, error) {
	if err := tcp.SetKeepAlive(true); err != nil {
		panic(err)
	}
	c := &codecConn{
		Frame: f,
		tcp:   tcp,
		key:   key,
	}
	if err := c.auth(); err != nil {
		return nil, err
	}
	return c, nil
}

// Pick random half-key Ying for RC4 symmetric encryption.
// Send HMAC_key(Ying), RC4_key(Ying)
// Receive H, R
// Decode remote half-seed Yang = RC4INV_key(R)
// Verify that HMAC_key(Yang) = H
// Compute send half-key as KSEND = (Ying, Yang)
// Computer receive half-key as KRECV = (Yang, Ying)
// Initialize RC4 send and receive coders with keys KSEND and KRECV, respectively
//
func (c *codecConn) auth() error {
	// Prepare our half of a random pad, the ying
	ying := pickHalfKey()
	// Sign the plaintext ying with HMAC
	yingmac := hmac.New(sha512.New, c.key)
	yingmac.Write(ying)
	// Encrypt the ying random pad, using RC4 and the shared private key.
	authcipher, err := rc4.NewCipher(c.key)
	if err != nil {
		panic(err)
	}
	yingcipher := make([]byte, len(ying))
	authcipher.XORKeyStream(yingcipher, ying)
	// Send our authentication message
	p := &authMsg{
		Sign: yingmac.Sum(nil),
		Yang: yingcipher,
	}
	if err = p.Write(c.tcp); err != nil {
		return err
	}
	// Prepare reader
	br := bufio.NewReader(c.tcp)
	// Receive the reciprocal authentication message
	q := &authMsg{}
	if err = q.Read(br); err != nil {
		return err
	}
	// Decipher yang
	authcipher, err = rc4.NewCipher(c.key)
	if err != nil {
		panic(err)
	}
	yang := make([]byte, len(q.Yang))
	authcipher.XORKeyStream(yang, q.Yang)
	// Verify MAC
	yangmac := hmac.New(sha512.New, c.key)
	yangmac.Write(yang)
	if !hmac.Equal(q.Sign, yangmac.Sum(nil)) {
		return errors.New("authentication error")
	}
	// Create encryption streams
	c.r = newRC4Reader(br, append(ying, yang...))
	c.w = newRC4Writer(c.tcp, append(yang, ying...))
	return nil
}

type authMsg struct {
	Sign []byte `json:"hmac"`
	Yang []byte `json:"pad"`
}

func (m *authMsg) String() string {
	return fmt.Sprintf("hmac=%s\npad=%s",
		base64.StdEncoding.EncodeToString(m.Sign),
		base64.StdEncoding.EncodeToString(m.Yang),
	)
}

func marshal(v interface{}) []byte {
	r, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return r
}

func (m *authMsg) Write(w io.Writer) (err error) {
	_, err = w.Write(marshal(m))
	return
}

func (m *authMsg) Read(r *bufio.Reader) (err error) {
	return json.NewDecoder(r).Decode(m)
}

func pickHalfKey() []byte {
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(rand.Int31())
	}
	key := sha512.Sum512(seed)
	return key[:]
}

func (c *codecConn) RemoteAddr() net.Addr {
	return c.tcp.RemoteAddr()
}

func (c *codecConn) Read() (chunk []byte, err error) {
	k, err := binary.ReadUvarint(c.r)
	if err != nil {
		return nil, err
	}
	var q = make([]byte, k)
	var n, m int
	for m < len(q) && err == nil {
		n, err = c.r.Read(q[m:])
		m += n
	}
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (c *codecConn) Write(chunk []byte) (err error) {
	q := make([]byte, len(chunk)+8)
	n := binary.PutUvarint(q, uint64(len(chunk)))
	m := copy(q[n:], chunk)
	_, err = c.w.Write(q[:n+m])
	return err
}

func (c *codecConn) Close() (err error) {
	return c.tcp.Close()
}

type rc4Writer struct {
	cipher *rc4.Cipher
	w      io.Writer
}

func newRC4Writer(w io.Writer, key []byte) *rc4Writer {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &rc4Writer{
		cipher: cipher,
		w:      w,
	}
}

func (w *rc4Writer) Write(p []byte) (n int, err error) {
	w.cipher.XORKeyStream(p, p)
	return w.w.Write(p)
}

type rc4Reader struct {
	cipher *rc4.Cipher
	r      *bufio.Reader
}

func newRC4Reader(r *bufio.Reader, key []byte) *rc4Reader {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &rc4Reader{
		cipher: cipher,
		r:      r,
	}
}

func (r *rc4Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.cipher.XORKeyStream(p[:n], p[:n])
	return
}

func (r *rc4Reader) ReadByte() (c byte, err error) {
	if c, err = r.r.ReadByte(); err != nil {
		return
	}
	s := []byte{c}
	r.cipher.XORKeyStream(s, s)
	return s[0], nil
}
