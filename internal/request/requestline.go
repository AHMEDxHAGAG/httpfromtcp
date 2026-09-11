package request

import (
	"bytes"
	"errors"
	"slices"
	"strings"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func parseRequestLine(buffer []byte) (rqline RequestLine, noOfBytes int, err error) {
	if !bytes.Contains(buffer, []byte(crlf)) {
		return RequestLine{}, 0, nil
	}
	reqLineBytes, _, _ := bytes.Cut(buffer, []byte(crlf))
	reqLineElems := strings.Split(string(reqLineBytes), " ")
	if len(reqLineElems) != 3 {
		return RequestLine{}, 0, errors.New("bad structured request-line")
	}
	reqMethod, reqTarget, reqHTTPVersion := reqLineElems[0], reqLineElems[1], reqLineElems[2]
	err = validateRequestLine(reqMethod, reqTarget, reqHTTPVersion)
	if err != nil {
		return RequestLine{}, 0, err
	}
	reqHTTPVersionNumber := strings.Split(reqHTTPVersion, "/")[1]
	return RequestLine{
		Method:        reqMethod,
		RequestTarget: reqTarget,
		HttpVersion:   reqHTTPVersionNumber,
	}, len(reqLineBytes) + 2, nil
}

func validateRequestLine(reqMethod, reqTarget, reqHTTPVersion string) error {
	err := validateMethod(reqMethod)
	if err != nil {
		return err
	}
	err = validateVersion(reqHTTPVersion)
	if err != nil {
		return err
	}
	err = validateTarget(reqTarget)
	if err != nil {
		return err
	}
	return nil
}

func validateTarget(reqTarget string) error {
	if !strings.Contains(reqTarget, "/") {
		return errors.New("wrong target must contain '/'")
	}
	return nil
}

func validateVersion(reqVersion string) error {
	allowedVersions := []string{"1.1"}
	parts := strings.Split(reqVersion, "/")
	if len(parts) != 2 {
		return errors.New("wrong http-version")
	}
	if parts[0] != "HTTP" {
		return errors.New("wrong http-name")
	}
	if !slices.Contains(allowedVersions, parts[1]) {
		return errors.New("wrong http-version-number")
	}
	return nil
}

func validateMethod(reqMethod string) error {
	for _, char := range reqMethod {
		if char >= 'A' && char <= 'Z' {
			continue
		}
		return errors.New("not all characters in the method are capital alphabetic")
	}
	return nil
}
