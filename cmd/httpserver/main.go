package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/request"
	"github.com/AHMEDxHAGAG/httpfromtcp/internal/response"
	"github.com/AHMEDxHAGAG/httpfromtcp/internal/server"
)

const (
	port        = 42069
	asyncSignal = 1
)

func main() {
	server, err := server.Serve(Handler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer func() { _ = server.Close() }()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, asyncSignal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func Handler(res *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		subHandler400(res, req)
	case "/myproblem":
		subHandler500(res, req)
	case "/video":
		VideoHandler(res, req)
	default:
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
			ProxyHandler(res, req)
		} else {
			subHandler200(res, req)
		}
	}
}

func subHandler400(res *response.Writer, req *request.Request) {
	_ = res.WriteStatusLine(response.ClientError)
	body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)

	header := response.GetDefaultHeaders(len(body))
	header.Set("Content-Type", "text/html")
	_ = res.WriteHeaders(header)
	_, _ = res.WriteBody(body)

}

func subHandler200(res *response.Writer, req *request.Request) {
	_ = res.WriteStatusLine(response.Success)

	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)

	header := response.GetDefaultHeaders(len(body))
	header.Set("Content-Type", "text/html")
	_ = res.WriteHeaders(header)
	_, _ = res.WriteBody(body)
}

func subHandler500(res *response.Writer, req *request.Request) {
	_ = res.WriteStatusLine(response.ServerError)
	body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)

	header := response.GetDefaultHeaders(len(body))
	header.Set("Content-Type", "text/html")
	_ = res.WriteHeaders(header)
	_, _ = res.WriteBody(body)
}

func ProxyHandler(res *response.Writer, req *request.Request) {
	if err := res.WriteStatusLine(response.Success); err != nil {
		response.WriteServerError(res, err)
		return
	}
	header := response.GetDefaultHeaders(0)
	header.UnSet("Content-Length")
	header.Set("Transfer-Encoding", "chunked")
	header.Add("Trailer", "X-Content-SHA256")
	header.Add("Trailer", "X-Content-Length")
	if err := res.WriteHeaders(header); err != nil {
		response.WriteServerError(res, err)
		return
	}
	resp, err := http.Get("https://httpbingo.org/" + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/"))
	if err != nil {
		response.WriteServerError(res, err)
		return
	}
	defer resp.Body.Close()
	body := []byte{}
	buffer := make([]byte, 1024)
	for {
		n, readerr := resp.Body.Read(buffer)
		body = append(body, buffer[:n]...)
		if _, err := res.WriteChunkedBody(buffer[:n]); err != nil {
			response.WriteServerError(res, err)
			return
		}
		if readerr != nil {
			if errors.Is(readerr, io.EOF) {
				break
			}
			response.WriteServerError(res, readerr)
			return
		}
	}
	if _, err := res.WriteChunkedBodyDone(header); err != nil {
		response.WriteServerError(res, err)
		return
	}
	hash := sha256.Sum256(body)
	header.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
	conLen := strconv.Itoa(len(body))
	header.Set("X-Content-Length", conLen)
	if _, err := res.WriteTrailers(header); err != nil {
		response.WriteServerError(res, err)
		return
	}
}

func VideoHandler(res *response.Writer, req *request.Request) {
	body, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		response.WriteServerError(res, err)
		return
	}
	_ = res.WriteStatusLine(response.Success)
	header := response.GetDefaultHeaders(len(body))
	header.Set("Content-Type", "video/mp4")
	_ = res.WriteHeaders(header)
	_, _ = res.WriteBody(body)
}
