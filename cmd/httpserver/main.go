package main

import (
	"io"
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

func ToyHandler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: string(response.ClientError),
			Msg:        "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: string(response.ServerError),
			Msg:        "Woopsie, my bad\n",
		}
	default:
		_, err := w.Write([]byte("All good, frfr\n"))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
}
