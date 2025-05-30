package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckReportHeader(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"Dear", " ", "Princess", " ", "Celestia", ":"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
func CheckReportFooter(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"Your", " ", "faithful", " ", "student", ","}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
