package cmd

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	Type  string
	Key   string
	Value string
}

func encode(tokens []string) (string, error) {

	if len(tokens) < 2 {
		return "", errors.New("invalid command, at least two tokens required")
	}

	command := &Command{
		Type: strings.ToUpper(tokens[0]),
	}

	switch command.Type {
	case "SET":
		if len(tokens) != 3 {
			return "", errors.New("SET command requires a key and a value")
		}
		command.Key = tokens[1]
		command.Value = tokens[2]
	case "GET":
		if len(tokens) != 2 {
			return "", errors.New("GET command requires a key")
		}
		command.Key = tokens[1]
	case "DEL":
		if len(tokens) != 2 {
			return "", errors.New("DEL commands requires only a key")
		}
		command.Key = tokens[1]
	case "INCR":
		if len(tokens) != 2 {
			return "", errors.New("INCR command only requires a key")
		}
		command.Key = tokens[1]
	case "DECR":
		if len(tokens) != 2 {
			return "", errors.New("DECR command only requires a key")
		}
		command.Key = tokens[1]
	default:
		return "", errors.New("unsupported command: " + tokens[0])
	}

	return encoder(command)
}

func encoder(command *Command) (string, error) {
	var out strings.Builder

	if command.Key == "" {
		return "", errors.New("key cannot be empty")
	}

	switch command.Type {
	case "SET":
		if command.Value == "" {
			return "", errors.New("SET command requires a value")
		}
		out.WriteString(fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(command.Key), command.Key, len(command.Value), command.Value))
	case "GET":
		out.WriteString(fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%s\r\n", len(command.Key), command.Key))
	case "DEL":
		out.WriteString(fmt.Sprintf("*2\r\n$3\r\nDEL\r\n$%d\r\n%s\r\n", len(command.Key), command.Key))
	case "INCR":
		out.WriteString(fmt.Sprintf("*2\r\n$4\r\nINCR\r\n$%d\r\n%s\r\n", len(command.Key), command.Key))
	case "DECR":
		out.WriteString(fmt.Sprintf("*2\r\n$4\r\nDECR\r\n$%d\r\n%s\r\n", len(command.Key), command.Key))
	default:
		return "", errors.New("unsupported command: " + command.Type)
	}

	return out.String(), nil
}
