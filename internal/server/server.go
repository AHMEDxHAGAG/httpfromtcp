package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/response"
)

type Server struct {
	state    atomic.Bool
	listener net.Listener
}

const (
	listening bool = true
	closed    bool = false
)

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Fatalf("%s", err)
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Fatalf("%s", err)
	}
}
