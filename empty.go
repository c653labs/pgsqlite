package pgproto

import (
	"fmt"
	"io"
)

type EmptyQueryResponse struct{}

func (e *EmptyQueryResponse) server() {}

func ParseEmptyQueryResponse(r io.Reader) (*EmptyQueryResponse, error) {
	b := newReadBuffer(r)

	err := b.ReadTag('I')
	if err != nil {
		return nil, err
	}

	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}

	if l != 4 {
		return nil, fmt.Errorf("expected message length of 4")
	}

	return &EmptyQueryResponse{}, nil
}

func (e *EmptyQueryResponse) Encode() []byte {
	b := newWriteBuffer()
	b.Wrap('I')
	return b.Bytes()
}

func (e *EmptyQueryResponse) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "EmptyQueryResponse",
		"Payload": nil,
	}
}

func (e *EmptyQueryResponse) String() string { return messageToString(e) }
