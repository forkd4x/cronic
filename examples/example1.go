///usr/bin/env go run "$0" "$@"; exit "$?"
//go:build ignore

// cronic:
//   name: Example Go Job
//	 desc: Say hello every 3 seconds
//	 cron: */3 * * * * *

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
	time.Sleep(time.Second)
	fmt.Printf("Hello, again, from %s using %s\n", filepath.Base(filename), goPath)
}
