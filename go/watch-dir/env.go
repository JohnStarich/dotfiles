package main

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseEnvFile(contents []byte) ([]string, error) {
	var env []string
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	for scanner.Scan() {
		line := scanner.Text()
		keyValue, ok, err := parseKeyValueAssignment(line)
		if err != nil {
			return nil, err
		}
		if ok {
			env = append(env, keyValue)
		}
	}
	return env, nil
}

func parseKeyValueAssignment(line string) (string, bool, error) {
	if line == "" || line[0] == '#' {
		return "", false, nil
	}
	equalIndex := strings.IndexRune(line, '=')
	if equalIndex == -1 {
		return "", false, errors.New("no '=' sign found")
	}
	value := line[equalIndex+1:]
	if value == "" {
		return line, false, nil
	}
	var remaining string
	if value[0] == '"' {
		endQuoteIndex := findEndQuote(value[1:])
		if endQuoteIndex == -1 {
			return "", false, errors.Errorf("found start quote, but no end quote: %q", line)
		}
		endIndex := endQuoteIndex + 2 // 1 for start offset + 1 to include end quote
		var err error
		value, err = strconv.Unquote(value[:endIndex])
		if err != nil {
			return "", false, err
		}
		remaining = value[endIndex:]
	}
	if remaining != "" {
		remaining = strings.TrimSpace(remaining)
		if remaining != "" && remaining[0] != '#' {
			return "", false, errors.Errorf("values must not contain spaces and in-line comments must be prefixed with a '#': %q", remaining)
		}
	}
	return line, true, nil
}

func findEndQuote(str string) int {
	escaped := false
	for index, r := range str {
		if r == '"' && !escaped {
			return index
		}
		if escaped {
			escaped = false
		} else if r == '\\' {
			escaped = true
		}
	}
	return -1
}
