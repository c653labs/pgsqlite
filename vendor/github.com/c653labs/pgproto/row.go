package pgproto

import (
	"io"
)

type RowField struct {
	ColumnName   []byte
	TableOID     int
	ColumnIndex  int // int16
	TypeOID      int
	ColumnLength int //int16
	TypeModifier int
	Format       Format
}

func (f RowField) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"ColumnName":   f.ColumnName,
		"TableOID":     f.TableOID,
		"ColumnIndex":  f.ColumnIndex,
		"TypeOID":      f.TypeOID,
		"ColumnLength": f.ColumnLength,
		"TypeModifier": f.TypeModifier,
		"Format":       f.Format,
	}
}

type RowDescription struct {
	Fields []RowField
}

func (r *RowDescription) server() {}

func ParseRowDescription(r io.Reader) (*RowDescription, error) {
	b := newReadBuffer(r)

	// 'T' [int32 - length] [int32 - field count] ([string - column name]\0 [int32 - table oid] [int16 - column index] [int32 - type oid] [int16 - column length] [int32 - type modifier] [int16 - format])
	err := b.ReadTag('T')
	if err != nil {
		return nil, err
	}

	// Length - int
	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	// Field count - int16
	c, err := b.ReadInt16()
	if err != nil {
		return nil, err
	}

	rd := &RowDescription{
		Fields: make([]RowField, c),
	}
	for i := 0; i < c; i++ {
		// Column Name - string
		rd.Fields[i].ColumnName, err = b.ReadString(true)
		if err != nil {
			return nil, err
		}

		// Table OID - int
		rd.Fields[i].TableOID, err = b.ReadInt()
		if err != nil {
			return nil, err
		}

		// Column Index - int16
		rd.Fields[i].ColumnIndex, err = b.ReadInt16()
		if err != nil {
			return nil, err
		}

		// Type OID - int
		rd.Fields[i].TypeOID, err = b.ReadInt()
		if err != nil {
			return nil, err
		}

		// Column Length - int16
		rd.Fields[i].ColumnLength, err = b.ReadInt16()
		if err != nil {
			return nil, err
		}

		// Type Modifier - int
		rd.Fields[i].TypeModifier, err = b.ReadInt()
		if err != nil {
			return nil, err
		}

		// Format - int16
		format, err := b.ReadInt16()
		rd.Fields[i].Format = Format(format)
		if err != nil {
			return nil, err
		}
	}

	return rd, nil
}

func (r *RowDescription) Encode() []byte {
	b := newWriteBuffer()
	// Field count - int16
	b.WriteInt16(len(r.Fields))
	for _, f := range r.Fields {
		// Column Name - string
		b.WriteString(f.ColumnName, true)

		// Table OID - int
		b.WriteInt(f.TableOID)

		// Column Index - int16
		b.WriteInt16(f.ColumnIndex)

		// Type OID - int
		b.WriteInt(f.TypeOID)

		// Column Length - int16
		b.WriteInt16(f.ColumnLength)

		// Type Modifier - int
		b.WriteInt(f.TypeModifier)

		// Format - int16
		b.WriteInt16(int(f.Format))
	}

	b.Wrap('T')
	return b.Bytes()
}

func (r *RowDescription) AsMap() map[string]interface{} {
	fields := make([]map[string]interface{}, 0)
	for _, f := range r.Fields {
		fields = append(fields, f.AsMap())
	}
	return map[string]interface{}{
		"Type": "RowDescription",
		"Payload": map[string]interface{}{
			"Fields": fields,
		},
	}
}

func (r *RowDescription) String() string { return messageToString(r) }
