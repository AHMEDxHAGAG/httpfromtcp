// Package request
package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
)

type status int

const (
	intialized status = iota
	done
)

const (
	bufferSize = 8
	crlf       = "\r\n"
)

type Request struct {
	Status      status
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{Status: intialized}
	buffer := make([]byte, bufferSize)
	readInd, sizeT := 0, bufferSize
	for r.Status != done {
		n, err := reader.Read(buffer[readInd:cap(buffer)])
		readInd += n
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		if errors.Is(err, io.EOF) {
			r.Status = done
			break
		}
		_, err = r.parse(buffer)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			buffer = slices.Grow(buffer, sizeT)
			sizeT *= 2
			buffer = buffer[:cap(buffer)]
		}

	}
	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.Status {
	case intialized:
		rq, no, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if no != 0 {
			r.RequestLine = rq
			r.Status = done
			return no, nil
		}
		return 0, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
