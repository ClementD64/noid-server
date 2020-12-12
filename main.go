package main

import (
	"fmt"
	"os"

	"github.com/ClementD64/noid-server/noid"
)

func main() {
	root := "."

	if len(os.Args) > 2 {
		fmt.Println("Usage: noid-server PATH")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		root = os.Args[1]
	}

	noid.New(root)
}
