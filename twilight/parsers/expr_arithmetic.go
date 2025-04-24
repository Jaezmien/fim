package parsers

import (
	"slices"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckInfixAddition(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"plus"}
	if slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 1
	}

	ExpectedTokens := []string{"added", " ", "to"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}

func CheckPrefixAddition(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"add"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func CheckInfixSubtraction(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"minus", "without"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func CheckPrefixSubtraction(tokens *queue.Queue[*token.Token]) int {
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

func CheckInfixMultiplication(tokens *queue.Queue[*token.Token]) int {
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

func CheckPrefixMultiplication(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"multiply"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func CheckInfixDivision(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"divided", " ", "by"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}

func CheckPrefixDivision(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"divide"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func CheckInfixModulo(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"mod", "modulo", "remainder"}

	if slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 1
	}

	return 0
}

func CheckPrefixModulo(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"remainder", " ", "of"}
	if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return len(ExpectedTokens)
	}

	return 0
}
