package server

import (
	"fmt"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

//const responseString = "HTTP/1.1 200 OK\r\n" +
//	"Content-Type: text/plain\r\n" +
//	"Content-Length: 13\r\n" +
//	"\r\n" +
//	"Hello World!\n"

//var response []byte = []byte(string(responseString))

type HandlerError struct {
	Message string
	Status  response.StatusCode
}

func (h HandlerError) writeError(w io.Writer) error {
	m := []byte(h.Message)
	writer := response.NewWriter(w)
	err := writer.WriteStatusLine(h.Status)
	if err != nil {
		return err
	}

	headers := response.GetDefaultHeaders(len(m))

	err = writer.WriteHeaders(headers)
	if err != nil {
		return err
	}

	_, err = writer.WriteBody(m)
	if err != nil {
		return err
	}

	return nil
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	s := &Server{listener: l, handler: handler}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			if !s.closed.Load() {
				log.Println("Connection error:", err)
				continue
			}
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	writer := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)

	if err != nil {
		errbytes := []byte(err.Error())
		writer.WriteStatusLine(response.BadRequest)
		writer.WriteHeaders(response.GetDefaultHeaders(len(errbytes)))
		writer.WriteBody(errbytes)
		return
	}

	//buffer := bytes.NewBuffer([]byte{})

	s.handler(writer, req)

	//if h != nil {
	//	err = h.writeError(conn)
	//	if err != nil && !s.closed.Load() {
	//		log.Println("error during response writing:", err)
	//	}
	//	return
	//}

	//headers := response.GetDefaultHeaders(buffer.Len())
	//err = writer.WriteStatusLine(200)
	//if err != nil && !s.closed.Load() {
	//	log.Println("error during response writing:", err)
	//}
	//err = writer.WriteHeaders(headers)
	//if err != nil && !s.closed.Load() {
	//	log.Println("error during response writing:", err)
	//}
	//_, err = writer.WriteBody(buffer.Bytes())
	//if err != nil && !s.closed.Load() {
	//	log.Println("error during response writing:", err)
	//}
}
