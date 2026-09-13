package main

import (
	"fmt"
	"net"
)

const (
	address = "127.0.0.1:42069"
	network = "tcp"
)

func main() {
	con, err := net.Dial(network, address)
	defer func() { _ = con.Close() }()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	_, err = con.Write([]byte("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
