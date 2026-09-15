package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	state    atomic.Bool
	listener net.Listener
}

const crlf = "\r\n"

const (
	listening bool = true
	closed    bool = false
)

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener: listener}
	server.state.Store(listening)
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.state.Store(closed)
	return nil
}

func (s *Server) listen() {
	for s.state.Load() {
		con, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("error: %s", err)
		}
		go s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	response := []byte("HTTP/1.1 200 OK" + crlf +
		"Content-Type: text/plain" + crlf +
		"Content-Length: 13" + crlf +
		crlf +
		"Hello World!\n")

	_, _ = conn.Write(response)
}
