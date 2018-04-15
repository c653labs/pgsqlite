package pgproto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type nullStrip bool

const (
	stripNull     nullStrip = true
	dontStripNull           = false
)

type nullWrite bool

const (
	writeNull     nullWrite = true
	dontWriteNull           = false
)

type readBuffer struct {
	io.Reader
}

func newReadBuffer(r io.Reader) *readBuffer {
	// If we already have a read buffer, don't create a new one
	if buf, ok := r.(*readBuffer); ok {
		return buf
	}

	buf := &readBuffer{
		Reader: r,
	}
	return buf
}

func (b *readBuffer) ReadInt() (int, error) {
	buf := make([]byte, 4)
	n, err := b.Read(buf)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, io.EOF
	}

	return bytesToInt(buf), nil
}

func (b *readBuffer) ReadInt16() (int, error) {
	buf := make([]byte, 2)
	n, err := b.Read(buf)
	if err != nil {
		return 0, err
	}
	if n != 2 {
		return 0, io.EOF
	}

	return bytesToInt16(buf), nil
}

func (b *readBuffer) ReadLength() (*readBuffer, error) {
	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}
	if l <= 0 {
		return nil, fmt.Errorf("unable to parse length from message")
	}

	// Length needs to account for the 4 bytes of the length value that have already been parsed
	l = l - 4
	if l == 0 {
		return nil, nil
	}

	buf := make([]byte, l)
	n, err := b.Read(buf)
	if err != nil {
		return nil, err
	}

	if l != n {
		return nil, fmt.Errorf("could not parse required bytes from message")
	}

	return newReadBuffer(bytes.NewReader(buf)), nil
}

func (b *readBuffer) ReadByte() (c byte, err error) {
	buf := make([]byte, 1)
	l, err := b.Read(buf)
	if l == 1 {
		c = buf[0]
	}
	return
}

func (b *readBuffer) ReadUntil(c byte) ([]byte, error) {
	buf := make([]byte, 0)
	for {
		n, err := b.ReadByte()
		if err == io.EOF {
			return buf, err
		} else if err != nil {
			return nil, err
		}

		buf = append(buf, n)
		if n == c {
			break
		}
	}

	return buf, nil
}

func (b *readBuffer) ReadString(trimNull nullStrip) ([]byte, error) {
	str, err := b.ReadUntil('\x00')
	if err != nil && err != io.EOF {
		return nil, err
	}

	if trimNull {
		str = bytes.TrimRight(str, "\x00")
	}

	return str, nil
}

func (b *readBuffer) ReadTag(t byte) error {
	tag, err := b.ReadByte()
	if err != nil {
		return err
	}
	if tag != t {
		return fmt.Errorf("invalid tag '%c' for message, must be '%c'", tag, t)
	}
	return nil
}

type writeBuffer struct {
	bytes []byte
}

func newWriteBuffer() *writeBuffer {
	return &writeBuffer{}
}

func (b *writeBuffer) Bytes() []byte {
	return b.bytes
}

func (b *writeBuffer) WriteInt(i int) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	b.bytes = append(b.bytes, buf...)
}

func (b *writeBuffer) WriteInt16(i int) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(i))
	b.bytes = append(b.bytes, buf...)
}

func (b *writeBuffer) WriteBytes(buf []byte) {
	b.bytes = append(b.bytes, buf...)
}

func (b *writeBuffer) WriteByte(c byte) error {
	b.bytes = append(b.bytes, c)
	return nil
}

func (b *writeBuffer) WriteString(buf []byte, null nullWrite) {
	b.bytes = append(b.bytes, buf...)
	if null {
		b.bytes = append(b.bytes, '\x00')
	}
}

func (b *writeBuffer) PrependLength() {
	// Need to include the 4 bytes as part of the length
	l := len(b.bytes) + 4
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(l))
	b.bytes = append(buf, b.bytes...)
}

func (b *writeBuffer) PrependByte(c byte) {
	b.bytes = append([]byte{c}, b.bytes...)
}

func (b *writeBuffer) Wrap(t byte) {
	b.PrependLength()
	b.PrependByte(t)
}

func (b *writeBuffer) Reader() *readBuffer {
	return newReadBuffer(bytes.NewReader(b.Bytes()))
}
