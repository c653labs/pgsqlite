package pgproto

import (
	"io"
)

// CloseComplete represents a server response message
type CloseComplete struct{}

func (c *CloseComplete) server() {}

// ParseCloseComplete will attempt to read a CloseComplete message from the io.Reader
func ParseCloseComplete(r io.Reader) (*CloseComplete, error) {
	buf := newReadBuffer(r)

	// '3' [int32 - length]
	err := buf.ReadTag('3')
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	return &CloseComplete{}, nil
}

// Encode will return the byte representation of this message
func (c *CloseComplete) Encode() []byte {
	// '3' [int32 - length]
	buf := newWriteBuffer()
	buf.Wrap('3')
	return buf.Bytes()
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "CloseComplete",
//     "Payload": nil,
//   }
func (c *CloseComplete) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "CloseComplete",
		"Payload": nil,
	}
}

func (c *CloseComplete) String() string { return messageToString(c) }
