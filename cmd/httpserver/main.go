package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/flames31/httpfromtcp/internal/request"
	"github.com/flames31/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		hErr := &server.HandlerError{
			StatusCode: 200,
		}

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			hErr.StatusCode = 400
			hErr.Msg = "Your problem is not my problem\n"
		case "/myproblem":
			hErr.StatusCode = 500
			hErr.Msg = "Woopsie, my bad\n"
		default:
			hErr.Msg = "All good, frfr\n"
		}

		return hErr
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
