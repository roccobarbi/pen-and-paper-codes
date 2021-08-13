package main

import (
	"errors"
	"os"
	"strconv"
)

type config struct {
	Key   byte
	Mode  byte
	Input string
}

/*
Usage:
	caesar_encrypt {key} -t {text}
	caesar_encrypt {key} -f {text file}
*/
func main() {
	configuration, err := validateConfig(os.Args)
}

func validateConfig(args []string) (config, error) {
	var configuration config
	if len(args) != 4 {
		return configuration, errors.New("Wrong number of arguments: expected 3, got " + strconv.Itoa(len(args)))
	}
	if len(args[1]) != 1 {
		return configuration, errors.New("Wrong length of key: expected 1, got " + strconv.Itoa(len(args[1])))
	} else {
		configuration.Key = args[1][0]
	}
	if len(args[2]) != 2 || args[2][1] != '-' || (args[2][2] != 'f' && args[2][2] != 't') {
		return configuration, errors.New("Invalid mode or bad mode formatting: expected -f or -t, got " + strconv.Itoa(len(args[1])))
	} else {
		configuration.Mode = args[2][0]
	}
	configuration.Input = args[3]
	return configuration, nil
}
