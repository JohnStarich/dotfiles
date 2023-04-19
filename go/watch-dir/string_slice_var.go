package main

import "strings"

type stringSliceVar []string

func (s *stringSliceVar) Get() interface{} {
	return []string(*s)
}

func (s *stringSliceVar) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *stringSliceVar) String() string {
	return strings.Join(*s, ", ")
}
