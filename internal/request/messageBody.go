package request

import (
	"fmt"
	"strconv"
)

func parseMessageBody(r *Request, data []byte) (int, error) {
	contentLenString, _ := r.Headers.Get("Content-Length")
	contentLen, err := getContentLen(contentLenString)
	if err != nil {
		return 0, err
	}
	r.Body = append(r.Body, data...)
	if len(r.Body) == contentLen {
		r.status = done
	}
	if len(r.Body) > contentLen {
		return 0, fmt.Errorf("Length Of The Body Is Larger Than The Content Length")
	}
	return len(data), nil
}

func getContentLen(s string) (int, error) {
	for _, char := range s {
		if !(char >= '0' && char <= '9') {
			return 0, fmt.Errorf("content-length value is non numeric")
		}
	}
	no, err := strconv.Atoi(s)
	return no, err
}
