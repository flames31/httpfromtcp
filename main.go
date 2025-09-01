package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		fmt.Println("Connection has been accepted!")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Println("Connection has been closed!")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineCh := make(chan (string))
	go func() {
		defer close(lineCh)
		defer f.Close()
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
	}()
	return lineCh
}
