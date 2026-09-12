package main

import (
	"fmt"
	"net"

	"github.com/AHMEDxHAGAG/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
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
			return
		}
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
		fmt.Println("Connection Closed")
	}
}
