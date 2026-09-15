package server

import (
	"io"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/request"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode string
	Msg        string
}

func (h HandlerError) Text() string {
	return h.Msg
}
