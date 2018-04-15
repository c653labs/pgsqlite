package pgproto

import (
	"io"
)

type Parse struct {
	Name  []byte
	Query []byte
	OIDs  []int
}

func (p *Parse) client() {}

func ParseParse(r io.Reader) (*Parse, error) {
	b := newReadBuffer(r)

	// 'P' [int32 - length] [string - Name] \0 [string - Query] \0 [int16 - parameter count] [int32 - parameter] ...
	err := b.ReadTag('P')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	p := &Parse{}

	p.Name, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	p.Query, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	count, err := buf.ReadInt16()
	if err != nil {
		return nil, err
	}

	p.OIDs = make([]int, count)
	for i := 0; i < count; i++ {
		p.OIDs[i], err = buf.ReadInt()
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Parse) Encode() []byte {
	// 'P' [int32 - length] [string - Name] \0 [string - Query] \0 [int16 - parameter count] [int32 - parameter] ...
	w := newWriteBuffer()
	w.WriteString(p.Name, true)
	w.WriteString(p.Query, true)
	w.WriteInt16(len(p.OIDs))
	for _, oid := range p.OIDs {
		w.WriteInt(oid)
	}
	w.Wrap('P')
	return w.Bytes()
}

func (p *Parse) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "Parse",
		"Payload": map[string]interface{}{
			"Name":  string(p.Name),
			"Query": string(p.Query),
			"OIDs":  p.OIDs,
		},
	}
}

func (p *Parse) String() string { return messageToString(p) }
