package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/request"
)

type Server struct {
	state    atomic.Bool
	listener net.Listener
	handler  Handler
}

const (
	listening bool = true
	closed    bool = false
)

func newServer(listener net.Listener, handler Handler) *Server {
	server := &Server{listener: listener, handler: handler}
	server.state.Store(listening)
	return server
}

func Serve(handler Handler, port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := newServer(listener, handler)
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
			log.Fatal(err)
			return
		}
		go s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		_, err = conn.Write([]byte(err.Error()))
		if err != nil {
			log.Fatal(err)
			return
		}
		return

	}
	buffer := bytes.NewBuffer([]byte{})
	handlererr := s.handler(buffer, req)

	if handlererr != nil {
		writeErrHandlerOutput(conn, handlererr)
	} else {
		writeNormalHandlerOutput(conn, buffer)
	}
}
