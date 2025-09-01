package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file : %v", err)
		os.Exit(1)
	}
	defer file.Close()

	line := ""
	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			break
		}

		data = data[:n]
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			line += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", line)
			line = ""
		}
		line += string(data)
	}
	if len(line) != 0 {
		fmt.Printf("read: %s\n", line)
	}
}
