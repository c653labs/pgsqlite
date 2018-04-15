package pgproto

// ObjectType represents an object type
type ObjectType byte

const (
	// ObjectTypePreparedStatement represents a Prepared statement object type
	ObjectTypePreparedStatement ObjectType = 'S'

	// ObjectTypePortal represents a Portalobject type
	ObjectTypePortal = 'P'
)

func (o ObjectType) String() string {
	switch o {
	case ObjectTypePreparedStatement:
		return "PreparedStatement"
	case ObjectTypePortal:
		return "Portal"
	}
	return "Uknown"
}
