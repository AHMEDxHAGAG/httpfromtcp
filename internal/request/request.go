// Package request
package request

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/headers"
)

type status int

const (
	initialized status = iota
	parsingHeaders
	done
)

const (
	bufferSize = 4096
	crlf       = "\r\n"
)

type Request struct {
	status      status
	RequestLine RequestLine
	Headers     headers.Headers
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{status: initialized}
	buffer := make([]byte, bufferSize)
	readInd := 0
	for r.status != done {
		if readInd >= len(buffer) {
			buffer = slices.Grow(buffer, len(buffer))
			buffer = buffer[:cap(buffer)]
		}
		read, readErr := reader.Read(buffer[readInd:cap(buffer)])
		readInd += read
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return nil, readErr
		}
		totalCons := 0
		for r.status != done {
			consumed, err := r.parse(buffer[totalCons:readInd])
			if err != nil {
				return nil, err
			}
			if consumed == 0 {
				break
			}
			totalCons += consumed
		}
		copy(buffer, buffer[totalCons:readInd])
		readInd -= totalCons
		if errors.Is(readErr, io.EOF) && r.status != done {
			return nil, readErr
		}
	}
	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.status {
	case initialized:
		rq, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n != 0 {
			r.RequestLine = rq
			r.status = parsingHeaders
			r.Headers = headers.NewHeaders()
			return n, nil
		}
		return 0, nil
	case parsingHeaders:
		n, ok, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if ok {
			r.status = done
			return n, nil
		}
		return n, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
