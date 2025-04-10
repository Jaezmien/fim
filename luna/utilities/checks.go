package utilities

import (
	"slices"
	"strings"
)

var	booleanTrueStrings = [...]string{"yes", "true", "right", "correct"}
var	booleanFalseStrings = [...]string{"no", "false", "wrong", "incorrect"}

func AsBooleanValue(str string) (bool, bool) {
	if slices.Contains(booleanTrueStrings[:], str) {
		return true, true
	}
	if slices.Contains(booleanFalseStrings[:], str) {
		return false, true
	}

	return false, false
}

func AsCharacterValue(str string) (string, bool) {
	if len(str) == 1 {
		return str, true
	}

	if len(str) == 2 {
		if !strings.HasPrefix(str, "\\") {
			return "", false
		}

		switch str[1] {
			case '0':
				return string(byte(0)), true
			case 'r':
				return "\r", true
			case 'n':
				return "\n", true
			case 't':
				return "\t", true
			default:
				return string(str[1]), true
		}
	}

	return "", false
}
