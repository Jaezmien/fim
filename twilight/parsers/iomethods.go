package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/twilight/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	_ "git.jaezmien.com/Jaezmien/fim/twilight/utilities"
	"slices"
)

func IsPrintNewlineMethod(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() < 3 {
		return 0
	}
	if tokens.First().Value.Value != "I" {
		return 0
	}
	if tokens.Peek(1).Value.Value != " " {
		return 0
	}

	expectedTokens := []string{"said", "sang", "wrote"}
	if slices.Contains((expectedTokens), tokens.Peek(2).Value.Value) {
		return 3
	}

	return 0
}

func IsPrintMethod(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() < 5 {
		return 0
	}
	if tokens.First().Value.Value != "I" {
		return 0
	}
	if tokens.Peek(1).Value.Value != " " {
		return 0
	}
	if tokens.Peek(2).Value.Value != "quickly" {
		return 0
	}
	if tokens.Peek(3).Value.Value != " " {
		return 0
	}

	expectedTokens := []string{"said", "sang", "wrote"}
	if slices.Contains((expectedTokens), tokens.Peek(4).Value.Value) {
		return 5
	}

	return 0
}

func IsReadMethod(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() < 3 {
		return 0
	}
	if tokens.First().Value.Value != "I" {
		return 0
	}
	if tokens.Peek(1).Value.Value != " " {
		return 0
	}

	expectedTokens := []string{"heard", "read", "asked"}
	if slices.Contains((expectedTokens), tokens.Peek(2).Value.Value) {
		return 3
	}

	return 0
}
