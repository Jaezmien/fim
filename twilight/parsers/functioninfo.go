package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/twilight/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func IsFunctionHeaderMain(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"Today", " ", "I", " ", "learned"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}

func IsFunctionHeader(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"I", " ", "learned"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}

func IsFunctionFooter(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"That", "'", "s", " ", "all", " ", "about"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
func IsFunctionParameter(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "using" {
		return 0
	}

	return 1
}
func IsFunctionReturn(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value == "with" {
		return 1
	}

	ExpectedTokens := []string{"to", " ", "get"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}

func IsReturnKeyword(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"Then", " ", "you", " ", "get"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}

func IsFunctionCallMethod(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "I" {
		return 0
	}
	if tokens.Len() < 3 {
		return 0
	}

	ExpectedTokens := [][]string{
		{"I", " ", "remembered"},
		{"I", " ", "would"},
	}
	for _, sequence := range ExpectedTokens {
		if utilities.CheckTokenSequence(tokens, sequence) {
			return len(sequence)
		}
	}

	return 0
}
