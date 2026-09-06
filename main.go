package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()
	chulk, curr := make([]byte, 8), ""
	for {
		n, err := file.Read(chulk)
		if err != nil && !errors.Is(err, io.EOF) {
			return
		}
		if n > 0 {
			buffers := strings.Split(string(chulk[:n]), "\n")
			curr += buffers[0]
			for i := 1; i < len(buffers); i++ {
				fmt.Printf("read: %s\n", curr)
				curr = buffers[i]
			}
		}
		if errors.Is(err, io.EOF) {
			if curr != "" {
				fmt.Printf("read: %s\n", curr)
			}
			break
		}
	}
}
