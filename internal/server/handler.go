package server

import (
	"bytes"
	"io"
	"log"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/request"
	"github.com/AHMEDxHAGAG/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Msg        string
}

func writeErrHandlerOutput(w io.Writer, handlererr *HandlerError) {
	if err := response.WriteStatusLine(w, handlererr.StatusCode); err != nil {
		log.Fatal(err)
		return
	}
	if err := response.WriteHeaders(w, response.GetDefaultHeaders(len(handlererr.Msg))); err != nil {
		log.Fatal(err)
		return
	}
	_, err := w.Write([]byte(handlererr.Msg))
	if err != nil {
		log.Fatal(err)
	}
}

func writeNormalHandlerOutput(w io.Writer, b *bytes.Buffer, resCode response.StatusCode) {
	if err := response.WriteStatusLine(w, resCode); err != nil {
		log.Fatal(err)
		return
	}
	if err := response.WriteHeaders(w, response.GetDefaultHeaders(b.Len())); err != nil {
		log.Fatal(err)
		return
	}
	_, err := w.Write(b.Bytes())
	if err != nil {
		log.Fatal(err)
		return
	}
}
