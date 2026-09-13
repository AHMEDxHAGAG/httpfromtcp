package request

import (
	"fmt"
	"strconv"
)

func parseMessageBody(r *Request, data []byte) (int, error) {
	contentLenString, _ := r.Headers.Get("Content-Length")
	contentLen, err := strconv.Atoi(contentLenString)
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
