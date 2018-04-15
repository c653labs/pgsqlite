package pgproto

import (
	"fmt"
	"io"
)

type ReadyStatus int

const (
	READY_IDLE ReadyStatus = 73
)

func (r ReadyStatus) String() string {
	switch r {
	case READY_IDLE:
		return "Idle"
	}
	return "Unknown"
}

type ReadyForQuery struct {
	Status ReadyStatus
}

func (r *ReadyForQuery) server() {}

func ParseReadyForQuery(r io.Reader) (*ReadyForQuery, error) {
	b := newReadBuffer(r)

	// 'Z' [int32 - length] [byte - status]
	err := b.ReadTag('Z')
	if err != nil {
		return nil, err
	}

	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}
	if l != 5 {
		return nil, fmt.Errorf("unexpected message length")
	}

	i, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ReadyForQuery{
		Status: ReadyStatus(i),
	}, nil
}

func (r *ReadyForQuery) Encode() []byte {
	b := newWriteBuffer()
	b.WriteByte(byte(r.Status))
	b.Wrap('Z')
	return b.Bytes()
}

func (r *ReadyForQuery) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "ReadyForQuery",
		"Payload": map[string]interface{}{
			"Status": r.Status,
		},
	}
}

func (r *ReadyForQuery) String() string { return messageToString(r) }
