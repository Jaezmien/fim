package utilities

import (
	"strings"
)

type ErrorPair struct {
	Line   int
	Column int
}

func GetErrorIndexPair(source string, index int) *ErrorPair {
	content := source[0 : index+1]
	lines := strings.Split(content, "\n")

	return &ErrorPair{
		Line:   len(lines),
		Column: len(lines[len(lines)-1]),
	}
}

func UnsanitizeString(value string, trim bool) string {
	sb := strings.Builder{}

	start := 0
	amount := len(value)
	if trim {
		start += 1
		amount -= 1
	}

	for idx := start; idx < amount; idx++ {
		if value[idx] != '\\' || idx+1 >= amount {
			sb.WriteByte(value[idx])
			continue
		}

		nextChar := string(value[idx+1])
		switch nextChar {
		case "0":
			sb.WriteByte(byte(rune(0)))
			break
		case "r":
			sb.WriteByte('\r')
			break
		case "n":
			sb.WriteByte('\n')
			break
		case "t":
			sb.WriteByte('\t')
			break
		case "\"":
			sb.WriteByte('"')
			break
		default:
			sb.WriteByte(value[idx])
			sb.WriteByte(value[idx+1])
			break
		}
		idx++
	}

	return sb.String()
}
