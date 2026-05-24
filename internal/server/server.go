package server

import (
	"fmt"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
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

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":"+fmt.Sprint(port))

	if err != nil {
		return nil, err
	}

	s := &Server{listener: l}

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
			}
			return
		}

		_, err = request.RequestFromReader(conn)

		if err != nil {
			log.Println("Invalid request:", err)
		}

		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, 200)
	if err != nil && !s.closed.Load() {
		log.Println("error during response writing:", err)
	}
	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil && !s.closed.Load() {
		log.Println("error during response writing:", err)
	}
}
