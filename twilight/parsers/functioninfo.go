package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckFunctionHeaderMain(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"Today", " ", "I", " ", "learned"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}

func CheckFunctionHeader(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"I", " ", "learned"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}

func CheckFunctionFooter(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"That", "'", "s", " ", "all", " ", "about"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
func CheckFunctionParameter(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "using" {
		return 0
	}

	return 1
}
func CheckFunctionReturn(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value == "with" {
		return 1
	}

	ExpectedTokens := []string{"to", " ", "get"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}

func CheckReturnKeyword(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"Then", " ", "you", " ", "get"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}

func CheckFunctionCallMethod(tokens *queue.Queue[*token.Token]) int {
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
