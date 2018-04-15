package pgproto

import (
	"io"
)

type CopyInResponse struct {
	Format        Format
	ColumnFormats []int
}

func (c *CopyInResponse) server() {}

func ParseCopyInResponse(r io.Reader) (*CopyInResponse, error) {
	b := newReadBuffer(r)

	// 'G' [int32 - length] [int16 - count] [int16 - format] ...
	err := b.ReadTag('G')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	format, err := buf.ReadByte()

	count, err := buf.ReadInt16()
	if err != nil {
		return nil, err
	}

	c := &CopyInResponse{
		Format:        Format(format),
		ColumnFormats: make([]int, count),
	}

	for i := 0; i < count; i++ {
		c.ColumnFormats[i], err = buf.ReadInt16()
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *CopyInResponse) Encode() []byte {
	// 'G' [int32 - length] [int16 - count] [int16 - format] ...
	w := newWriteBuffer()
	w.WriteByte(byte(c.Format))
	w.WriteInt16(len(c.ColumnFormats))
	for _, format := range c.ColumnFormats {
		w.WriteInt16(format)
	}
	w.Wrap('G')
	return w.Bytes()
}

func (c *CopyInResponse) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "CopyInResponse",
		"Payload": map[string]interface{}{
			"Format":        c.Format,
			"ColumnFormats": c.ColumnFormats,
		},
	}
}

func (c *CopyInResponse) String() string { return messageToString(c) }
