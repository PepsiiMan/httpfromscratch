package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromscratch/internal/headers"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"httpfromscratch/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		html := writeHTMLStub("500 Internal Server Error", "Internal Server Error", "Okay, you know what? this one is on me.")
		headers := response.GetDefaultHeaders(len(html))
		headers.Reset("Content-Type", "text/html")
		w.WriteStatusLine(response.InternalServerError)
		w.WriteHeaders(headers)
		w.WriteBody(html)
		return
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {

		target := req.RequestLine.RequestTarget

		res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
		if err != nil {
			html := writeHTMLStub("500 Internal Server Error", "Internal Server Error", "Okay, you know what? this one is on me.")
			headers := response.GetDefaultHeaders(len(html))
			headers.Reset("Content-Type", "text/html")
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(headers)
			w.WriteBody(html)

		} else {
			headers := response.GetDefaultHeaders(0)
			headers.Delete("Content-Length")
			headers.Set("Transfer-Encoding", "chunked")
			w.WriteStatusLine(response.Ok)
			w.WriteHeaders(headers)

			for {
				data := make([]byte, 1024)
				n, err := res.Body.Read(data)
				if err != nil {
					break
				}
				w.WriteChunkedBody([]byte(fmt.Sprintf("%x", n)))
				w.WriteChunkedBody(data[:n])

			}
			w.WriteChunkedBodyDone()
			return
		}
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/html") {

		target := req.RequestLine.RequestTarget

		res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
		if err != nil {
			html := writeHTMLStub("500 Internal Server Error", "Internal Server Error", "Okay, you know what? this one is on me.")
			headers := response.GetDefaultHeaders(len(html))
			headers.Reset("Content-Type", "text/html")
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(headers)
			w.WriteBody(html)
		} else {
			header := response.GetDefaultHeaders(0)
			header.Delete("Content-Length")
			w.WriteStatusLine(response.Ok)
			w.WriteHeaders(header)

			fullBody := []byte{}

			for {
				data := make([]byte, 1024)
				n, err := res.Body.Read(data)

				if err != nil {
					break
				}
				fullBody = append(fullBody, data[:n]...)
			}

			w.WriteBody(fullBody)
			hash := sha256.Sum256(fullBody)
			trailers := headers.NewHeaders()
			trailers.Set("X-Content-SHA256", string(hash[:]))
			trailers.Set("X-Content-Length", fmt.Sprint(len(fullBody)))

		}
	} else if req.RequestLine.Method == "GET" && req.RequestLine.RequestTarget == "/video" {
		video, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			log.Println(err)
			html := writeHTMLStub("500 Internal Server Error", "Internal Server Error", "Okay, you know what? this one is on me.")
			headers := response.GetDefaultHeaders(len(html))
			headers.Reset("Content-Type", "text/html")
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(headers)
			w.WriteBody(html)
		}
		headers := response.GetDefaultHeaders(len(video))
		headers.Reset("Content-Type", "video/mp4")
		w.WriteStatusLine(response.Ok)
		w.WriteHeaders(headers)
		w.WriteBody(video)
	} else {

		html := writeHTMLStub("200 OK", "Success!", "Your request was an absolute banger.")
		headers := response.GetDefaultHeaders(len(html))
		headers.Reset("Content-Type", "text/html")
		w.WriteStatusLine(response.Ok)
		w.WriteHeaders(headers)
		w.WriteBody(html)
	}
}

func writeHTMLStub(title string, h1 string, p string) []byte {
	return fmt.Appendf(nil, "<html>\n<head>\n<title>%s</title>\n</head>\n<body>\n<h1>%s</h1>\n<p>%s</p>\n</body>\n</html>\n", title, h1, p)
}
