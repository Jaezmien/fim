package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type ReportNode struct {
	Node
	Title  string
	Author string

	Body []DynamicNode
}

func (r *ReportNode) Type() NodeType {
	return TYPE_REPORT
}

func ParseReportNode(ast *ast.AST) (*ReportNode, error) {
	report := &ReportNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_ReportHeader, token.TokenType_ReportHeader.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	nameTokens, err := ast.ConsumeUntilTokenMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	if err != nil {
		return nil, err
	}

	firstNameToken := nameTokens[0]
	lastNameToken := nameTokens[len(nameTokens)-1]
	report.Title = ast.GetSourceText(firstNameToken.Start, lastNameToken.Start+lastNameToken.Length-firstNameToken.Start)

	_, err = ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	for {
		if ast.CheckType(token.TokenType_ReportFooter) {
			break
		}
		if ast.CheckType(token.TokenType_EndOfFile) {
			return nil, ast.Peek().CreateError(token.TokenType_FunctionFooter.Message("Could not find %s"), ast.Source)
		}

		if ast.CheckType(token.TokenType_NewLine) {
			continue
		}

		if ast.CheckType(token.TokenType_FunctionMain) || ast.CheckType(token.TokenType_FunctionHeader) {
			functionNode, err := ParseFunctionNode(ast)

			if err != nil {
				return nil, err
			}

			report.Body = append(report.Body, functionNode)

			continue
		}

		if ast.CheckType(token.TokenType_Declaration) {
			declarationNode, err := ParseVariableDeclarationNode(ast)

			if err != nil {
				return nil, err
			}

			report.Body = append(report.Body, declarationNode)

			continue
		}

		return nil, ast.Peek().CreateError(ast.Peek().Type.Message("Unxpected token: %s"), ast.Source)
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
		return nil, ast.Peek().CreateError(token.TokenType_EndOfFile.Message("Expected %s"), ast.Source)
	}

	report.Start = startToken.Start
	report.Length = endToken.Start + endToken.Length - startToken.Start

	return report, nil
}
