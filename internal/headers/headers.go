// Package headers
package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/constants"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (string, bool) {
	value, avaliable := h[strings.ToLower(key)]
	return value, avaliable
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Add(key, value string) {
	h[strings.ToLower(key)] = handleDuplicate(h, key, value)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	ind := bytes.Index(data, []byte(constants.CRLF))
	if ind == -1 {
		return 0, false, nil
	}
	if ind == 0 {
		return len(constants.CRLF), true, nil
	}
	key, value, err := parseFieldLine(data[:ind])
	if err != nil {
		return 0, false, err
	}
	h.Add(key, value)
	return ind + len(constants.CRLF), false, nil
}

func parseFieldLine(data []byte) (key, value string, err error) {
	head := string(data)
	fieldName, fieldValue, colon := strings.Cut(head, ":")
	if !colon {
		return "", "", fmt.Errorf("error: colon isn't found in the field-name")
	}
	if strings.Contains(fieldName, " ") {
		return "", "", fmt.Errorf("error: whitespace is found in the field-name")
	}
	fieldName = strings.ToLower(fieldName)
	if err := validateFieldName(fieldName); err != nil {
		return "", "", err
	}
	fieldValue = strings.TrimSpace(fieldValue)
	return fieldName, fieldValue, nil
}

func validateFieldName(fieldName string) error {
	if fieldName == "" {
		return fmt.Errorf("wrong fieldline structure: empty fieldname")
	}
	for _, char := range fieldName {
		if (char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			slices.Contains([]rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}, char) {
			continue
		}
		return fmt.Errorf("invalid character in the header key")
	}
	return nil
}

func handleDuplicate(h Headers, key, value string) string {
	if _, ok := h[key]; !ok {
		return value
	}
	valuesList := strings.Split(h[key], ", ")
	if found := slices.Contains(valuesList, value); found {
		value = h[key]
	} else {
		value = fmt.Sprintf("%s, %s", h[key], value)
	}
	return value
}
