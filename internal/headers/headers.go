package headers

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

const specialChars = "!#$%&'*+-.^_`|~"

func isSpecialChar(r rune) bool {
	return strings.ContainsRune(specialChars, r)
}

var SEPARATOR = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Reset(key string, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Set(key string, value string) {
	key = strings.ToLower(key)
	val, ok := h[key]
	if ok {
		h[key] = val + ", " + value
	} else {
		h[key] = value
	}
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
		return 0, false, errors.New("malformed field-line")
	}

	idx = bytes.Index(headerparts[0], []byte(" "))

	if idx != -1 {
		return 0, false, errors.New("malformed field-line")
	}
	key := strings.ToLower(string(headerparts[0]))
	value := string(bytes.ReplaceAll(headerparts[1], []byte(" "), []byte{}))

	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !isSpecialChar(r) {
			return 0, false, errors.New("invalid character in field-line")
		}
	}

	if _, ok := h[key]; ok {
		h[key] = h[key] + ", " + value
	} else {
		h[key] = value
	}

	return n, false, nil
}
