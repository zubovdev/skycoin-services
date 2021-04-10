package main

import (
	"encoding/json"
	"fmt"
)

type output struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	Error   error  `json:"error"`
}

// show prints the output in format, defined by flag JSONOut.
func (o output) show() {
	if JSONOut {
		b, _ := json.Marshal(o)
		fmt.Printf("%s", b)
		return
	}

	fmt.Println("Path:", o.Path)
	fmt.Println("Error:", o.Error)
	fmt.Println("Version:", o.Version)
}
