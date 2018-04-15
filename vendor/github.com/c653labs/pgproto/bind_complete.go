package pgproto

import (
	"io"
)

// BindComplete represents a server response message
type BindComplete struct{}

func (b *BindComplete) server() {}

// ParseBindComplete will attempt to read an BindComplete message from the io.Reader
func ParseBindComplete(r io.Reader) (*BindComplete, error) {
	buf := newReadBuffer(r)

	// '2' [int32 - length]
	err := buf.ReadTag('2')
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	return &BindComplete{}, nil
}

// Encode will return the byte representation of this message
func (b *BindComplete) Encode() []byte {
	// '2' [int32 - length]
	buf := newWriteBuffer()
	buf.Wrap('2')
	return buf.Bytes()
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "BindComplete",
//     "Payload": nil,
//     },
//   }
func (b *BindComplete) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "BindComplete",
		"Payload": nil,
	}
}

func (b *BindComplete) String() string { return messageToString(b) }
