package parsers

import (
	"slices"

	"git.jaezmien.com/Jaezmien/fim/twilight/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func IsInfixAddition(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"plus", "added"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func IsPrefixAddition(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"add"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func IsInfixSubtraction(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"minus", "without"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func IsPrefixSubtraction(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"subtract"}
	if slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 1
	}

	ExpectedTokens := []string{"the", " ", "difference", " ", "between"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}

func IsInfixMultiplication(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"times"}
	if slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 1
	}

	ExpectedTokens := []string{"multiplied", " ", "with"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}

func IsPrefixMultiplication(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"multiply"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func IsInfixDivision(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"divided", " ", "by"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}

func IsPrefixDivision(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"divide"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}
