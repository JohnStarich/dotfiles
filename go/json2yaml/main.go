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
	decoder := json.NewDecoder(in)
	encoder := yaml.NewEncoder(os.Stdout)
	for decoder.More() {
		var value any
		err := decoder.Decode(&value)
		if err != nil {
			panic(err)
		}
		err = encoder.Encode(value)
		if err != nil {
			panic(err)
		}
	}
}
