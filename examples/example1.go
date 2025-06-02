///usr/bin/env go run "$0" "$@"; exit "$?"
// cronic:
//   name: Example Go Job
//   desc: Say hello every 3 seconds
//   cron: */3 0 0 0 0 0

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
