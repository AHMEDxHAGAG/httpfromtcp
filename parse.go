package main

import (
	"errors"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer func() {
			_ = f.Close()
			close(ch)
		}()
		chulk, curr := make([]byte, 8), ""
		for {
			n, err := f.Read(chulk)
			if err != nil && !errors.Is(err, io.EOF) {
				return
			}
			if n > 0 {
				buffers := strings.Split(string(chulk[:n]), "\n")
				curr += buffers[0]
				for i := 1; i < len(buffers); i++ {
					ch <- curr
					curr = buffers[i]
				}
			}
			if errors.Is(err, io.EOF) {
				if curr != "" {
					ch <- curr
				}
				return
			}
		}
	}()
	return ch
}
