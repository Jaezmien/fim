package nodes

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type FunctionNode struct {
	Node

	Name string
	Main bool

	Body *StatementsNode
}

func (f *FunctionNode) Type() NodeType {
	return TYPE_FUNCTION
}
func (f *FunctionNode) ToNode() Node {
	return Node{
		Start:  f.Start,
		Length: f.Length,
	}
}

func ParseFunctionNode(ast *ast.AST) (*FunctionNode, error) {
	function := &FunctionNode{
		Main: false,
	}

	var startToken *token.Token
	var err error

	if ast.CheckType(token.TokenType_FunctionMain) {
		function.Main = true

		startToken, err = ast.ConsumeToken(token.TokenType_FunctionMain, token.TokenType_FunctionMain.Message("Expected %s"))
	} else {
		startToken, err = ast.ConsumeToken(token.TokenType_FunctionHeader, token.TokenType_FunctionHeader.Message("Expected %s"))
	}

	if err != nil {
		return nil, err
	}

	startNameToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	function.Name = startNameToken.Value

	// TODO: Parameter check
	// TODO: Return check

	_, err = ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	// TODO: Body check
	bodyStatement, err := ParseStatementsNode(ast, token.TokenType_FunctionFooter)
	if err != nil {
		return nil, err
	}
	function.Body = bodyStatement

	_, err = ast.ConsumeToken(token.TokenType_FunctionFooter, token.TokenType_FunctionFooter.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	endNameToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	if startNameToken.Value != endNameToken.Value {
		return nil, ast.CreateErrorFromToken(endNameToken, fmt.Sprintf("Mismatch method name. Expected '%s', got '%s'", startNameToken.Value, endNameToken.Value))
	}

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	function.Start = startToken.Start
	function.Length = endToken.Start + endToken.Length - startToken.Start

	return function, nil
}
