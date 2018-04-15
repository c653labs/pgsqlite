package pgproto

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

func bytesToInt(buf []byte) int {
	return int(int32(binary.BigEndian.Uint32(buf)))
}

func bytesToInt16(buf []byte) int {
	return int(int16(binary.BigEndian.Uint16(buf)))
}

// HashPassword helper function is used to compute the hash of a user's password
func HashPassword(user []byte, password []byte, salt []byte) []byte {
	digest := md5.New()
	digest.Write(password)
	digest.Write(user)
	pwdhash := digest.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(pwdhash)))
	hex.Encode(dst, pwdhash)

	digest = md5.New()
	digest.Write(dst)
	digest.Write(salt)

	hash := digest.Sum(nil)
	dst = make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(dst, hash)

	return append([]byte("md5"), dst...)
}

// WriteTo helper function is used to write the binary representation of a Message to an io.Writer
func WriteTo(m Message, w io.Writer) (int64, error) {
	n, err := w.Write(m.Encode())
	return int64(n), err
}

func messageToString(m Message) string {
	data := m.AsMap()

	t, _ := data["Type"]
	p, _ := data["Payload"]

	str := fmt.Sprintf("%s<", t)
	if p, ok := p.(map[string]interface{}); ok {
		first := true
		for k, v := range p {
			if !first {
				str += ", "
			}
			str += fmt.Sprintf("%s=%#v", k, v)
			first = false
		}
	}

	return str + ">"
}
