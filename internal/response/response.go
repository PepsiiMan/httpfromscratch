package response

import (
	"fmt"
	"httpfromscratch/internal/headers"
	"io"
	"strings"
)

type StatusCode int

const (
	Ok                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := strings.Builder{}
	statusLine.WriteString("HTTP/1.1 ")
	fmt.Fprint(&statusLine, statusCode)
	switch statusCode {
	case Ok:
		statusLine.WriteString(" OK")
	case BadRequest:
		statusLine.WriteString(" Bad Request")
	case InternalServerError:
		statusLine.WriteString(" IntenralServerError")

	}
	statusLine.WriteString("\r\n")
	_, err := w.Write([]byte(statusLine.String()))

	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h["Content-Length"] = fmt.Sprint(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headersString := strings.Builder{}
	for k, v := range headers {
		headersString.WriteString(k)
		headersString.WriteString(": ")
		headersString.WriteString(v)
		headersString.WriteString("\r\n")
	}
	headersString.WriteString("\r\n")

	_, err := w.Write([]byte(headersString.String()))

	if err != nil {
		return err
	}

	return nil
}
