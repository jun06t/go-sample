package main

import (
	"fmt"
	"runtime/debug"

	library "github.com/jun06t/go-sample/toolchain-library"
)

func main() {
	info, ok := debug.ReadBuildInfo()
	if ok && info != nil {
		fmt.Println("Go Version:", info.GoVersion)
	} else {
		fmt.Println("Unable to determine Go version")
	}
	fmt.Println(library.Echo())
}
