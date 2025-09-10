package main

import (
	"fmt"
	"net"
	"os"

	"github.com/flames31/httpfromtcp/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		fmt.Println("Connection has been accepted!")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		request.Headers.ForEach(func(k, v string) {
			fmt.Printf("- %s: %s\n", k, v)
		})
		fmt.Printf("Body:\n%s\n", string(request.Body))
	}
}
