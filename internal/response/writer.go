package response

import (
	"fmt"
	"io"

	"github.com/flames31/httpfromtcp/internal/headers"
)

type Writer struct {
	Writer io.Writer
}

func (w *Writer) WriteStatusLine(statusCode int) error {
	switch statusCode {
	case StatusOK:
		w.Writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		w.Writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusInternalSrvErr:
		w.Writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}

	return nil
}
func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	b := make([]byte, 0)
	headers.ForEach(func(k, v string) {
		b = append(b, []byte(fmt.Sprintf("%v: %v\r\n", k, v))...)
	})

	b = append(b, []byte("\r\n")...)
	w.Writer.Write(b)
	return nil
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	return n, err
}
