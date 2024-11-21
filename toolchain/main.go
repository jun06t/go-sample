package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	info, ok := debug.ReadBuildInfo()
	if ok && info != nil {
		fmt.Println("Go Version:", info.GoVersion)
	} else {
		fmt.Println("Unable to determine Go version")
	}
}
