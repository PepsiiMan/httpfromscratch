package response

import (
	"errors"
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

type WriterState int

const (
	writingStatusLine WriterState = iota
	writingHeaders
	writingBody
	writingDone
)

type Writer struct {
	w           io.Writer
	writerState WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:           w,
		writerState: writingStatusLine,
	}
}

func (writer *Writer) WriteStatusLine(statusCode StatusCode) error {
	if writer.writerState != writingStatusLine {
		return errors.New("writing response out of order")
	}

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
	_, err := writer.w.Write([]byte(statusLine.String()))

	if err != nil {
		return err
	}

	writer.writerState = writingHeaders

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprint(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func (writer *Writer) WriteHeaders(headers headers.Headers) error {
	if writer.writerState != writingHeaders {
		return errors.New("writing response out of order")
	}

	headersString := strings.Builder{}
	for k, v := range headers {
		headersString.WriteString(k)
		headersString.WriteString(": ")
		headersString.WriteString(v)
		headersString.WriteString("\r\n")
	}
	headersString.WriteString("\r\n")

	_, err := writer.w.Write([]byte(headersString.String()))

	if err != nil {
		return err
	}

	writer.writerState = writingBody

	return nil
}

func (writer *Writer) WriteBody(p []byte) (int, error) {
	if writer.writerState != writingBody {
		return 0, errors.New("writing response out of order")
	}

	n, err := writer.w.Write(p)

	writer.writerState = writingDone
	return n, err
}
