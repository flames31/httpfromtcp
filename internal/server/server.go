package server

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/flames31/httpfromtcp/internal/request"
	"github.com/flames31/httpfromtcp/internal/response"
)

type Server struct {
	listner  net.Listener
	isClosed atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	portStr := strconv.Itoa(port)
	l, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	srv := &Server{
		listner:  l,
		isClosed: atomic.Bool{},
		handler:  handler,
	}

	go srv.listen()

	return srv, nil

}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listner.Close()
}

func (s *Server) listen() {
	for !s.isClosed.Load() {
		conn, err := s.listner.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go s.handle(conn)
	}
}
func (s *Server) handle(conn net.Conn) {
	if s.isClosed.Load() {
		return
	}
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := bytes.Buffer{}
	hErr := s.handler(&buf, req)
	hErr.Write(&buf)

	defHeaders := response.GetDefaultHeaders(buf.Len())
	response.WriteStatusLine(conn, response.StatusCode(hErr.StatusCode))
	response.WriteHeaders(conn, defHeaders)
	conn.Write(buf.Bytes())
	conn.Close()
}
