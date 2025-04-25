package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckForEveryKeyword(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"For", " ", "every"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
