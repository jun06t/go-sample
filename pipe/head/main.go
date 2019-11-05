package main

import (
	"errors"
	"fmt"
	"syscall"
)

func main() {
	for {
		_, err := fmt.Println("Wow!")
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				break
			} else {
				panic(err)
			}
		}
	}
}
