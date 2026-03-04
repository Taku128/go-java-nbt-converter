package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	javanbt "github.com/ntaku256/go-java-nbt-converter"
)

func main() {
	outputFile := flag.String("o", "", "Output NBT file path (defaults to <input>.nbt)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: convert [options] <input_file>\n")
		fmt.Fprintf(os.Stderr, "Supported formats: .litematic, .schem, .nbt (Java Structure)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputFile := args[0]
	if *outputFile == "" {
		ext := filepath.Ext(inputFile)
		*outputFile = inputFile[:len(inputFile)-len(ext)] + ".converted.nbt"
	}

	fmt.Printf("Converting: %s → %s\n", inputFile, *outputFile)

	nbtBytes, err := javanbt.ConvertAny(context.Background(), inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, nbtBytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Wrote %d bytes to %s\n", len(nbtBytes), *outputFile)
}
