package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type StartupMessageTestSuite struct {
	suite.Suite
}

func TestStartupMessageTestSuite(t *testing.T) {
	suite.Run(t, new(StartupMessageTestSuite))
}

func (s *StartupMessageTestSuite) Test_ParseStartupMessage() {
	raw := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x27',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// "database" \0
		'\x64', '\x61', '\x74', '\x61', '\x62', '\x61', '\x73', '\x65', '\x00',
		// "db_name" \0
		'\x64', '\x62', '\x5f', '\x6e', '\x61', '\x6d', '\x65', '\x00',
		// "user" \0
		'\x75', '\x73', '\x65', '\x72', '\x00',
		// "pgproto" \0
		'\x70', '\x67', '\x70', '\x72', '\x6f', '\x74', '\x6f', '\x00',
		// ending
		'\x00',
	}
	startup, err := pgproto.ParseStartupMessage(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(startup)
	s.Equal(startup.Options["user"], []byte("pgproto"))
	s.Equal(startup.Options["database"], []byte("db_name"))
	s.Equal(raw, startup.Encode())
}

func (s *StartupMessageTestSuite) Test_ParseStartupMessage_NoOptions() {
	raw := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x09',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// ending
		'\x00',
	}

	startup, err := pgproto.ParseStartupMessage(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(startup)
	s.Equal(raw, startup.Encode())
}

func (s *StartupMessageTestSuite) Test_ParseStartupMessage_InvalidProtocolVersion() {
	raw := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x09',
		// Protocol
		'\x00', '\x00', '\x04', '\x00',
		// ending
		'\x00',
	}

	startup, err := pgproto.ParseStartupMessage(bytes.NewReader(raw))
	s.NotNil(err)
	s.Nil(startup)
}

func (s *StartupMessageTestSuite) Test_StartupMessageEncode() {
	expected := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x27',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// "database" \0
		'\x64', '\x61', '\x74', '\x61', '\x62', '\x61', '\x73', '\x65', '\x00',
		// "db_name" \0
		'\x64', '\x62', '\x5f', '\x6e', '\x61', '\x6d', '\x65', '\x00',
		// "user" \0
		'\x75', '\x73', '\x65', '\x72', '\x00',
		// "pgproto" \0
		'\x70', '\x67', '\x70', '\x72', '\x6f', '\x74', '\x6f', '\x00',
		// ending
		'\x00',
	}
	startup := &pgproto.StartupMessage{
		Options: make(map[string][]byte),
	}
	startup.Options["user"] = []byte("pgproto")
	startup.Options["database"] = []byte("db_name")
	s.Equal(expected, startup.Encode())
}
func (s *StartupMessageTestSuite) Test_StartupMessageEncode_NoOptions() {
	expected := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x09',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// ending
		'\x00',
	}
	startup := &pgproto.StartupMessage{}
	s.Equal(expected, startup.Encode())
}
