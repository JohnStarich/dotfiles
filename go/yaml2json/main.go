package main

import (
	"encoding/json"
	"io"
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
	decoder := yaml.NewDecoder(in)
	encoder := json.NewEncoder(os.Stdout)
	for {
		var value any
		err := decoder.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		err = encoder.Encode(value)
		if err != nil {
			panic(err)
		}
	}
}
