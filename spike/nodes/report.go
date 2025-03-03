package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type ReportNode struct {
	Node
	Name string
	Author string

	Body []*Node
}

func ParseReportNode(ast *spike.AST) (*ReportNode, error) {
	report := &ReportNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_ReportHeader, token.TokenType_ReportHeader.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	nameTokens, err := ast.ConsumeTokenUntilMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	if err != nil {
		return nil, err
	}

	firstNameToken := nameTokens[0]
	lastNameToken := nameTokens[len(nameTokens)-1]
	report.Name = ast.GetSourceText(firstNameToken.Start, lastNameToken.Start + lastNameToken.Length - firstNameToken.Start)

	_, err = ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	for {
		if ast.CheckType(token.TokenType_ReportFooter) { break }
		if ast.CheckType(token.TokenType_EndOfFile) {
			return nil, ast.CreateErrorFromToken(ast.Peek(), token.TokenType_FunctionFooter.Message("Could not find %s"))
		}

		if ast.CheckType(token.TokenType_NewLine) {
			continue
		}

		ast.Next()
	}

	_, err = ast.ConsumeToken(token.TokenType_ReportFooter, token.TokenType_ReportFooter.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	authorToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	report.Author = authorToken.Value

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	if !ast.EndOfFile() {
		return nil, ast.CreateErrorFromToken(ast.Peek(), token.TokenType_EndOfFile.Message("Expected %s"))
	}

	report.Start = startToken.Start
	report.Length = endToken.Start + endToken.Length - startToken.Start

	return report, nil
}
