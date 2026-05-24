package request

import (
	"bytes"
	"errors"
	"httpfromscratch/internal/headers"
	"io"
	"unicode"
)

type RequestState int

const (
	Initialized RequestState = iota
	ParsingHeaders
	Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       RequestState
}

func newRequest() *Request {
	return &Request{
		state:   Initialized,
		Headers: headers.NewHeaders(),
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var SEPARATOR = []byte("\r\n")

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.state {
		case Initialized:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.state = ParsingHeaders

		case ParsingHeaders:

			n, b, err := r.Headers.Parse(data[read:])

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if b {
				r.state = Done
			}

		case Done:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == Done
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	header := b[:idx]
	read := idx + len(SEPARATOR)

	headerparts := bytes.Split(header, []byte(" "))

	if len(headerparts) != 3 {
		return nil, 0, errors.New("malformed header")
	}

	// check method
	for _, r := range string(headerparts[0]) {
		if unicode.IsNumber(r) || unicode.IsLower(r) {
			return nil, 0, errors.New("malformed method")
		}
	}

	// check HTTP version
	httpVersion := bytes.Split(headerparts[2], []byte("/"))
	if len(httpVersion) != 2 || string(httpVersion[0]) != "HTTP" || string(httpVersion[1]) != "1.1" {
		return nil, 0, errors.New("unsupported http version")
	}

	rl := &RequestLine{
		Method:        string(headerparts[0]),
		RequestTarget: string(headerparts[1]),
		HttpVersion:   string(httpVersion[1]),
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {

		length, err := reader.Read(buf[bufLen:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = Done
				break
			}
			return nil, err
		}

		bufLen += length

		readN, err := request.parse(buf[:bufLen])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}

	return request, nil
}
