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
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
