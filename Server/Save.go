package main

import (
	"errors"
	"strconv"
	"strings"
)

type Command struct {
	Type  string
	Key   string
	Value string
}

var store = make(map[string]string)

func ParseRESP(resp string) ([]string, error) {
	var tokens []string
	lines := strings.Split(resp, "\r\n")
	if len(lines) < 3 || lines[0][0] != '*' {
		return nil, errors.New("invalid RESP format")
	}

	numElements, err := strconv.Atoi(lines[0][1:])
	if err != nil || numElements <= 0 {
		return nil, errors.New("invalid number of elements in RESP")
	}

	for i := 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "$") {

			elementLength, err := strconv.Atoi(lines[i][1:])
			if err != nil || elementLength < 0 {
				return nil, errors.New("invalid element length in RESP")
			}

			i++
			if i >= len(lines) || len(lines[i]) != elementLength {
				return nil, errors.New("mismatched element length in RESP")
			}

			tokens = append(tokens, lines[i])
		}
	}

	if len(tokens) != numElements {
		return nil, errors.New("mismatched number of elements in RESP")
	}

	return tokens, nil
}

func handleCommand(tokens []string) (string, error) {
	if len(tokens) < 2 {
		return "", errors.New("invalid command, at least two tokens required")
	}

	command := tokens[0]
	var response string

	switch strings.ToUpper(command) {
	case "SET":
		if len(tokens) != 3 {
			return "", errors.New("SET command requires a key and a value")
		}
		key := tokens[1]
		value := tokens[2]
		store[key] = value
		response = "+OK"

	case "GET":
		if len(tokens) != 2 {
			return "", errors.New("GET command requires a key")
		}
		key := tokens[1]
		if val, exists := store[key]; exists {
			response = val
		} else {
			response = "nil"
		}
	case "DEL":
		if len(tokens) != 2 {
			return "", errors.New("Del command reqires a key")
		}
		key := tokens[1]
		delete(store, key)
		response = "+OK"
	default:
		return "", errors.New("unsupported command: " + command)
	}

	return response, nil
}
