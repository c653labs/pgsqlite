package pgproto

import (
	"io"
)

// Close represents a client request message
type Close struct {
	ObjectType ObjectType
	Name       []byte
}

func (c *Close) client() {}

// ParseClose will attempt to read a Close message from the io.Reader
func ParseClose(r io.Reader) (*Close, error) {
	b := newReadBuffer(r)

	err := b.ReadTag('C')
	if err != nil {
		return nil, err
	}

	c := &Close{}
	t, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	c.ObjectType = ObjectType(t)
	c.Name, err = b.ReadString(stripNull)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Encode will return the byte representation of this message
func (c *Close) Encode() []byte {
	b := newWriteBuffer()
	b.WriteByte(byte(c.ObjectType))
	b.WriteString(c.Name, writeNull)
	b.Wrap('C')
	return b.Bytes()
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "Close",
//     "Payload": map[string]interface{}{
//       "ObjectType": <Close.ObjectType>,
//       "Name": <Close.Name>,
//     },
//   }
func (c *Close) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "Close",
		"Payload": map[string]interface{}{
			"ObjectType": c.ObjectType,
			"Name":       c.Name,
		},
	}
}

func (c *Close) String() string { return messageToString(c) }
