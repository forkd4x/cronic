///usr/bin/env go run "$0" "$@"; exit "$?"
//go:build ignore

// cronic:
//   name: Example Go Job
//	 desc: Say hello every 4 seconds
//	 cron: */4 * * * * *

package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	goPath, _ := exec.LookPath("go")
	fmt.Printf("Hello, from %s using %s\n", filepath.Base(filename), goPath)
	time.Sleep(2 * time.Second)
	fmt.Printf("Bye, from %s using %s\n", filepath.Base(filename), goPath)
}
