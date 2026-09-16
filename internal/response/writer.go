package response

import (
	"fmt"
	"io"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/constants"
	"github.com/AHMEDxHAGAG/httpfromtcp/internal/headers"
)

type writerState int

const (
	writingStatusLine writerState = iota
	writingFieldLines
	writingMessageBody
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
			w.writerState = writingMessageBody
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
	headersBuffer = append(headersBuffer, []byte(constants.CRLF)...)
	_, err = w.w.Write(headersBuffer)
	return err
}

func (w *Writer) WriteBody(p []byte) (n int, err error) {
	defer func() {
		if err == nil {
			w.writerState = finished
		}
	}()
	if w.writerState != writingMessageBody {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingMessageBody)
	}
	return w.w.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writingMessageBody {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, writingMessageBody)
	}
	line := fmt.Sprintf("%X%s%s%s", len(p), constants.CRLF, p, constants.CRLF)
	n, _ := w.w.Write([]byte(line))
	return n, nil
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != finished {
		return 0, fmt.Errorf("unordered writing state your current order is: %d and your request order is: %d", w.writerState, finished)
	}
	bodyDoneString := "0" + constants.CRLF + constants.CRLF
	bodyDone := []byte(bodyDoneString)
	return w.w.Write(bodyDone)
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
