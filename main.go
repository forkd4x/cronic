package main

import (
	"os"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	cronic, err := NewCronic(root)
	if err != nil {
		panic(err)
	}

	if err := cronic.LoadJobs(); err != nil {
		panic(err)
	}

	cronic.Start()
	if err := cronic.Shutdown(); err != nil {
		panic(err)
	}
}
