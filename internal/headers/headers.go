package headers

import (
	"bytes"
	"errors"
)

var SEPARATOR = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, SEPARATOR)
	switch idx {
	case -1:
		return 0, false, nil
	case 0:
		return idx + len(SEPARATOR), true, nil
	}

	n = idx + len(SEPARATOR)

	header := data[:idx]

	headerparts := bytes.SplitN(header, []byte(":"), 2)

	if len(headerparts) != 2 {
		return
	}

	idx = bytes.Index(headerparts[0], []byte(" "))

	if idx != -1 {
		return 0, false, errors.New("malformed field-line")
	}
	key := headerparts[0]
	value := bytes.ReplaceAll(headerparts[1], []byte(" "), []byte{})
	h[string(key)] = string(value)

	return n, false, nil
}
