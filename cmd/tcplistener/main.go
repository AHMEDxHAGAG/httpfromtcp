package main

import (
	"fmt"
	"io"
	"net"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/request"
)

const (
	address = "127.0.0.1:42069"
	network = "tcp"
)

func main() {
	listener, err := net.Listen(network, address)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer func() {
		_ = listener.Close()
	}()
	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		go perConnection(con)
	}
}

func perConnection(con io.ReadCloser) {
	defer func() {
		fmt.Println("Connection Closed")
		err := con.Close()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}()
	fmt.Println("Connection Started")
	req, err := request.RequestFromReader(con)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
		req.RequestLine.Method,
		req.RequestLine.RequestTarget,
		req.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for key, val := range req.Headers {
		fmt.Printf("- %s: %s\n", key, val)
	}
}
