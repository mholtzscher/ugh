package main

import (
	"fmt"
	"os"

	"github.com/mholtzscher/ugh/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
