package pgproto

import (
	"io"
)

// CommandCompletion represents a server response message
type CommandCompletion struct {
	Tag []byte
}

func (c *CommandCompletion) server() {}

// ParseCommandCompletion will attempt to read a CommandCompletion message from the io.Reader
func ParseCommandCompletion(r io.Reader) (*CommandCompletion, error) {
	b := newReadBuffer(r)

	// 'C' [int32 - length] [tag] \0
	err := b.ReadTag('C')
	if err != nil {
		return nil, err
	}

	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	c := &CommandCompletion{}
	c.Tag, err = b.ReadString(true)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Encode will return the byte representation of this message
func (c *CommandCompletion) Encode() []byte {
	b := newWriteBuffer()
	b.WriteString(c.Tag, true)
	b.Wrap('C')
	return b.Bytes()
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "CommandCompletion",
//     "Payload": map[string]interface{}{
//       "Tag": <CommandCompletion.Tag>,
//     },
//   }
func (c *CommandCompletion) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "CommandCompletion",
		"Payload": map[string]string{
			"Tag": string(c.Tag),
		},
	}
}

func (c *CommandCompletion) String() string { return messageToString(c) }
