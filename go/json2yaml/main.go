package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	in := os.Stdin
	if len(os.Args) >= 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		in = file
	}
	var value any
	err := json.NewDecoder(in).Decode(&value)
	if err != nil {
		panic(err)
	}
	err = yaml.NewEncoder(os.Stdout).Encode(value)
	if err != nil {
		panic(err)
	}
}
