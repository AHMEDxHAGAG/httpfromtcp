// Package request
package request

import (
	"fmt"
	"io"
)

type status int

const (
	intialized status = iota
	done
)

type Request struct {
	Status      status
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reqLine, err := parseRequestLine(buffer)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: reqLine,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.Status == intialized {
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
	} else if r.Status == done {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	} else {
		return 0, fmt.Errorf("error: unknown state")
	}
}
