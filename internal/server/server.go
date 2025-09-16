package server

import (
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
}

func Serve(port int) (*Server, error) {
	portStr := strconv.Itoa(port)
	l, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	srv := &Server{
		listner:  l,
		isClosed: atomic.Bool{},
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
	_, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.WriteStatusLine(conn, 200)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	conn.Close()
}
