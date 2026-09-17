package response

import (
	"fmt"
	"io"
	"strings"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/constants"
	"github.com/AHMEDxHAGAG/httpfromtcp/internal/headers"
)

type writerState int

const (
	writingStatusLine writerState = iota
	writingFieldLines
	writingBody
	writingTrailers
	finished
)

type Writer struct {
	w           io.Writer
	writerState writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, writerState: writingStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) (err error) {
	defer func() {
		if err == nil {
			w.writerState = writingFieldLines
		}
	}()
	if w.writerState != writingStatusLine {
		return fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingStatusLine)
	}
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
	_, err = w.w.Write(statusLine)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) (err error) {
	defer func() {
		if err == nil {
			w.w.Write([]byte(constants.CRLF))
			w.writerState = writingBody
		}
	}()
	if w.writerState != writingFieldLines {
		return fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingFieldLines)
	}
	headersBuffer := []byte{}
	for key, value := range headers {
		header := fmt.Sprintf("%s: %s"+constants.CRLF, key, value)
		headersBuffer = append(headersBuffer, []byte(header)...)
	}
	_, err = w.w.Write(headersBuffer)
	return err
}

func (w *Writer) WriteBody(p []byte) (n int, err error) {
	defer func() {
		if err == nil {
			w.writerState = finished
		}
	}()
	if w.writerState != writingBody {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingBody)
	}
	return w.w.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (n int, err error) {
	if w.writerState != writingBody {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingBody)
	}
	if len(p) == 0 {
		return w.w.Write([]byte("0" + constants.CRLF))
	}
	line := fmt.Sprintf("%X%s%s%s", len(p), constants.CRLF, p, constants.CRLF)
	return w.w.Write([]byte(line))
}
func (w *Writer) WriteChunkedBodyDone(h headers.Headers) (n int, err error) {
	defer func() {
		if err == nil {
			w.writerState = writingTrailers
		}
	}()
	if w.writerState != writingBody {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingBody)
	}
	bodyDoneString := "0" + constants.CRLF
	if _, ok := h.Get("Trailer"); !ok {
		bodyDoneString += constants.CRLF
	}
	bodyDone := []byte(bodyDoneString)
	return w.w.Write(bodyDone)
}

func (w *Writer) WriteTrailers(h headers.Headers) (n int, err error) {
	defer func() {
		if err == nil {
			w.w.Write([]byte(constants.CRLF))
			w.writerState = finished
		}
	}()
	if w.writerState != writingTrailers {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingTrailers)
	}
	if _, found := h.Get("Trailer"); !found {
		return 0, nil
	}
	trailersString, _ := h.Get("Trailer")
	trailers := strings.Split(trailersString, ", ")
	headersBuffer := []byte{}
	for _, key := range trailers {
		value, _ := h.Get(key)
		header := fmt.Sprintf("%s: %s"+constants.CRLF, key, value)
		headersBuffer = append(headersBuffer, []byte(header)...)
	}
	return w.w.Write(headersBuffer)
}

func WriteClientError(writer *Writer, err error) {
	body := []byte(err.Error())
	_ = writer.WriteStatusLine(ClientError)
	header := GetDefaultHeaders(len(body))
	_ = writer.WriteHeaders(header)
	_, _ = writer.WriteBody(body)
}

func WriteServerError(writer *Writer, err error) {
	body := []byte(err.Error())
	_ = writer.WriteStatusLine(ServerError)
	header := GetDefaultHeaders(len(body))
	_ = writer.WriteHeaders(header)
	_, _ = writer.WriteBody(body)
}
