package main

import (
	"fmt"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"httpfromscratch/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 8182

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

var handler server.Handler = func(w *response.Writer, req *request.Request) {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		html := writeHTMLStub("400 Bad Request", "Bad Request", "Your request honestly kinda sucked.")
		headers := response.GetDefaultHeaders(len(html))
		headers.Reset("Content-Type", "text/html")
		w.WriteStatusLine(response.BadRequest)
		w.WriteHeaders(headers)
		w.WriteBody(html)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		html := writeHTMLStub("500 Internal Server Error", "Internal Server Error", "Okay, you know what? this one is on me.")
		headers := response.GetDefaultHeaders(len(html))
		headers.Reset("Content-Type", "text/html")
		w.WriteStatusLine(response.InternalServerError)
		w.WriteHeaders(headers)
		w.WriteBody(html)
		return
	}

	html := writeHTMLStub("200 OK", "Success!", "Your request was an absolute banger.")
	headers := response.GetDefaultHeaders(len(html))
	headers.Reset("Content-Type", "text/html")
	w.WriteStatusLine(response.Ok)
	w.WriteHeaders(headers)
	w.WriteBody(html)
}

func writeHTMLStub(title string, h1 string, p string) []byte {
	return fmt.Appendf(nil, "<html>\n<head>\n<title>%s</title>\n</head>\n<body>\n<h1>%s</h1>\n<p>%s</p>\n</body>\n</html>\n", title, h1, p)
}
