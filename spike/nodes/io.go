package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type PrintNode struct {
	Node

	NewLine bool
	Value   INode
}

func (p *PrintNode) Type() NodeType {
	return TYPE_PRINT
}

func ParsePrintNode(ast *ast.AST) (*PrintNode, error) {
	printNode := &PrintNode{}

	startToken := ast.Peek()
	if startToken.Type != token.TokenType_Print && startToken.Type != token.TokenType_PrintNewline {
		return nil, ast.CreateErrorFromToken(startToken, "Expected print token")
	}

	printNode.NewLine = startToken.Type == token.TokenType_PrintNewline
	ast.Next()

	valueTokens, err := ast.ConsumeTokenUntilMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	if err != nil {
		return nil, err
	}

	printNode.Value = CreateValueNode(valueTokens, CreateValueNodeOptions{})

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	printNode.Start = startToken.Start
	printNode.Length = endToken.Start + endToken.Length - startToken.Start

	return printNode, nil
}
