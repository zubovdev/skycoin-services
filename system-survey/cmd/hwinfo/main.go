package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jaypipes/ghw"
)

// JSONOut flag, which defines, that output must be printed as JSON string.
var JSONOut bool

func init() {
	flag.BoolVar(&JSONOut, "json", false, "Output in JSON.")
	flag.Parse()
}

func main() {
	info, _ := ghw.Host()

	if JSONOut {
		b, _ := json.Marshal(info)
		fmt.Printf("%s", b)
	} else {
		fmt.Println(info)
	}
}
