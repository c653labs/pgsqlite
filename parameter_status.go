package pgproto

import (
	"io"
)

type ParameterStatus struct {
	Name  []byte
	Value []byte
}

func (p *ParameterStatus) server() {}

func ParseParameterStatus(r io.Reader) (*ParameterStatus, error) {
	b := newReadBuffer(r)

	// 'S' [int32 - length] [string] \0 [string] \0
	err := b.ReadTag('S')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	p := &ParameterStatus{}

	p.Name, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	p.Value, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *ParameterStatus) Encode() []byte {
	// 'S' [int32 - length] [string] \0 [string] \0
	w := newWriteBuffer()
	w.WriteString(p.Name, true)
	w.WriteString(p.Value, true)
	w.Wrap('S')
	return w.Bytes()
}

func (p *ParameterStatus) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "ParameterStatus",
		"Payload": map[string]interface{}{
			"Name":  string(p.Name),
			"Value": string(p.Value),
		},
	}
}

func (p *ParameterStatus) String() string { return messageToString(p) }
