package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type AuthenticationRequestTestSuite struct {
	suite.Suite
}

func TestAuthenticationRequestTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationRequestTestSuite))
}

func (s *AuthenticationRequestTestSuite) Test_ParseAuthenticationRequest_MD5() {
	raw := []byte{
		// Tag
		'R',
		// Length
		'\x00', '\x00', '\x00', '\x0c',
		// Method
		'\x00', '\x00', '\x00', '\x05',
		// Salt
		'\xd1', '\x5b', '\x0e', '\x4f',
	}

	auth, err := pgproto.ParseAuthenticationRequest(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(auth)
	s.Equal(auth.Method, pgproto.AuthenticationMethodMD5)
	s.Equal(auth.Salt, []byte{'\xd1', '\x5b', '\x0e', '\x4f'})
	s.Equal(raw, auth.Encode())
}

func (s *AuthenticationRequestTestSuite) Test_ParseAuthenticationRequest_Empty() {
	auth, err := pgproto.ParseAuthenticationRequest(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(auth)
}

func (s *AuthenticationRequestTestSuite) Test_AuthenticationRequestEncode_MD5() {
	expected := []byte{
		// Tag
		'R',
		// Length
		'\x00', '\x00', '\x00', '\x0c',
		// Method
		'\x00', '\x00', '\x00', '\x05',
		// Salt
		'\xd1', '\x5b', '\x0e', '\x4f',
	}

	a := &pgproto.AuthenticationRequest{
		Method: pgproto.AuthenticationMethodMD5,
		Salt:   []byte{'\xd1', '\x5b', '\x0e', '\x4f'},
	}
	s.Equal(expected, a.Encode())
}

func (s *AuthenticationRequestTestSuite) Test_ParseAuthenticationRequest_OK() {
	raw := []byte{
		// Tag
		'R',
		// Length
		'\x00', '\x00', '\x00', '\x08',
		// Method
		'\x00', '\x00', '\x00', '\x00',
	}

	a, err := pgproto.ParseAuthenticationRequest(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(a)
	s.Equal(a.Method, pgproto.AuthenticationMethodOK)
	s.Nil(a.Salt)
	s.Equal(raw, a.Encode())
}

func (s *AuthenticationRequestTestSuite) Test_AuthenticationRequestEncode_OK() {
	expected := []byte{
		// Tag
		'R',
		// Length
		'\x00', '\x00', '\x00', '\x08',
		// Method
		'\x00', '\x00', '\x00', '\x00',
	}

	a := &pgproto.AuthenticationRequest{
		Method: pgproto.AuthenticationMethodOK,
		Salt:   nil,
	}
	s.Equal(expected, a.Encode())
}
