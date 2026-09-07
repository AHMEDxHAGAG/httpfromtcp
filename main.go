package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		return
	}
	ch := getLinesChannel(file)
	for str := range ch {
		fmt.Printf("read: %s\n", str)
	}
}
