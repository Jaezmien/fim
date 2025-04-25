package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckConstantKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "always" {
		return 0
	}

	return 1
}

func CheckAndKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "and" {
		return 0
	}

	return 1
}

func CheckOrKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "or" {
		return 0
	}

	return 1
}

func CheckOfKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "of" {
		return 0
	}

	return 1
}

func CheckInKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "in" {
		return 0
	}

	return 1
}

func CheckFromKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "from" {
		return 0
	}

	return 1
}

func CheckToKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "to" {
		return 0
	}

	return 1
}

func CheckStatementEndKeyword(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"That", "'", "s", " ", "what", " ", "I", " ", "did"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
