package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
)

// JSONOut flag, which defines, that output must be printed as JSON string.
var JSONOut bool

func init() {
	flag.BoolVar(&JSONOut, "json", false, "Output in JSON.")
	flag.Parse()
}

func main() {
	appInfo := make(map[string]bool)
	for _, app := range []string{"wget", "git", "go"} {
		_, err := exec.LookPath(app)
		appInfo[app] = err == nil
	}

	if JSONOut {
		b, _ := json.Marshal(appInfo)
		fmt.Printf("%s\n", b)
	} else {
		for app, exist := range appInfo {
			fmt.Println(app, exist)
		}
	}
}
