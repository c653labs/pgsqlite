package pgproto

import (
	"bytes"
	"io"
)

type PasswordMessage struct {
	Password []byte
}

func (p *PasswordMessage) client() {}

func ParsePasswordMessage(r io.Reader) (*PasswordMessage, error) {
	b := newReadBuffer(r)

	// 'p' [int32 - length] [string] \0
	err := b.ReadTag('p')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	p := &PasswordMessage{}

	p.Password, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *PasswordMessage) PasswordValid(user []byte, password []byte, salt []byte) bool {
	hash := HashPassword(user, password, salt)
	return bytes.Equal(p.Password, hash)
}

func (p *PasswordMessage) SetPassword(user []byte, password []byte, salt []byte) {
	p.Password = HashPassword(user, password, salt)
}

func (p *PasswordMessage) Encode() []byte {
	// 'p' [int32 - length] [string] \0
	w := newWriteBuffer()
	w.WriteString(p.Password, true)
	w.Wrap('p')
	return w.Bytes()
}

func (p *PasswordMessage) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "PasswordMessage",
		"Payload": map[string]interface{}{
			"Password": p.Password,
		},
	}
}
func (p *PasswordMessage) String() string { return messageToString(p) }
