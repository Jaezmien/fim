package nodes

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type FunctionNode struct {
	Node

	Name string
	Main bool

	Body *StatementsNode

	Parameters []FunctionNodeParameter
	ReturnType vartype.VariableType
}

type FunctionNodeParameter struct {
	Name         string
	VariableType vartype.VariableType
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
		Main:       false,
		Parameters: make([]FunctionNodeParameter, 0),
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

	for {
		if ast.CheckType(token.TokenType_FunctionParameter) {
			parameterToken, _ := ast.ConsumeToken(token.TokenType_FunctionParameter, token.TokenType_FunctionParameter.Message("Expected %s"))

			if len(function.Parameters) > 0 {
				return nil, ast.CreateErrorFromToken(parameterToken, "Return type already exists")
			}

			continue
		}
		if ast.CheckType(token.TokenType_FunctionReturn) {
			returnToken, _ := ast.ConsumeToken(token.TokenType_FunctionReturn, token.TokenType_FunctionReturn.Message("Expected %s"))

			if function.ReturnType != vartype.UNKNOWN {
				return nil, ast.CreateErrorFromToken(returnToken, "Return type already exists")
			}

			returnTypeHintToken := ast.Peek()
			possibleTypeHint := vartype.FromTokenTypeHint(returnTypeHintToken.Type)
			if possibleTypeHint == vartype.UNKNOWN {
				return nil, ast.CreateErrorFromToken(returnTypeHintToken, "Expected variable type hint")
			}
			function.ReturnType = possibleTypeHint

			ast.Next()

			continue
		}

		break
	}

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

// --- //

type FunctionCallNode struct {
	Node

	Identifier string
}

func (f *FunctionCallNode) Type() NodeType {
	return TYPE_FUNCTION_CALL
}

func (f *FunctionCallNode) ToNode() Node {
	return Node{
		Start:  f.Start,
		Length: f.Length,
	}
}

func ParseFunctionCallNode(ast *ast.AST) (*FunctionCallNode, error) {
	call := &FunctionCallNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_FunctionCall, token.TokenType_FunctionCall.Message("Expected %s"))

	nameToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	call.Identifier = nameToken.Value

	// TODO: Parameter check
	// TODO: Return check

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	call.Start = startToken.Start
	call.Length = endToken.Start + endToken.Length - startToken.Start

	return call, nil
}

// --- //

type FunctionReturnNode struct {
	Node

	Value INode
}

func (f *FunctionReturnNode) Type() NodeType {
	return TYPE_FUNCTION_RETURN
}

func (f *FunctionReturnNode) ToNode() Node {
	return Node{
		Start:  f.Start,
		Length: f.Length,
	}
}

func ParseFunctionReturnNode(ast *ast.AST) (*FunctionReturnNode, error) {
	returnNode := &FunctionReturnNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_KeywordReturn, token.TokenType_KeywordReturn.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	valueTokens, err := ast.ConsumeUntilTokenMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	if err != nil {
		return nil, err
	}

	returnNode.Value, err = CreateValueNode(valueTokens, CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	returnNode.Start = startToken.Start
	returnNode.Length = endToken.Start + endToken.Length - startToken.Start

	return returnNode, nil
}
