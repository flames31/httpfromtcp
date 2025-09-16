package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/flames31/httpfromtcp/internal/request"
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

	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"))
	conn.Close()
}
