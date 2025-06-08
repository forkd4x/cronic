///usr/bin/env go run "$0" "$@"; exit "$?"
//go:build ignore

package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	goPath, _ := exec.LookPath("go")
	fmt.Printf("Hello, from %s using %s\n", filepath.Base(filename), goPath)
}
