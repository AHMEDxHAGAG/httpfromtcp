// Package headers
package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	ind := bytes.Index(data, []byte(crlf))
	if ind == -1 {
		return 0, false, nil
	}
	if ind == 0 {
		return len(data) - 2, true, nil
	}
	key, value, err := parseFieldLine(data)
	if err != nil {
		return 0, false, err
	}
	h[key] = value
	return len(data) - 2, false, nil
}

func parseFieldLine(data []byte) (key, value string, err error) {
	head := string(data)
	fieldName, fieldValue, colon := strings.Cut(head, ":")
	if !colon {
		return "", "", fmt.Errorf("error: colon is found in the field-name")
	}
	if strings.Contains(fieldName, " ") {
		return "", "", fmt.Errorf("error: whitespace is found in the field-name")
	}
	fieldValue = strings.TrimSpace(fieldValue)
	return fieldName, fieldValue, nil
}
