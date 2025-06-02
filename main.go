package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	err := os.Chdir("examples")
	if err != nil {
		panic(err)
	}

	dirEntries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		file, err := os.Open(dirEntry.Name())
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing file: %v\n", err)
			}
		}()

		// Read the first 10 kB of the file looking for `cronic:` yaml
		buffer := make([]byte, 10240)
		n, err := file.Read(buffer)
		if err != nil {
			panic(err)
		}
		if strings.Contains(string(buffer[:n]), "cronic:") {
			fmt.Println("Found cronic yaml in", dirEntry.Name())
		}
	}
}
