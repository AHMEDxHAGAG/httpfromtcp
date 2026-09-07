package main

import (
	"fmt"
	"net"
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
		ch := getLinesChannel(con)
		for str := range ch {
			fmt.Println(str)
		}
		fmt.Println("Connection Closed")
	}
}
