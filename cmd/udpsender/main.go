package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	con, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer func() { _ = con.Close() }()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		_, err = con.Write(line)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}
}
