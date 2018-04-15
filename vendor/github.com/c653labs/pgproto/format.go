package pgproto

type Format uint8

const (
	FormatText   Format = 0
	FormatBinary        = 1
)

func (f Format) String() string {
	switch f {
	case FormatText:
		return "Text"
	case FormatBinary:
		return "Binary"
	}
	return "Unknown"
}
