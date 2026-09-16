package response

import (
	"fmt"
	"io"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/constants"
	"github.com/AHMEDxHAGAG/httpfromtcp/internal/headers"
)

type writerState int

const (
	writingStatusLine  writerState = iota
	writingFieldLines  writerState = iota
	writingMessageBody writerState = iota
)

type Writer struct {
	w           io.Writer
	writerState writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, writerState: writingStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
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
	_, err := w.w.Write(statusLine)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writingFieldLines {
		return fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingFieldLines)
	}
	headersBuffer := []byte{}
	for key, value := range headers {
		header := fmt.Sprintf("%s: %s"+constants.CRLF, key, value)
		headersBuffer = append(headersBuffer, []byte(header)...)
	}
	headersBuffer = append(headersBuffer, []byte(constants.CRLF)...)
	_, err := w.w.Write(headersBuffer)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writingMessageBody {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingMessageBody)
	}
	n, err := w.w.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}
