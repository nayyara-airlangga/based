package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/nayyara-airlangga/basedlang/repl"
)

func main() {
	fmt.Printf("Basedlang v0.0.1 on %s %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Type away!")
	repl.Start(os.Stdin, os.Stdout)
}
