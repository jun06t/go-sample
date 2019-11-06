package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"syscall"
	"time"
)

func main() {
	pr, pw := io.Pipe()
	go func() {
		for {
			buf, err := ioutil.ReadAll(pr)
			if err != nil {
				panic(err)
			}
			fmt.Println(buf)
			time.Sleep(10 * time.Millisecond)
		}
	}()
	fmt.Fprintf(pw, "hoge")
	pw.Close()
	_, err := fmt.Fprintf(pw, "fuga")
	if errors.Is(err, syscall.EPIPE) {
		fmt.Println("syscall error")
	}
	if errors.Is(err, io.ErrClosedPipe) {
		fmt.Println("close pipe error")
	}
	fmt.Printf("%v", err)
}
