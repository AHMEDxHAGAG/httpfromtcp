// Package repsonse
package response

import (
	"strconv"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/headers"
)

type StatusCode string

const (
	Success     StatusCode = "200"
	ClientError StatusCode = "400"
	ServerError StatusCode = "500"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}
