package request

import (
	"bytes"
	"errors"
	"httpfromscratch/internal/headers"
	"io"
	"strconv"
	"unicode"
)

type RequestState int

const (
	Initialized RequestState = iota
	ParsingHeaders
	ParsingBody
	Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       RequestState
	Body        []byte
}

func newRequest() *Request {
	return &Request{
		state:   Initialized,
		Headers: headers.NewHeaders(),
		Body:    []byte{},
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
				r.state = ParsingBody
			}

		case ParsingBody:
			if _, ok := r.Headers["content-length"]; !ok {
				r.state = Done
				break outer
			}
			r.Body = append(r.Body, data[read:]...)

			read += len(data[read:])

			cl, err := strconv.Atoi(r.Headers.Get("Content-Length"))
			if err != nil {
				return read, errors.New("malformed content-length")
			}

			if len(r.Body) > cl {
				return read, errors.New("body larger than declared content-length")
			}

			if len(r.Body) == cl {
				r.state = Done
			}

			if len(r.Body) < cl {
				break outer
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
				if request.state == ParsingBody {
					return nil, errors.New("body shorter than declared content-length")
				}
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
