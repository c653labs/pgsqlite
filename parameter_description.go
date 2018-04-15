package pgproto

import (
	"io"
)

type ParameterDescription struct {
	OIDs []int
}

func (p *ParameterDescription) client() {}

func ParseParameterDescription(r io.Reader) (*ParameterDescription, error) {
	b := newReadBuffer(r)

	// 't' [int32 - length] [int16 - parameter count] [int32 - parameter] ...
	err := b.ReadTag('S')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	count, err := buf.ReadInt16()
	if err != nil {
		return nil, err
	}

	p := &ParameterDescription{
		OIDs: make([]int, count),
	}

	for i := 0; i < count; i++ {
		p.OIDs[i], err = b.ReadInt()
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *ParameterDescription) Encode() []byte {
	// 't' [int32 - length] [int16 - parameter count] [in32 - parameter] ...
	w := newWriteBuffer()
	w.WriteInt16(len(p.OIDs))
	for _, oid := range p.OIDs {
		w.WriteInt(oid)
	}
	w.Wrap('t')
	return w.Bytes()
}

func (p *ParameterDescription) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "ParameterDescription",
		"Payload": map[string]interface{}{
			"OIDs": p.OIDs,
		},
	}
}

func (p *ParameterDescription) String() string { return messageToString(p) }
