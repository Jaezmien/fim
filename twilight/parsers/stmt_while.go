package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckWhileKeyword(tokens *queue.Queue[*token.Token]) int {
	ExpectedMultiTokens := [][]string{
		{"As", " ", "long", " ", "as"},
		{"Here", "'", "s", " ", "what", " ", "I", " ", "did", " ", "while"},
	}
	for _, sequence := range ExpectedMultiTokens {
		if utilities.CheckTokenSequence(tokens, sequence) {
			return len(sequence)
		}
	}

	return 0
}
