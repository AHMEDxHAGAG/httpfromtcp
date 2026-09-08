package request

import (
	"errors"
	"strings"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func parseRequestLine(lines string) (RequestLine, error) {
	reqLineString := strings.Split(lines, "\r\n")[0]
	reqLineElems := strings.Split(reqLineString, " ")
	if len(reqLineElems) != 3 {
		return RequestLine{}, errors.New("wrong structure of request-line")
	}
	reqMethod, reqTarget, reqHTTPVersion := reqLineElems[0], reqLineElems[1], reqLineElems[2]
	err := validateRequestLine(reqMethod, reqTarget, reqHTTPVersion)
	if err != nil {
		return RequestLine{}, err
	}
	reqHTTPVersionNumber := strings.Split(reqHTTPVersion, "/")[1]
	return RequestLine{
		Method:        reqMethod,
		RequestTarget: reqTarget,
		HttpVersion:   reqHTTPVersionNumber,
	}, nil
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
	for i, v := range allowedVersions {
		if (i == len(allowedVersions)-1) && (parts[1] != v) {
			return errors.New("wrong http-version-number")
		}
		if parts[1] == v {
			break
		}
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
