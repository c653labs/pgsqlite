package pgproto

import (
	"fmt"
	"io"
)

type CopyData struct {
	Data []byte
}

func (c *CopyData) server() {}

func ParseCopyData(r io.Reader) (*CopyData, error) {
	b := newReadBuffer(r)

	// 'd' [int32 - length] [bytes - data]
	err := b.ReadTag('d')
	if err != nil {
		return nil, err
	}

	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}

	c := &CopyData{
		Data: make([]byte, l),
	}

	n, err := b.Read(c.Data)
	if n != l && err == nil {
		return nil, fmt.Errorf("expected to read length %d, instead read %d", l, n)
	}
	return c, err
}

// Encode will return the byte representation of this message
func (c *CopyData) Encode() []byte {
	// 'd' [int32 - length] [bytes - data]
	w := newWriteBuffer()
	w.WriteBytes(c.Data)
	w.Wrap('d')
	return w.Bytes()
}

func (c *CopyData) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "CopyData",
		"Payload": map[string]interface{}{
			"Data": c.Data,
		},
	}
}

func (c *CopyData) String() string { return messageToString(c) }
