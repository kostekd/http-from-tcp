package main

import (
	"fmt"
)

func listen[T any](chn <- chan T) {
	for result := range chn {
		fmt.Printf("%s\n", result)
	}
}