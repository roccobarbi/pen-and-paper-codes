package main

import (
	"errors"
	"os"
	"pen-and-paper-codes/utils"
	"strconv"
	"strings"
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

func (c *config) setPlainKey(key string) {
	c.IsPlainKeySet = true
	c.PlainKey = key
}

func (c *config) setCypherKey(key string) {
	c.IsCypherKeySet = true
	c.CypherKey = key
}

func (c *config) setOffset(offset byte) {
	c.IsOffsetSet = true
	c.Offset = offset
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
	if len(arg) != 2 {
		return configuration, errors.New("-f option used without a filename")
	}
	configuration.Input = arg[2]
	return configuration, nil
}

func validateFlagC(configuration config, arg []string) (config, error) {
	if len(arg) != 2 {
		return configuration, errors.New("-c option used without a key")
	}
	if !utils.IsAlphabeticString(arg[2]) {
		return configuration, errors.New("-c key uses non-alphabetic characters")
	}
	configuration.setCypherKey(strings.ToLower(arg[2])) // all lowercase
	return configuration, nil
}

func validateFlagP(configuration config, arg []string) (config, error) {
	if len(arg) != 2 {
		return configuration, errors.New("-p option used without a key")
	}
	if !utils.IsAlphabeticString(arg[2]) {
		return configuration, errors.New("-p key uses non-alphabetic characters")
	}
	configuration.setPlainKey(strings.ToLower(arg[2])) // all lowercase
	return configuration, nil
}

func validateFlagO(configuration config, arg []string) (config, error) {
	if len(arg) != 2 {
		return configuration, errors.New("-o option used without an offset")
	}
	if utils.IsAlphabeticString(arg[2]) {
		if len(arg[2]) != 1 {
			return configuration, errors.New("invalid offset: if a character is used, it must be exactly one character")
		}
		configuration.setOffset(strings.ToLower(arg[2])[0] - 'a') // the offset is a number between 0 and 25
	} else if utils.IsIntegerString(arg[2]) {
		offset, err := strconv.Atoi(arg[2])
		if err != nil {
			return configuration, errors.New("could not convert integrer offset to int")
		}
		if offset < 0 || offset > 25 {
			return configuration, errors.New("integer offset out of range, must be between 0 and 25")
		}
		configuration.setOffset(byte(offset))
	} else {
		return configuration, errors.New("invalid offset: it must be either a character or an integer between 0 and 25")
	}
	return configuration, errors.New("not implemented")
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
		case "o":
			configuration, err = validateFlagO(configuration, arg)
		default:
			err = errors.New("Unknown option: " + arg[0])
		}
		if err != nil {
			return configuration, err
		}
	}
	return configuration, nil
}

func setPlainAlphabet(configuration config) []rune {
	var plainAlphabet []rune
	var used map[rune]bool
	var alphabet = "abcdefghijklmnopqrstuvwxyz"
	if configuration.IsPlainKeySet {
		key := utils.StripDuplicateCharacters(configuration.PlainKey)
		for _, letter := range key {
			if _, inUse := used[letter]; !inUse {
				plainAlphabet = append(plainAlphabet, letter)
				used[letter] = true
			}
		}
	}
	for _, letter := range alphabet {
		if _, inUse := used[letter]; !inUse {
			plainAlphabet = append(plainAlphabet, letter)
			used[letter] = true
		}
	}
	return plainAlphabet
}

func setCypherAlphabet(configuration config) []rune {
	var cypherAlphabet []rune
	var used map[rune]bool
	var alphabet = "abcdefghijklmnopqrstuvwxyz"
	if configuration.IsOffsetSet {
		// start with the overflowed letter
		// add all other letters
	}
	if configuration.IsCypherKeySet {
		key := utils.StripDuplicateCharacters(configuration.CypherKey)
		for _, letter := range key {
			if _, inUse := used[letter]; !inUse {
				cypherAlphabet = append(cypherAlphabet, letter)
				used[letter] = true
			}
		}
	}
	for _, letter := range alphabet {
		if _, inUse := used[letter]; !inUse {
			cypherAlphabet = append(cypherAlphabet, letter)
			used[letter] = true
		}
	}
	return cypherAlphabet
}

/*
createCypher builds a map between the plaintext alphabet and the cyphertext's alphabet. Since a caesar's cypher is a
simple substitution cypher, that's all we need to encrypt the plaintext.
This version of the caesar cypher incorporates the (optional) features of aristocrats used by the American Cryptogram
Association, together with an (optional) offset that works like the original caesar cypher. At least one encryption
method must be set or an error will be returned:
- offset, the cyphertext is shifted by n characters to the right (e.g. with offset 3 the plaintext A is encrypted as D);
- plainKey, the key is written (if there are duplicate letters, only their first instance is kept, so that e.g. the key
  "sassy" is rendered as "say"), then the rest of the plaintext alphabet is written;
- cypherkey, the key is written (if there are duplicate letters, only their first instance is kept, so that e.g. the key
  "sassy" is rendered as "say"), then the rest of the cyphertext alphabet is written.
The 2 keys may or may not be equal, and they may or may not be present. If both a cypherkey and an offset are present,
the offset will be applied after writing the key, to the remaining letters of the cypher's alphabet.
If the combination of keys/offset causes a plaintext letter to be represented by itself in the cyphertext, an error is
returned as the code would be much easier to break.
*/
func createCypher(configuration config) (map[rune]rune, error) {
	var cypher map[rune]rune
	if !configuration.IsPlainKeySet && !configuration.IsCypherKeySet && !configuration.IsOffsetSet {
		return nil, errors.New("no key or offset has been set")
	}
	plainAlphabet := setPlainAlphabet(configuration)
	cypherAlphabet := setCypherAlphabet(configuration)
	for i, l := range plainAlphabet {
		cypher[l] = cypherAlphabet[i]
	}
	return cypher, nil
}

/*
encrypt encrypts the plaintext, displaying it on screen and saving it in the plaintext's location using the .caesar
file extension
*/
func encrypt(configuration config, cypher map[rune]rune) {
}

/*
Usage:
	caesar_encrypt {key} -t {text}
	caesar_encrypt {key} -f {text file}
*/
func main() {
	configuration, err := validateConfig(os.Args)
	utils.ExitIfError(err)
	cypher, err := createCypher(configuration)
	utils.ExitIfError(err)
	encrypt(configuration, cypher)
}
