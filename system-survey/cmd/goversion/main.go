package main

import (
	"flag"
	"os"
	"os/exec"
)

// JSONOut flag, which defines, that output must be printed as JSON string.
var JSONOut bool

func init() {
	flag.BoolVar(&JSONOut, "json", false, "Output in JSON.")
	flag.Parse()
}

func main() {
	// Lookup for command `go`.
	path, err := exec.LookPath("go")
	if err != nil {
		output{Error: err}.show()
		os.Exit(0)
	}

	// Execute `go version` command and get output.
	b, err := exec.Command("go", "version").Output()
	if err != nil {
		output{Error: err, Path: path}.show()
		os.Exit(0)
	}

	output{Error: nil, Path: path, Version: string(b)}.show()
}
