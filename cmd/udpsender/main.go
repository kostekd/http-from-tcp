package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)



func main() {
	local, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Print("Error: failed to accept connection\n")
	}

	conn, err := net.DialUDP("udp4", nil, local)
	if err != nil {
		fmt.Print("Error: failed to accept connection\n")
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print(">")
		text, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(text))

		if err != nil {
			fmt.Print("Error: failed sending a message\n")
		}

	}
}