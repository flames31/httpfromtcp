package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/flames31/httpfromtcp/internal/request"
	"github.com/flames31/httpfromtcp/internal/response"
	"github.com/flames31/httpfromtcp/internal/server"
)

const port = 42069

var (
	req200 = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	req400 = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	req500 = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
)

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := req200
		statusCode := response.StatusOK
		if req.RequestLine.RequestTarget == "/yourproblem" {
			body = req400
			statusCode = response.StatusBadRequest
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			body = req500
			statusCode = response.StatusInternalSrvErr
		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {
			target := req.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
			if err != nil {
				body = req500
				statusCode = response.StatusInternalSrvErr
			} else {
				h.Delete("content-length")
				h.Set("transfer-encoding", "chunked")
				w.WriteHeaders(h)

				buf := make([]byte, 1024)
				for {
					_, err := res.Body.Read(buf)
					if err != nil {
						break
					}

					w.WriteChunkedBody(buf)
				}

				w.WriteChunkedBodyDone()
				return
			}
		}

		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html")
		w.WriteStatusLine(statusCode)
		w.WriteHeaders(h)
		w.WriteBody(body)

	})
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
