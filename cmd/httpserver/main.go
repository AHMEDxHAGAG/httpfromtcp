package main

import (
	"log"
	"os"
	"os/signal"
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
	server, err := server.Serve(ToyHandler, port)
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

func ToyHandler(res *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {

	case "/yourproblem":
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
		_, _ = res.WriteBody([]byte(body))

	case "/myproblem":
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
		_, _ = res.WriteBody([]byte(body))

	default:
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
		_, _ = res.WriteBody([]byte(body))

	}
}
