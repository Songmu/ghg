package main

import (
	"os"

	"github.com/Songmu/ghg"
)

func main() {
	os.Exit((&ghg.CLI{ErrStream: os.Stderr, OutStream: os.Stdout}).Run(os.Args[1:]))
}
