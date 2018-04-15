package pgproto

import (
	"fmt"
	"io"
)

type Flush struct{}

func (f *Flush) client() {}

func ParseFlush(r io.Reader) (*Flush, error) {
	b := newReadBuffer(r)

	// 'H' [int32 - length]
	err := b.ReadTag('H')
	if err != nil {
		return nil, err
	}

	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}

	if l != 4 {
		return nil, fmt.Errorf("expected message length of 4")
	}

	return &Flush{}, nil
}

func (f *Flush) Encode() []byte {
	b := newWriteBuffer()
	b.Wrap('H')
	return b.Bytes()
}

func (f *Flush) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "Flush",
		"Payload": nil,
	}
}

func (f *Flush) String() string { return messageToString(f) }
