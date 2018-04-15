package pgproto

import (
	"io"
)

type CopyOutResponse struct {
	Format        Format
	ColumnFormats []int
}

func (c *CopyOutResponse) server() {}

func ParseCopyOutResponse(r io.Reader) (*CopyOutResponse, error) {
	b := newReadBuffer(r)

	// 'H' [int32 - length] [int16 - count] [int16 - format] ...
	err := b.ReadTag('H')
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

	c := &CopyOutResponse{
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

func (c *CopyOutResponse) Encode() []byte {
	// 'H' [int32 - length] [int16 - count] [int16 - format] ...
	w := newWriteBuffer()
	w.WriteByte(byte(c.Format))
	w.WriteInt16(len(c.ColumnFormats))
	for _, format := range c.ColumnFormats {
		w.WriteInt16(format)
	}
	w.Wrap('H')
	return w.Bytes()
}

func (c *CopyOutResponse) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "CopyOutResponse",
		"Payload": map[string]interface{}{
			"Format":        c.Format,
			"ColumnFormats": c.ColumnFormats,
		},
	}
}

func (c *CopyOutResponse) String() string { return messageToString(c) }
