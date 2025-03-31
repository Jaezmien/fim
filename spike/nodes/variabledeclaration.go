package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type VariableDeclarationNode struct {
	Node

	Identifier string
	Constant   bool

	Value     INode
	ValueType vartype.VariableType
}

func (d *VariableDeclarationNode) Type() NodeType {
	return TYPE_VARIABLE_DECLARATION
}
func (f *VariableDeclarationNode) ToNode() Node {
	return Node{
		Start:  f.Start,
		Length: f.Length,
	}
}

func ParseVariableDeclarationNode(ast *ast.AST) (*VariableDeclarationNode, error) {
	node := &VariableDeclarationNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_Declaration, token.TokenType_Declaration.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	identifier, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	node.Identifier = identifier.Value

	_, err = ast.ConsumeToken(token.TokenType_OperatorEq, token.TokenType_OperatorEq.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	if ast.Peek().Type == token.TokenType_KeywordConst {
		node.Constant = true
		ast.Next()
	}

	typeToken := ast.Peek()
	node.ValueType = vartype.FromTokenTypeHint(typeToken.Type)
	if node.ValueType == vartype.UNKNOWN {
		return nil, ast.CreateErrorFromToken(typeToken, "Expected variable type hint")
	}
	ast.Next()

	var valueTokens []*token.Token
	if node.ValueType.IsArray() {
		valueTokens, err = ast.ConsumeUntilFuncMatch(func(t *token.Token) bool {
			return t.Type == token.TokenType_Punctuation && t.Value != ","
		}, token.TokenType_Punctuation.Message("Could not find %s"))
	} else {
		valueTokens, err = ast.ConsumeUntilTokenMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	}
	if err != nil {
		return nil, err
	}

	node.Value, err = CreateValueNode(valueTokens, CreateValueNodeOptions{
		possibleNullType: &node.ValueType,
		intoArray:        node.ValueType.IsArray(),
	})
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
