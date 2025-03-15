package parsers

import (
	"slices"

	"git.jaezmien.com/Jaezmien/fim/twilight/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func IsLessThanEqualOperator(tokens *queue.Queue[*token.Token]) int {
	MultiTokens := [][]string{
		{"had", " ", "no", " ", "more", " ", "than"},
		{"has", " ", "no", " ", "more", " ", "than"},
		{"is", " ", "no", " ", "greater", " ", "than"},
		{"is", " ", "no", " ", "more", " ", "than"},
		{"is", " ", "not", " ", "greater", " ", "than"},
		{"is", " ", "not", " ", "more", " ", "than"},
		{"isn", "'", "t", " ", "greater", " ", "than"},
		{"isn", "'", "t", " ", "more", " ", "than"},
		{"was", " ", "no", " ", "greater", " ", "than"},
		{"was", " ", "no", " ", "more", " ", "than"},
		{"was", " ", "not", " ", "greater", " ", "than"},
		{"was", " ", "not", " ", "more", " ", "than"},
		{"wasn", "'", "t", " ", "greater", " ", "than"},
		{"wasn", "'", "t", " ", "more", " ", "than"},
		{"were", " ", "no", " ", "greater", " ", "than"},
		{"were", " ", "no", " ", "more", " ", "than"},
		{"were", " ", "not", " ", "greater", " ", "than"},
		{"were", " ", "not", " ", "more", " ", "than"},
		{"weren", "'", "t", " ", "greater", " ", "than"},
		{"weren", "'", "t", " ", "more", " ", "than"},
	}

	for _, ExpectedTokens := range MultiTokens {
		if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
			return len(ExpectedTokens)
		}
	}

	return 0
}

func IsGreaterThanEqualOperator(tokens *queue.Queue[*token.Token]) int {
	MultiTokens := [][]string{
		{"had", " ", "no", " ", "less", " ", "than"},
		{"has", " ", "no", " ", "less", " ", "than"},
		{"is", " ", "no", " ", "less", " ", "than"},
		{"is", " ", "not", " ", "less", " ", "than"},
		{"isn", "'", "t", " ", "less", " ", "than"},
		{"was", " ", "no", " ", "less", " ", "than"},
		{"was", " ", "not", " ", "less", " ", "than"},
		{"wasn", "'", "t", " ", "less", " ", "than"},
		{"were", " ", "no", " ", "less", " ", "than"},
		{"were", " ", "not", " ", "less", " ", "than"},
		{"weren", "'", "t", " ", "less", " ", "than"},
	}

	for _, ExpectedTokens := range MultiTokens {
		if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
			return len(ExpectedTokens)
		}
	}

	return 0
}

func IsGreaterThanOperator(tokens *queue.Queue[*token.Token]) int {
	MultiTokens := [][]string{
		{"had", " ", "more", " ", "than"},
		{"has", " ", "more", " ", "than"},
		{"were", " ", "more", " ", "than"},
		{"was", " ", "more", " ", "than"},

		{"is", " ", "greater", " ", "than"},
		{"was", " ", "greater", " ", "than"},
		{"were", " ", "greater", " ", "than"},
	}

	for _, ExpectedTokens := range MultiTokens {
		if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
			return len(ExpectedTokens)
		}
	}

	return 0
}

func IsLessThanOperator(tokens *queue.Queue[*token.Token]) int {
	MultiTokens := [][]string{
		{"had", " ", "less", " ", "than"},
		{"has", " ", "less", " ", "than"},
		{"is", " ", "less", " ", "than"},
		{"was", " ", "less", " ", "than"},
		{"were", " ", "less", " ", "than"},
	}

	for _, ExpectedTokens := range MultiTokens {
		if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
			return len(ExpectedTokens)
		}
	}

	return 0
}

func IsNotEqualOperator(tokens *queue.Queue[*token.Token]) int {
	MultiTokens := [][]string{
		{"wasn", "'", "t", " ", "equal", " ", "to"},
		{"isn", "'", "t", " ", "equal", " ", "to"},
		{"weren", "'", "t", " ", "equal", " ", "to"},
		{"had", "'", "t"},
		{"has", "'", "t"},
		{"isn", "'", "t"},
		{"wasn", "'", "t"},
		{"weren", "'", "t"},
		{"had", " ", "not"},
		{"has", " ", "not"},
		{"is", " ", "not"},
		{"was", " ", "not"},
		{"were", " ", "not"},
	}

	for _, ExpectedTokens := range MultiTokens {
		if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
			return len(ExpectedTokens)
		}
	}

	return 0
}

func IsEqualOperator(tokens *queue.Queue[*token.Token]) int {
	MultiTokens := [][]string{
		{"is", " ", "equal", " ", "to"},
		{"was", " ", "equal", " ", "to"},
		{"were", " ", "equal", " ", "to"},
	}

	for _, ExpectedTokens := range MultiTokens {
		if utilities.CheckTokenSequence(tokens, ExpectedTokens) {
			return len(ExpectedTokens)
		}
	}

	SingleTokens := []string{"is", "was", "were", "had", "has", "has", "likes", "like"}

	if slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 1
	}

	return 0
}
