package pgproto

import "io"

type NoticeResponse Error

func ParseNoticeResponse(r io.Reader) (*NoticeResponse, error) {
	e, err := ParseError(r)
	if err != nil {
		return nil, err
	}

	return (*NoticeResponse)(e), nil
}

func (n *NoticeResponse) server() {}

func (n *NoticeResponse) Encode() []byte {
	return encodeError((*Error)(n), 'N')
}

func (n *NoticeResponse) AsMap() map[string]interface{} {
	return errorMap((*Error)(n), "NoticeResponse")
}
func (n *NoticeResponse) String() string { return messageToString(n) }
