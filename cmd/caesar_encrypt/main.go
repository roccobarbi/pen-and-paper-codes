package main

import (
	"errors"
	"os"
	"strconv"
)

type config struct {
	IsPlainKeySet  bool
	PlainKey       string
	IsCypherKeySet bool
	CypherKey      string
	IsOffsetSet    bool
	Offset         byte
	Input          string
}

func (c *config) init() {
	c.IsPlainKeySet = false
	c.PlainKey = ""
	c.IsCypherKeySet = false
	c.CypherKey = ""
	c.IsOffsetSet = false
	c.Offset = 0
	c.Input = ""
}

/*
ingestArgs transforms an array of arguments into a slice of string slices. Each element of the external slice is made of
a first element that represents the flag's name (without leading -), followed by any additional string found before the
next argument.

Long-form "--" arguments or trailing elements that are not preceded by a flag are not supported (neither are they
needed) at this time.
*/
func ingestArgs(args []string) [][]string {
	var parsedArgs [][]string
	var argGroup []string
	for _, arg := range args {
		if arg[0] == '-' {
			// This argument defines a new option
			if len(argGroup) > 0 {
				// Some argument was being parsed already, it can be closed
				parsedArgs = append(parsedArgs, argGroup)
				argGroup = []string{}
			}
			argGroup = append(argGroup, arg[1:])
		} else {
			argGroup = append(argGroup, arg)
		}
	}
	parsedArgs = append(parsedArgs, argGroup) // append the last arg to be parsed
	return parsedArgs
}

func validateFlagF(configuration config, arg []string) (config, error) {
	return configuration, errors.New("Not implemented.")
}

func validateFlagC(configuration config, arg []string) (config, error) {
	return configuration, errors.New("Not implemented.")
}

func validateFlagP(configuration config, arg []string) (config, error) {
	return configuration, errors.New("Not implemented.")
}

func validateFlagM(configuration config, arg []string) (config, error) {
	return configuration, errors.New("Not implemented.")
}

func validateFlagO(configuration config, arg []string) (config, error) {
	return configuration, errors.New("Not implemented.")
}

/**
validateConfig parses the command line arguments and returns a configuration structure. The command line arguments are
not read directly, but rather received as validateConfig's only argument (a string slice). validateConfig returns a
tuple with the config structure and an error object (or nil).
*/
func validateConfig(args []string) (config, error) {
	var configuration config
	var err error
	configuration.init()
	parsedArgs := ingestArgs(args)
	for _, arg := range parsedArgs {
		if len(arg[0]) > 1 {
			// Multiple flags have been grouped together, which may be acceptable to other programs, but in this case is
			// an error (all arguments require an additional value).
			return configuration, errors.New("Grouped flags found, but they are not supported by this program: -" + arg[0])
		}
		switch arg[0] {
		case "f":
			configuration, err = validateFlagF(configuration, arg)
		case "c":
			configuration, err = validateFlagC(configuration, arg)
		case "p":
			configuration, err = validateFlagP(configuration, arg)
		case "m":
			configuration, err = validateFlagM(configuration, arg)
		case "o":
			configuration, err = validateFlagO(configuration, arg)
		}
		if err != nil {
			return configuration, err
		}
	}
	return configuration, nil
}

/*
Usage:
	caesar_encrypt {key} -t {text}
	caesar_encrypt {key} -f {text file}
*/
func main() {
	configuration, err := validateConfig(os.Args)
}
