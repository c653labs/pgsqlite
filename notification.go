package pgproto

import (
	"io"
)

type Notification struct {
	PID     int
	Channel []byte
	Payload []byte
}

func (n *Notification) server() {}

func ParseNotification(r io.Reader) (*Notification, error) {
	buf := newReadBuffer(r)

	// 'A' [int32 - length] [int32 - pid] [string - channel] \0 [string - payload] \0
	err := buf.ReadTag('A')
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

	channel, err := buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	payload, err := buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	return &Notification{
		PID:     pid,
		Channel: channel,
		Payload: payload,
	}, nil
}

func (n *Notification) Encode() []byte {
	// 'A' [int32 - length] [int32 - pid] [string - channel] \0 [string - payload] \0
	buf := newWriteBuffer()
	buf.WriteInt(n.PID)
	buf.WriteString(n.Channel, true)
	buf.WriteString(n.Payload, true)
	buf.Wrap('N')
	return buf.Bytes()
}

func (n *Notification) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "Notification",
		"Payload": map[string]interface{}{
			"PID":     n.PID,
			"Channel": string(n.Channel),
			"Payload": string(n.Payload),
		},
	}
}

func (n *Notification) String() string { return messageToString(n) }
