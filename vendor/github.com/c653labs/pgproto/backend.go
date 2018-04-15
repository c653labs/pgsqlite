package pgproto

import (
	"io"
)

// BackendKeyData is a server response message
type BackendKeyData struct {
	PID int
	Key int
}

func (b *BackendKeyData) server() {}

// ParseBackendKeyData is used to parse a BackendKeyData message from an io.Reader
func ParseBackendKeyData(r io.Reader) (*BackendKeyData, error) {
	buf := newReadBuffer(r)

	// 'K' [int32 - length] [int32 - pid] [in32 - key]
	err := buf.ReadTag('K')
	if err != nil {
		return nil, err
	}

	buf, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	pid, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}

	key, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	return &BackendKeyData{
		PID: pid,
		Key: key,
	}, nil
}

// Encode will return the byte representation of this message
func (b *BackendKeyData) Encode() []byte {
	buf := newWriteBuffer()
	buf.WriteInt(b.PID)
	buf.WriteInt(b.Key)
	buf.Wrap('K')
	return buf.Bytes()
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "BackendKeyData",
//     "Payload": map[string]interface{}{
//       "PID": <BackendKeyData.PID>,
//       "Key": <BackendKeyData.Key>,
//     },
//   }
func (b *BackendKeyData) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "BackendKeyData",
		"Payload": map[string]interface{}{
			"PID": b.PID,
			"Key": b.Key,
		},
	}
}

func (b *BackendKeyData) String() string { return messageToString(b) }
