// Package repsonse
package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/constants"
	"github.com/AHMEDxHAGAG/httpfromtcp/internal/headers"
)

type StatusCode string

const (
	Success     StatusCode = "200"
	ClientError StatusCode = "400"
	ServerError StatusCode = "500"
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := []byte("HTTP/1.1 " + statusCode + " ")
	switch statusCode {
	case Success:
		statusLine = append(statusLine, []byte("OK")...)
	case ClientError:
		statusLine = append(statusLine, []byte("Bad Request")...)
	case ServerError:
		statusLine = append(statusLine, []byte("Internal Server Error")...)
	}
	statusLine = append(statusLine, constants.CRLF...)
	_, err := w.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headersBuffer := []byte{}
	for key, value := range headers {
		header := fmt.Sprintf("%s: %s"+constants.CRLF, key, value)
		headersBuffer = append(headersBuffer, []byte(header)...)
	}
	headersBuffer = append(headersBuffer, []byte(constants.CRLF)...)
	_, err := w.Write(headersBuffer)
	return err
}
