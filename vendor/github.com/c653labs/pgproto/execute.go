package pgproto

import (
	"io"
)

type Execute struct {
	Portal  []byte
	MaxRows int
}

func (e *Execute) client() {}

func ParseExecute(r io.Reader) (*Execute, error) {
	b := newReadBuffer(r)

	// 'E' [int32 - length] [string - portal] \0 [int32 - max rows]
	err := b.ReadTag('E')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	e := &Execute{}

	e.Portal, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	e.MaxRows, err = buf.ReadInt()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Execute) Encode() []byte {
	// 'E' [int32 - length] [string - portal] \0 [int32 - max rows]
	w := newWriteBuffer()
	w.WriteString(e.Portal, true)
	w.WriteInt(e.MaxRows)
	w.Wrap('E')
	return w.Bytes()
}

func (e *Execute) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "Execute",
		"Payload": map[string]interface{}{
			"Portal":  string(e.Portal),
			"MaxRows": e.MaxRows,
		},
	}
}

func (e *Execute) String() string { return messageToString(e) }
