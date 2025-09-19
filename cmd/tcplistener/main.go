package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)


func listen[T any](chn <- chan T) {
	for result := range chn {
		fmt.Printf("%s\n", result)
	}
}

func getLinesChannel(connection net.Conn) <- chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)

		str := ""
		for {
			buf := make([]byte, 8)
			chunk, err := connection.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("dev:kostekd Error: %v\n", err)
			}
	
			if chunk == 0 {
				break
			}
	
			index := strings.Index(string(buf[:8]), "\n");
			if index == -1 {
				str += string(buf[:8]);
			} else {
				str += string(buf[:index + 1])
				lines <- str
				str = string(buf[index + 1:])
			}
		}
	}()
	return lines
}


func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Print("Failed to start TCP listener on :42069\n")
	}
	defer l.Close();
	
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Print("Error: failed to accept connection\n")
		}
		fmt.Print("connection accepted\n")
		channel := getLinesChannel(connection)
		listen(channel)
		connection.Close()
		fmt.Print("connection closed\n")
	}
}