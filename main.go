package main

import (
	"io"
	"log"
	"strings"
)

func getLinesChannel(file io.ReadCloser) <- chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)

		str := ""
		for {
			buf := make([]byte, 8)
			chunk, err := file.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
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
	file := readFile("messages.txt")
	channel := getLinesChannel(file)

	listen(channel)
}