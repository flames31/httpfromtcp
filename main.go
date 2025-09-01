package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file : %v", err)
		os.Exit(1)
	}
	defer file.Close()

	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineCh := make(chan (string))
	go func() {
		line := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				line += string(data[:i])
				data = data[i+1:]
				lineCh <- line
				line = ""
			}
			line += string(data)
		}
		if len(line) != 0 {
			lineCh <- line
		}

		close(lineCh)
	}()
	return lineCh
}
