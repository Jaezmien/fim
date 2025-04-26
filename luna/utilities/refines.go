package utilities

import "strings"

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
		case "r":
			sb.WriteByte('\r')
		case "n":
			sb.WriteByte('\n')
		case "t":
			sb.WriteByte('\t')
		case "\"":
			sb.WriteByte('"')
		case "\\":
			sb.WriteByte('\\')
		default:
			sb.WriteByte(value[idx])
			sb.WriteByte(value[idx+1])
		}
		idx++
	}

	return sb.String()
}
