package main

import (
	"os"

	"github.com/gdanko/clearsky/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	// return
}
