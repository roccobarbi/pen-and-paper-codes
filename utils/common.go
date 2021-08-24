package utils

import "os"

func IsAlphabeticString(s string) bool {
	for _, c := range s {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return true
}

func ExitIfError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
