package pgproto

import (
	"fmt"
	"io"
)

type Describe struct {
	ObjectType ObjectType
	Name       []byte
}

func (d *Describe) client() {}

func ParseDescribe(r io.Reader) (*Describe, error) {
	b := newReadBuffer(r)

	// 'D' [int32 - length] [byte - object type] [string - name] \0
	err := b.ReadTag('D')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	d := &Describe{}

	t, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	switch o := ObjectType(t); o {
	case ObjectTypePreparedStatement:
	case ObjectTypePortal:
		d.ObjectType = o
	default:
		return nil, fmt.Errorf("unknown describe object type %#v", t)
	}

	d.Name, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Describe) Encode() []byte {
	// 'D' [int32 - length] [byte - object type] [string - name] \0
	w := newWriteBuffer()
	w.WriteByte(byte(d.ObjectType))
	w.WriteString(d.Name, true)
	w.Wrap('D')
	return w.Bytes()
}

func (d *Describe) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "Describe",
		"Payload": map[string]interface{}{
			"ObjectType": byte(d.ObjectType),
			"Name":       string(d.Name),
		},
	}
}

func (d *Describe) String() string { return messageToString(d) }
