package server

import (
	"fmt"
	"io"

	"github.com/flames31/httpfromtcp/internal/request"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode int
	Msg        string
}

func (h *HandlerError) Write(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("%s", h.Msg)))
}
