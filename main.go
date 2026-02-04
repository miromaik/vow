package main

import (
	"os"
	"github.com/miromaik/vow/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
