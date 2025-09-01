package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file : %v", err)
		os.Exit(1)
	}

	bytes := make([]byte, 8)
	for {
		n, err := file.Read(bytes)
		if err != nil {
			break
		}
		fmt.Printf("read: %s\n", string(bytes[:n]))
	}
}
