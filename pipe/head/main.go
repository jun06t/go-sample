package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGPIPE)

	go func() {
		select {
		case <-signal_chan:
			os.Exit(1)
		}
	}()

	run()
}

func run() {
	for {
		_, err := fmt.Println("run!")
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				fmt.Fprintln(os.Stderr, "catch broken pipe error")
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
