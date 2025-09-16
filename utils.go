package main

import (
	"fmt"
	"log"
	"os"
)

func readFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err);
	}
	return file;
}

func listen[T any](chn <- chan T) {
	for result := range chn {
		fmt.Printf("read: %v\n", result)
	}
}