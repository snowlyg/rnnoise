package main

import (
	"fmt"
	"os"

	"github.com/snowlyg/rnnoise"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input noisy file> <output denoised file>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	rnnoise.PlayFile(inputFile)
}
