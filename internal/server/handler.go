package server

import (
	"fmt"
	"io"

	"github.com/flames31/httpfromtcp/internal/request"
	"github.com/flames31/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode int
	Msg        string
}

func (h *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	w.Write([]byte(fmt.Sprintf("%s", h.Msg)))
}
