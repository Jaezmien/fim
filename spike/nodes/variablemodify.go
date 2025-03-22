package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type VariableModifyNode struct {
	Node

	Identifier string

	Value     INode
}

func (d *VariableModifyNode) Type() NodeType {
	return TYPE_VARIABLE_MODIFY
}
func (f *VariableModifyNode) ToNode() Node {
	return Node{
		Start:  f.Start,
		Length: f.Length,
	}
}

func CheckVariableModifyNode(ast *ast.AST) (bool) {
	if ast.Peek() == nil {
		return false
	}	
	if ast.Peek().Type != token.TokenType_Identifier {
		return false
	}
	if ast.PeekAt(1) == nil {
		return false
	}	
	if ast.PeekAt(1).Type != token.TokenType_Identifier {
		return false
	}

	return true
}

func ParseVariableModifyNode(ast *ast.AST) (*VariableModifyNode, error) {
	node := &VariableModifyNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	node.Identifier = startToken.Value

	_, err = ast.ConsumeToken(token.TokenType_Modify, token.TokenType_Modify.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	valueTokens, err := ast.ConsumeTokenUntilMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	if err != nil {
		return nil, err
	}
	node.Value, err = CreateValueNode(valueTokens, CreateValueNodeOptions{})
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
