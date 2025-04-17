package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type VariableModifyNode struct {
	Node

	Identifier string

	Value             DynamicNode
	ReinforcementType vartype.VariableType
}

func (d *VariableModifyNode) Type() NodeType {
	return TYPE_VARIABLE_MODIFY
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

	possibleTypeToken := ast.Peek()
	possibleType := vartype.FromTokenTypeHint(possibleTypeToken.Type)
	if possibleType != vartype.UNKNOWN && !possibleType.IsArray() {
		node.ReinforcementType = vartype.UNKNOWN
		ast.Next()
	} else {
		node.ReinforcementType = vartype.UNKNOWN
	}

	valueTokens, err := ast.ConsumeUntilTokenMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
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

type ArrayModifyNode struct {
	Node

	Identifier string

	Index DynamicNode

	Value             DynamicNode
	ReinforcementType vartype.VariableType
}

func (d *ArrayModifyNode) Type() NodeType {
	return TYPE_ARRAY_MODIFY
}
func (f *ArrayModifyNode) ToNode() Node {
	return Node{
		Start:  f.Start,
		Length: f.Length,
	}
}

func ParseArrayModifyNode(ast *ast.AST) (*ArrayModifyNode, error) {
	node := &ArrayModifyNode{}

	indexTokens, err := ast.ConsumeUntilTokenMatch(token.TokenType_KeywordOf, token.TokenType_KeywordOf.Message("Could not find %s"))
	if err != nil {
		return nil, err
	}
	node.Index, err = CreateValueNode(indexTokens, CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}

	_, err = ast.ConsumeToken(token.TokenType_KeywordOf, token.TokenType_KeywordOf .Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	identifierToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	node.Identifier = identifierToken.Value

	_, err = ast.ConsumeToken(token.TokenType_OperatorEq, token.TokenType_OperatorEq.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	possibleTypeToken := ast.Peek()
	possibleType := vartype.FromTokenTypeHint(possibleTypeToken.Type)
	if possibleType != vartype.UNKNOWN && !possibleType.IsArray() {
		node.ReinforcementType = vartype.UNKNOWN
		ast.Next()
	} else {
		node.ReinforcementType = vartype.UNKNOWN
	}

	valueTokens, err := ast.ConsumeUntilTokenMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
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

	node.Start = node.Index.ToNode().Start
	node.Length = endToken.Start + endToken.Length - node.Index.ToNode().Start

	return node, nil
}
