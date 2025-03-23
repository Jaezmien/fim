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
func (p *PrintNode) ToNode() Node {
	return Node{
		Start:  p.Start,
		Length: p.Length,
	}
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

	printNode.Value, err = CreateValueNode(valueTokens, CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	printNode.Start = startToken.Start
	printNode.Length = endToken.Start + endToken.Length - startToken.Start

	return printNode, nil
}

type PromptNode struct {
	Node

	Identifier string
	Prompt   INode
}

func (p *PromptNode) Type() NodeType {
	return TYPE_PROMPT
}
func (p *PromptNode) ToNode() Node {
	return Node{
		Start:  p.Start,
		Length: p.Length,
	}
}

func ParsePromptNode(ast *ast.AST) (*PromptNode, error) {
	node := &PromptNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_Prompt, token.TokenType_Prompt.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	identifier, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	node.Identifier = identifier.Value

	_, err = ast.ConsumeFunc(func(t *token.Token) bool { return t.Type == token.TokenType_Punctuation && t.Value == ":" }, "Expected ':'")
	if err != nil {
		return nil, err
	}
	
	valueTokens, err := ast.ConsumeTokenUntilMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	if err != nil {
		return nil, err
	}
	node.Prompt, err = CreateValueNode(valueTokens, CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	node.Start = startToken.Start
	node.Length = endToken.Start + endToken.Length - startToken.Start

	return node, nil
}
