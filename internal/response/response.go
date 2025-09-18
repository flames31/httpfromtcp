package response

import (
	"fmt"
	"io"

	"github.com/flames31/httpfromtcp/internal/headers"
)

const (
	StatusOK             = 200
	StatusBadRequest     = 400
	StatusInternalSrvErr = 500
)

func WriteStatusLine(w io.Writer, statusCode int) error {
	switch statusCode {
	case StatusOK:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusInternalSrvErr:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}

	return nil
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	b := make([]byte, 0)
	headers.ForEach(func(k, v string) {
		b = append(b, []byte(fmt.Sprintf("%v: %v\r\n", k, v))...)
	})

	b = append(b, []byte("\r\n")...)
	w.Write(b)
	return nil
}
