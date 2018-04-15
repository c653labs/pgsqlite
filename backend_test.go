package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type BackendKeyDataTestSuite struct {
	suite.Suite
}

func TestBackendKeyDataTestSuite(t *testing.T) {
	suite.Run(t, new(BackendKeyDataTestSuite))
}

func (s *BackendKeyDataTestSuite) Test_ParseBackendKeyData_MD5() {
	raw := []byte{
		// Tag
		'K',
		// Length
		'\x00', '\x00', '\x00', '\x0c',
		// PID
		'\x00', '\x00', '\x04', '\xd2',
		// Key
		'\x00', '\x00', '\x04', '\xd2',
	}

	backend, err := pgproto.ParseBackendKeyData(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(backend)
	s.Equal(backend.PID, 1234)
	s.Equal(backend.Key, 1234)
	s.Equal(raw, backend.Encode())
}

func (s *BackendKeyDataTestSuite) Test_ParseBackendKeyData_Empty() {
	backend, err := pgproto.ParseBackendKeyData(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(backend)
}
func (s *BackendKeyDataTestSuite) Test_BackendKeyDataEncode() {
	expected := []byte{
		// Tag
		'K',
		// Length
		'\x00', '\x00', '\x00', '\x0c',
		// PID
		'\x00', '\x00', '\x04', '\xd2',
		// Key
		'\x00', '\x00', '\x04', '\xd2',
	}

	b := &pgproto.BackendKeyData{
		PID: 1234,
		Key: 1234,
	}
	s.Equal(expected, b.Encode())
}
