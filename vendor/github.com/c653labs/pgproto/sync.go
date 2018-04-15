package pgproto

import (
	"fmt"
	"io"
)

type Sync struct{}

func (s *Sync) client() {}

func ParseSync(r io.Reader) (*Sync, error) {
	b := newReadBuffer(r)

	err := b.ReadTag('S')
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

	return &Sync{}, nil
}

func (s *Sync) Encode() []byte {
	b := newWriteBuffer()
	b.Wrap('S')
	return b.Bytes()
}

func (s *Sync) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "Sync",
		"Payload": nil,
	}
}

func (s *Sync) String() string { return messageToString(s) }
