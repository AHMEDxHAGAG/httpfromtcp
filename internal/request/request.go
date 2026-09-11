// Package request
package request

import (
	"errors"
	"io"
	"slices"
)

type status int

const (
	intialized status = iota
	done
)

const (
	bufferSize = 4096
	crlf       = "\r\n"
)

type Request struct {
	status      status
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{status: intialized}
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
		consumed, err := r.parse(buffer)
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[consumed:])
		readInd -= consumed
		if errors.Is(readErr, io.EOF) {
			r.status = done
			break
		}
	}
	return r, nil
}
