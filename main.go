package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		return
	}
	chulk := make([]byte, 8)
	for {
		n, err := file.Read(chulk)
		if n > 0 {
			fmt.Printf("read: %s\n", chulk[:n])
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
}
