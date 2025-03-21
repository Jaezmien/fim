package parsers

import (
	"strings"

	"git.jaezmien.com/Jaezmien/fim/twilight/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

func IsPostscript(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() < 5 {
		return 0
	}
	if tokens.First().Value.Value != "P" {
		return 0
	}
	if tokens.Peek(1).Value.Value != "." {
		return 0
	}

	for idx := 2; idx < tokens.Len(); idx += 2 {
		if idx + 1 >= tokens.Len() {
			return 0
		}

		if tokens.Peek(idx).Value.Value != "S" {
			return 0
		}
		if tokens.Peek(idx + 1).Value.Value != "." {
			return 0
		}
		if idx + 2 < tokens.Len() && strings.TrimSuffix(tokens.Peek(idx + 2).Value.Value, " ") == "" {
			return idx + 1	
		}
	}

	return 0
}
