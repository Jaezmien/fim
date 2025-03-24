package parsers

import (
	"slices"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func IsNumberType(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() >= 1 {
		ExpectedSingleTokens := []string{"number"}
		if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
			return 1
		}
	}

	if tokens.Len() >= 3 {
		ExpectedMultiTokens := [][]string{
			{"a", " ", "number"},
			{"the", " ", "number"},
		}
		for _, sequence := range ExpectedMultiTokens {
			if utilities.CheckTokenSequence(tokens, sequence) {
				return len(sequence)
			}
		}
	}

	return 0
}

func IsNumberArrayType(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() >= 1 {
		ExpectedSingleTokens := []string{"numbers"}
		if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
			return 1
		}
	}

	if tokens.Len() >= 3 {
		ExpectedMultiTokens := [][]string{
			{"the", " ", "numbers"},
			{"many", " ", "numbers"},
		}
		for _, sequence := range ExpectedMultiTokens {
			if utilities.CheckTokenSequence(tokens, sequence) {
				return len(sequence)
			}
		}
	}

	return 0
}

func IsBooleanType(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() >= 1 {
		ExpectedSingleTokens := []string{"argument", "logic"}
		if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
			return 1
		}
	}

	if tokens.Len() >= 3 {
		ExpectedMultiTokens := [][]string{
			{"an", " ", "argument"},
			{"the", " ", "argument"},
			{"the", " ", "logic"},
		}
		for _, sequence := range ExpectedMultiTokens {
			if utilities.CheckTokenSequence(tokens, sequence) {
				return len(sequence)
			}
		}
	}

	return 0
}

func IsBooleanArrayType(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() >= 1 {
		ExpectedSingleTokens := []string{"arguments", "logics"}
		if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
			return 1
		}
	}

	if tokens.Len() >= 3 {
		ExpectedMultiTokens := [][]string{
			{"many", " ", "arguments"},
			{"many", " ", "logics"},
			{"the", " ", "arguments"},
			{"the", " ", "logics"},
		}
		for _, sequence := range ExpectedMultiTokens {
			if utilities.CheckTokenSequence(tokens, sequence) {
				return len(sequence)
			}
		}
	}

	return 0
}

func IsCharacterType(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() >= 1 {
		ExpectedSingleTokens := []string{"character", "letter"}
		if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
			return 1
		}
	}

	if tokens.Len() >= 3 {
		ExpectedMultiTokens := [][]string{
			{"a", " ", "character"},
			{"a", " ", "letter"},
			{"the", " ", "character"},
			{"the", " ", "letter"},
		}
		for _, sequence := range ExpectedMultiTokens {
			if utilities.CheckTokenSequence(tokens, sequence) {
				return len(sequence)
			}
		}
	}

	return 0
}

func IsStringType(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() >= 1 {
		ExpectedSingleTokens := []string{"characters", "letters", "phrase", "quote", "sentence", "word"}
		if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
			return 1
		}
	}

	if tokens.Len() >= 3 {
		ExpectedMultiTokens := [][]string{
			{"a", " ", "phrase"},
			{"a", " ", "quote"},
			{"a", " ", "sentence"},
			{"a", " ", "word"},
			{"the", " ", "characters"},
			{"the", " ", "letters"},
			{"the", " ", "phrase"},
			{"the", " ", "quote"},
			{"the", " ", "sentence"},
			{"the", " ", "word"},
		}
		for _, sequence := range ExpectedMultiTokens {
			if utilities.CheckTokenSequence(tokens, sequence) {
				return len(sequence)
			}
		}
	}

	return 0
}

func IsStringArrayType(tokens *queue.Queue[*token.Token]) int {
	if tokens.Len() >= 1 {
		ExpectedSingleTokens := []string{"phrases", "quotes", "sentences", "words"}
		if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
			return 1
		}
	}

	if tokens.Len() >= 3 {
		ExpectedMultiTokens := [][]string{
			{"many", " ", "phrases"},
			{"many", " ", "quotes"},
			{"many", " ", "sentences"},
			{"many", " ", "words"},
			{"the", " ", "phrases"},
			{"the", " ", "quotes"},
			{"the", " ", "sentences"},
			{"the", " ", "words"},
		}
		for _, sequence := range ExpectedMultiTokens {
			if utilities.CheckTokenSequence(tokens, sequence) {
				return len(sequence)
			}
		}
	}

	return 0
}
