package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type CommandCompletionTestSuite struct {
	suite.Suite
}

func TestCommandCompletionTestSuite(t *testing.T) {
	suite.Run(t, new(CommandCompletionTestSuite))
}

func (s *CommandCompletionTestSuite) Test_ParseCommandCompletion_MD5() {
	raw := []byte{
		// Tag
		'C',
		// Length
		'\x00', '\x00', '\x00', '\x0f',
		// Tag
		'\x73', '\x65', '\x6c', '\x65', '\x63', '\x74', '\x20', '\x31', '\x32', '\x31',
		// \0
		'\x00',
	}

	command, err := pgproto.ParseCommandCompletion(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(command)
	s.Equal(command.Tag, []byte("select 121"))
	s.Equal(raw, command.Encode())
}

func (s *CommandCompletionTestSuite) Test_ParseCommandCompletion_Empty() {
	command, err := pgproto.ParseCommandCompletion(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(command)
}
func (s *CommandCompletionTestSuite) Test_CommandCompletionEncode() {
	expected := []byte{
		// Tag
		'C',
		// Length
		'\x00', '\x00', '\x00', '\x0f',
		// Tag
		'\x73', '\x65', '\x6c', '\x65', '\x63', '\x74', '\x20', '\x31', '\x32', '\x31',
		// \0
		'\x00',
	}

	c := &pgproto.CommandCompletion{
		Tag: []byte("select 121"),
	}
	s.Equal(expected, c.Encode())
}
