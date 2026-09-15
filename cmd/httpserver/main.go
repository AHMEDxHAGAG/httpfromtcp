package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/server"
)

const (
	port        = 42069
	asyncSignal = 1
)

func main() {
	server, err := server.Serve(port)
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
