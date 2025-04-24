package nodes

import (
	"fmt"
	"slices"

	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type FunctionNode struct {
	Node

	Name string
	Main bool

	Body *StatementsNode

	Parameters []FunctionNodeParameter
	ReturnType variable.VariableType
}

type FunctionNodeParameter struct {
	Name         string
	VariableType variable.VariableType
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
				return nil, parameterToken.CreateError("Parameter already exists", ast.Source)
			}

			isExpectingComma := false
			for {
				if ast.Peek().Type == token.TokenType_FunctionReturn {
					break
				}

				if ast.Peek().Type == token.TokenType_Punctuation {
					if ast.Peek().Value == "," {
						ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected '%s'"))
						isExpectingComma = false
						continue
					}
					break
				}

				if isExpectingComma {
					return nil, parameterToken.CreateError("Expecting parameters separated by comma", ast.Source)
				}

				paramTypeToken := ast.Peek()
				paramType := variable.FromTokenTypeHint(paramTypeToken.Type)
				if paramType == variable.UNKNOWN {
					return nil, paramTypeToken.CreateError("Expected variable type", ast.Source)
				}
				ast.Next()

				literalIdentifier, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
				if err != nil {
					return nil, err
				}

				if idx := slices.IndexFunc(function.Parameters, func(p FunctionNodeParameter) bool { return p.Name == literalIdentifier.Value }); idx != -1 {
					return nil, literalIdentifier.CreateError("Parameter already exists", ast.Source)
				}

				function.Parameters = append(function.Parameters, FunctionNodeParameter{
					Name:         literalIdentifier.Value,
					VariableType: paramType,
				})

				isExpectingComma = true
			}

			continue
		}
		if ast.CheckType(token.TokenType_FunctionReturn) {
			returnToken, _ := ast.ConsumeToken(token.TokenType_FunctionReturn, token.TokenType_FunctionReturn.Message("Expected %s"))

			if function.ReturnType != variable.UNKNOWN {
				return nil, returnToken.CreateError("Return type already exists", ast.Source)
			}

			returnTypeHintToken := ast.Peek()
			possibleTypeHint := variable.FromTokenTypeHint(returnTypeHintToken.Type)
			if possibleTypeHint == variable.UNKNOWN {
				return nil, returnTypeHintToken.CreateError("Expected variable type hint", ast.Source)
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
		return nil, endNameToken.CreateError(fmt.Sprintf("Mismatch method name. Expected '%s', got '%s'", startNameToken.Value, endNameToken.Value), ast.Source)
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
	Parameters []DynamicNode
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
	if err != nil {
		return nil, err
	}

	nameToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	call.Identifier = nameToken.Value

	if ast.CheckType(token.TokenType_FunctionParameter) {
		parameterToken, _ := ast.ConsumeToken(token.TokenType_FunctionParameter, token.TokenType_FunctionParameter.Message("Expected %s"))

		isExpectingComma := false
		for {
			if ast.Peek().Type == token.TokenType_Punctuation {
				if ast.Peek().Value == "," {
					ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected '%s'"))
					isExpectingComma = false
					continue
				}
				break
			}

			if isExpectingComma {
				return nil, parameterToken.CreateError("Expecting parameters separated by comma", ast.Source)
			}

			paramTypeToken := ast.Peek()
			paramType := variable.FromTokenTypeHint(paramTypeToken.Type)
			if paramType != variable.UNKNOWN {
				// TODO: Maybe also store the type hint that we will compare against
				// the evaluated valueNode as a safety type hint.
				ast.Next()
			}

			valueTokens, err := ast.ConsumeUntilFuncMatch(func(t *token.Token) bool {
				return t.Type == token.TokenType_Punctuation
			}, token.TokenType_Punctuation.Message("Could not find %s"))
			if err != nil {
				return nil, err
			}

			valueNode, err := CreateValueNode(valueTokens, CreateValueNodeOptions{})
			if err != nil {
				return nil, err
			}

			call.Parameters = append(call.Parameters, valueNode)

			isExpectingComma = true
		}
	}

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

	Value DynamicNode
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
