package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func readFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err);
	}
	return file;
}

func getLinesChannel(file io.ReadCloser) <- chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		buf := make([]byte, 8)
		str := ""
		for {
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

	for result := range channel {
		fmt.Printf("read: %s", result);
	}
}