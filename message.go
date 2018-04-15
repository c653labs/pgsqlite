package pgproto

import (
	"fmt"
	"io"
)

// Message is the main interface for all PostgreSQL messages
type Message interface {
	Encode() []byte
	AsMap() map[string]interface{}
	String() string
}

// ClientMessage is an interface describing all client side PostgreSQL messages (messages sent to the server)
type ClientMessage interface {
	Message
	client()
}

// ServerMessage is an interface describing all server side PostgreSQL messages (messages sent to the client)
type ServerMessage interface {
	Message
	server()
}

// ParseClientMessage will read the next ClientMessage from the provided io.Reader
func ParseClientMessage(r io.Reader) (ClientMessage, error) {
	// Create a buffer
	buf := newReadBuffer(r)

	// Look at the first byte to determine the type of message we have
	start, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	// Startup message:
	//   [int32 - length] [int32 - protocol] [[string]\0[string]\0] \0
	// Regular message
	//   [char - tag] [int32 - length] [payload]
	switch start {
	// TODO: We need to handle this case better, it might not always start with \x00
	//       We could just make calling `ParseStartupMessage` explicit
	case '\x00':
		msgReader, err := readStartupMessage(start, buf)
		if err != nil {
			return nil, err
		}
		return ParseStartupMessage(msgReader)
	default:
		msgReader, err := readMessage(start, buf)
		if err != nil {
			return nil, err
		}
		switch start {
		case 'p':
			// Password message
			return ParsePasswordMessage(msgReader)
		case 'Q':
			// Simple query
			return ParseSimpleQuery(msgReader)
		case 't':
			// Parameter description
			return ParseParameterDescription(msgReader)
		case 'B':
			// Binary parameters
			return ParseBinaryParameters(msgReader)
		case 'P':
			// Parse
			return ParseParse(msgReader)
		case 'E':
			// Execute
			return ParseExecute(msgReader)
		case 'H':
			// Flush
			return ParseFlush(msgReader)
		case 'S':
			// Sync
			return ParseSync(msgReader)
		case 'C':
			// Close
			return ParseClose(msgReader)
		case 'D':
			// Describe
			return ParseDescribe(msgReader)
		case 'X':
			// Termination
			return ParseTermination(msgReader)
		default:
			return nil, fmt.Errorf("unknown message tag '%c'", start)
		}
	}
}

// ParseServerMessage will read the next ServerMessage from the provided io.Reader
func ParseServerMessage(r io.Reader) (ServerMessage, error) {
	// Create a buffer
	buf := newReadBuffer(r)

	// Look at the first byte to determine the type of message we have
	start, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	msgReader, err := readMessage(start, buf)
	if err != nil {
		return nil, err
	}

	// Message
	//   [char - tag] [int32 - length] [payload]
	switch start {
	case 'R':
		// Authentication request
		return ParseAuthenticationRequest(msgReader)
	case 'S':
		// Parameter status
		return ParseParameterStatus(msgReader)
	case 'K':
		// Backend key data
		return ParseBackendKeyData(msgReader)
	case 'Z':
		// Ready for query
		return ParseReadyForQuery(msgReader)
	case 'C':
		// Command completion
		return ParseCommandCompletion(msgReader)
	case 'T':
		// Row description
		return ParseRowDescription(msgReader)
	case 't':
		// Parameter description
		return nil, fmt.Errorf("unhandled message tag %#v", start)
	case 'D':
		// Data row
		return ParseDataRow(msgReader)
	case 'I':
		// Empty query response
		return ParseEmptyQueryResponse(msgReader)
	case '1':
		// Parse complete
		return ParseParseComplete(msgReader)
	case '2':
		// Bind complete
		return ParseBindComplete(msgReader)
	case '3':
		// Close complete
		return ParseCloseComplete(msgReader)
	case 'W':
		// Copy both response
		return ParseCopyBothResponse(msgReader)
	case 'd':
		// Copy data
		return ParseCopyData(msgReader)
	case 'G':
		// Copy in response
		return ParseCopyInResponse(msgReader)
	case 'H':
		// Copy out response
		return ParseCopyOutResponse(msgReader)
	case 'V':
		// Function call response
		return nil, fmt.Errorf("unhandled message tag %#v", start)
	case 'n':
		// No data
		return ParseNoData(msgReader)
	case 'N':
		// Notice response
		return ParseNoticeResponse(msgReader)
	case 'A':
		// Notification response
		return ParseNotification(msgReader)
	case 'E':
		// Error message
		return ParseError(msgReader)
	default:
		return nil, fmt.Errorf("unknown message tag '%c'", start)
	}
}

func readStartupMessage(start byte, buf *readBuffer) (io.Reader, error) {
	// [int32 - length] [payload]
	// StartupMessage
	// Read the next 3 bytes, prepend with the 1 we already read to parse the length from this message
	b := make([]byte, 3)
	_, err := buf.Read(b)
	if err != nil {
		return nil, err
	}
	b = append([]byte{start}, b...)
	l := bytesToInt(b)

	// Read the rest of the message into a []byte
	// DEV: Subtract 4 to account for the length of the in32 we just read
	b = make([]byte, l-4)
	_, err = buf.Read(b)
	if err != nil {
		return nil, err
	}

	// Rebuild the message into a []byte
	w := newWriteBuffer()
	w.WriteInt(l)
	w.WriteBytes(b)
	return w.Reader(), nil
}

func readMessage(start byte, buf *readBuffer) (io.Reader, error) {
	// [char tag] [int32 length] [payload]
	// Parse length from the message
	l, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}

	// Read the rest of the message into a []byte
	// DEV: Subtract 4 to account for the length of the int32 we just read
	b := make([]byte, l-4)
	_, err = buf.Read(b)
	if err != nil {
		return nil, err
	}

	// Rebuild the message into a []byte
	w := newWriteBuffer()
	w.WriteByte(start)
	w.WriteInt(l)
	w.WriteBytes(b)
	return w.Reader(), nil
}
