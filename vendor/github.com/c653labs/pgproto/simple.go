package pgproto

import (
	"io"
)

type SimpleQuery struct {
	Query []byte
}

func (q *SimpleQuery) client() {}

func ParseSimpleQuery(r io.Reader) (*SimpleQuery, error) {
	b := newReadBuffer(r)

	// 'Q' [int32 - length] [query] \0
	err := b.ReadTag('Q')
	if err != nil {
		return nil, err
	}

	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	q := &SimpleQuery{}
	q.Query, err = b.ReadString(true)
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (q *SimpleQuery) Encode() []byte {
	b := newWriteBuffer()
	b.WriteString(q.Query, true)
	b.Wrap('Q')
	return b.Bytes()
}

func (q *SimpleQuery) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "SimpleQuery",
		"Payload": map[string]interface{}{
			"Query": string(q.Query),
		},
	}
}

func (q *SimpleQuery) String() string { return messageToString(q) }
