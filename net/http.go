package http_util

import (
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Protocols struct {
	bits uint8
}

const (
	protoHTTP1 = 1 << iota
	protoHTTP2
	protoUnencryptedHTTP2
)

func (p Protocols) HTTP1() bool                  { return p.bits&protoHTTP1 != 0 }
func (p *Protocols) SetHTTP1(ok bool)            { p.setBit(protoHTTP1, ok) }
func (p Protocols) HTTP2() bool                  { return p.bits&protoHTTP2 != 0 }
func (p *Protocols) SetHTTP2(ok bool)            { p.setBit(protoHTTP2, ok) }
func (p Protocols) UnencryptedHTTP2() bool       { return p.bits&protoUnencryptedHTTP2 != 0 }
func (p *Protocols) SetUnencryptedHTTP2(ok bool) { p.setBit(protoUnencryptedHTTP2, ok) }

func (p *Protocols) setBit(bit uint8, ok bool) {
	if ok {
		p.bits |= bit
	} else {
		p.bits &^= bit
	}
}

func (p Protocols) String() string {
	var s []string
	if p.HTTP1() {
		s = append(s, "HTTP1")
	}
	if p.HTTP2() {
		s = append(s, "HTTP2")
	}
	if p.UnencryptedHTTP2() {
		s = append(s, "UnencryptedHTTP2")
	}
	return "{" + strings.Join(s, ",") + "}"
}

type incomparable [0]func()

const maxInt64 = 1<<63 - 1

var aLongTimeAgo = time.Unix(1, 0)

var omitBundledHTTP2 bool

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "net/http context value " + k.name }

func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}

func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}

func stringContainsCTLByte(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < ' ' || b == 0x7f {
			return true
		}
	}
	return false
}

func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	var pos int
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			if pos < i {
				b = append(b, s[pos:i]...)
			}
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
			pos = i + 1
		}
	}
	if pos < len(s) {
		b = append(b, s[pos:]...)
	}
	return string(b)
}

var NoBody = noBody{}

type noBody struct{}

func (noBody) Read([]byte) (int, error)         { return 0, io.EOF }
func (noBody) Close() error                     { return nil }
func (noBody) WriteTo(io.Writer) (int64, error) { return 0, nil }

var (
	_ io.WriterTo   = NoBody
	_ io.ReadCloser = NoBody
)

type PushOptions struct {
	Method string
	Header Header
}

type Pusher interface {
	Push(target string, opts *PushOptions) error
}

type HTTP2Config struct {
	MaxConcurrentStreams          int
	MaxDecoderHeaderTableSize     int
	MaxEncoderHeaderTableSize     int
	MaxReadFrameSize              int
	MaxReceiveBufferPerConnection int
	MaxReceiveBufferPerStream     int
	SendPingTimeout               time.Duration
	PingTimeout                   time.Duration
	WriteByteTimeout              time.Duration
	PermitProhibitedCipherSuites  bool
	CountError                    func(errType string)
}
