package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	flag.Parse()
}

func main() {
	// read
	b, err := readFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// transform
	out := strings.ToUpper(string(b))

	// write
	fmt.Print(out)
}

func readFile() ([]byte, error) {
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	var r io.Reader
	switch filename {
	case "":
		if terminal.IsTerminal(int(syscall.Stdin)) {
			return nil, errors.New("usage: capitalize path")
		}
		r = os.Stdin
	case "-":
		r = os.Stdin
	default:
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return b, nil
}
