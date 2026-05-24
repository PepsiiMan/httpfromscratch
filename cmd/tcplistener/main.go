package main

import (
	"fmt"
	"httpfromscratch/internal/request"
	"log"
	"net"
)

func main() {

	l, err := net.Listen("tcp", ":8182")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Connection Accepted")
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}
		fmt.Printf("Request line: \n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Printf("Headers: \n")
		for k, v := range r.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Printf("Body: \n")
		fmt.Printf("%s\n", string(r.Body))

		conn.Close()
		fmt.Println("Connection Closed")

	}

}
