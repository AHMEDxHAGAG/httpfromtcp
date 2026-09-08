// Package request
package request

import (
	"io"
)

type Request struct {
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lines := string(buffer)
	reqLine, err := parseRequestLine(lines)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: reqLine,
	}, nil
}
