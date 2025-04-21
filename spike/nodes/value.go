package nodes

import (
	"errors"
	"fmt"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type CreateValueNodeOptions struct {
	possibleNullType *variable.VariableType
	intoArray        bool
}

func wrapAsDictionaryNode(n DynamicNode, arrayType variable.VariableType, start int, length int) DynamicNode {
	values := make(map[int]DynamicNode, 0)
	values[1] = n

	dictionaryNode := &LiteralDictionaryNode{
		Node:      *NewNode(start, length),
		Values:    values,
		ArrayType: arrayType,
	}

	return dictionaryNode
}

func CreateValueNode(tokens []*token.Token, options CreateValueNodeOptions) (DynamicNode, error) {
	tempAST := ast.NewAST(tokens, "")

	if tempAST.Length() == 0 && options.possibleNullType != nil {
		if options.possibleNullType != nil && options.possibleNullType.IsArray() {
			return NewLiteralNode(
				0, 0,
				variable.NewDictionaryVariable(*options.possibleNullType),
			), nil
		}

		defaultValue, ok := options.possibleNullType.GetDefaultValue()
		if !ok {
			panic("AST@CreateValueNode called with no tokens, and unknown possibleNullType")
		}

		return NewLiteralNode(
			0, 0,
			variable.FromValueType(defaultValue, *options.possibleNullType),
		), nil
	}

	if tempAST.Length() == 0 {
		panic("AST@CreateValueNode called without any tokens")
	}

	if tempAST.Length() == 1 {
		if tempAST.Peek().Value == "nothing" {
			return NewLiteralNode(
				0, 0,
				variable.NewUnknownVariable(),
			), nil
		}

		if tempAST.CheckType(token.TokenType_Identifier) {
			t := tempAST.Consume()

			node := &IdentifierNode{
				Node:       *NewNode(t.Start, t.Length),
				Identifier: t.Value,
			}

			if options.possibleNullType != nil && options.possibleNullType.IsArray() {
				arrayNode := wrapAsDictionaryNode(node, *options.possibleNullType, t.Start, t.Length)
				return arrayNode, nil
			}

			return node, nil
		}

		literalNode := NewLiteralNode(0, 0, nil)
		t := tempAST.Consume()

		defaultType := variable.FromTokenType(t.Type)
		if defaultType != variable.UNKNOWN {
			literalNode.DynamicVariable = variable.FromValueType(t.Value, defaultType)

			literalNode.Start = t.Start
			literalNode.Length = t.Length

			if options.possibleNullType != nil && options.possibleNullType.IsArray() {
				arrayNode := wrapAsDictionaryNode(literalNode, *options.possibleNullType, t.Start, t.Length)
				return arrayNode, nil
			}

			return literalNode, nil
		}

		if t.Type == token.TokenType_Null && options.possibleNullType != nil {
			defaultValue, ok := options.possibleNullType.GetDefaultValue()
			if !ok {
				panic("AST@CreateValueNode literal null called with no possible default value")
			}

			literalNode.DynamicVariable = variable.FromValueType(
				defaultValue,
				*options.possibleNullType,
			)

			return literalNode, nil
		}
	}

	if tempAST.Length() > 1 {
		// FunctionCallNode
		if tempAST.Peek().Type == token.TokenType_Identifier && tempAST.PeekNext().Type == token.TokenType_FunctionParameter {
			identifier := tempAST.Consume()
			tempAST.Consume()

			callNode := &FunctionCallNode{
				Node: *NewNode(identifier.Start, tempAST.End().Start+tempAST.End().Length-identifier.Start),
				Identifier: identifier.Value,
				Parameters: make([]DynamicNode, 0),
			}

			isExpectingComma := false
			for tempAST.PeekIndex() < tempAST.Length() {
				if tempAST.CheckType(token.TokenType_Punctuation) {
					if tempAST.Peek().Value == "," {
						tempAST.Consume()
						isExpectingComma = false
						continue
					}
					break
				}

				if isExpectingComma {
					return nil, errors.New("Expecting parameters separated by comma")
				}

				paramTypeToken := tempAST.Peek()
				paramType := variable.FromTokenTypeHint(paramTypeToken.Type)
				if paramType != variable.UNKNOWN {
					// TODO: Maybe also store the type hint that we will compare against
					// the evaluated valueNode as a safety type hint.
					tempAST.Next()
				}

				valueTokens, err := tempAST.ConsumeUntilFuncMatch(func(t *token.Token) bool {
					return t.Type == token.TokenType_Punctuation
				}, token.TokenType_Punctuation.Message("Could not find %s"))
				if err != nil {
					return nil, err
				}

				valueNode, err := CreateValueNode(valueTokens, CreateValueNodeOptions{})
				if err != nil {
					return nil, err
				}

				callNode.Parameters = append(callNode.Parameters, valueNode)

				isExpectingComma = true
			}

			return callNode, nil
		}
		expressions := []struct {
			tokenType  token.TokenType
			operator   BinaryExpressionOperator
			binaryType BinaryExpressionType
		}{
			// Arithmetic
			{
				tokenType:  token.TokenType_OperatorMulInfix,
				operator:   BINARYOPERATOR_MUL,
				binaryType: BINARYTYPE_ARITHMETIC,
			},
			{
				tokenType:  token.TokenType_OperatorDivInfix,
				operator:   BINARYOPERATOR_DIV,
				binaryType: BINARYTYPE_ARITHMETIC,
			},
			{
				tokenType:  token.TokenType_OperatorAddInfix,
				operator:   BINARYOPERATOR_ADD,
				binaryType: BINARYTYPE_ARITHMETIC,
			},
			{
				tokenType:  token.TokenType_OperatorSubInfix,
				operator:   BINARYOPERATOR_SUB,
				binaryType: BINARYTYPE_ARITHMETIC,
			},

			// Relational
			{
				tokenType:  token.TokenType_OperatorGte,
				operator:   BINARYOPERATOR_GTE,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType:  token.TokenType_OperatorLte,
				operator:   BINARYOPERATOR_LTE,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType:  token.TokenType_OperatorGt,
				operator:   BINARYOPERATOR_GT,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType:  token.TokenType_OperatorLt,
				operator:   BINARYOPERATOR_LT,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType:  token.TokenType_OperatorNeq,
				operator:   BINARYOPERATOR_NEQ,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType:  token.TokenType_OperatorEq,
				operator:   BINARYOPERATOR_EQ,
				binaryType: BINARYTYPE_RELATIONAL,
			},

			{
				tokenType:  token.TokenType_KeywordAnd,
				operator:   BINARYOPERATOR_AND,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType:  token.TokenType_KeywordOr,
				operator:   BINARYOPERATOR_OR,
				binaryType: BINARYTYPE_RELATIONAL,
			},
		}

		for _, expression := range expressions {
			expressionNode, err := CreateExpression(tempAST.Tokens, expression.tokenType, expression.operator, expression.binaryType)
			if err != nil {
				return nil, err
			}

			if expressionNode != nil {
				return expressionNode, nil
			}
		}

		// LiteralDictionaryNode
		if tempAST.ContainsFunc(func(t *token.Token) bool {
			return t.Type == token.TokenType_Punctuation && t.Value == ","
		}) {
			if options.possibleNullType == nil || !options.possibleNullType.IsArray() {
				panic("AST@CreateValueNode LiteralDictionaryNode with invalid options.possibleNullType")
			}

			baseType := options.possibleNullType.AsBaseType()
			values := make(map[int]DynamicNode, 0)
			currentIndex := 1

			for tempAST.PeekIndex() < tempAST.Length() {
				if currentIndex != 1 {
					_, err := tempAST.ConsumeFunc(func(t *token.Token) bool {
						return t.Type == token.TokenType_Punctuation && t.Value == ","
					}, token.TokenType_Punctuation.Message("Expected '%s' (Comma)"))

					if err != nil {
						return nil, err 
					}
				}

				ts, _ := tempAST.ConsumeUntilFuncMatch(func(t *token.Token) bool {
					return t.Type == token.TokenType_Punctuation && t.Value == ","
				}, "")

				if len(ts) == 0 {
					break
				}

				value, err := CreateValueNode(ts, CreateValueNodeOptions{
					possibleNullType: &baseType,
				})
				if err != nil {
					return nil, err
				}

				values[currentIndex] = value
				currentIndex += 1
			}

			startToken := tempAST.Start()
			endToken := tempAST.End()

			dictionaryNode := &LiteralDictionaryNode{
				Node:      *NewNode(startToken.Start, endToken.Start+endToken.Length-startToken.Start),
				Values:    values,
				ArrayType: *options.possibleNullType,
			}

			return dictionaryNode, nil
		}

		// DictionaryIdentifierNode
		if tempAST.ContainsToken(token.TokenType_KeywordOf) && tempAST.End().Type != token.TokenType_KeywordOf {
			indexTokens, _ := tempAST.ConsumeUntilTokenMatch(token.TokenType_KeywordOf, "")
			tempAST.Consume()
			identifierTokens := tempAST.ConsumeRemaining()

			if len(indexTokens) < 1 {
				return nil, errors.New("Expected dictionary identifier index")
			}

			index, err := CreateValueNode(indexTokens, CreateValueNodeOptions{})
			if err != nil {
				return nil, err
			}

			if len(identifierTokens) != 1 || identifierTokens[0].Type != token.TokenType_Identifier {
				return nil, errors.New("Expected dictionary identifier")
			}

			startToken := tempAST.Start()
			endToken := tempAST.End()

			identifierNode := &DictionaryIdentifierNode{
				Node:       *NewNode(startToken.Start, endToken.Start+endToken.Length-startToken.Start),
				Identifier: identifierTokens[0].Value,
				Index:      index,
			}

			return identifierNode, nil
		}
	}

	unknownToken := strings.Builder{}
	for _, token := range tokens {
		unknownToken.WriteString(token.Value)
	}
	return nil, fmt.Errorf("Encountered unknown value token: '%s'", unknownToken.String())
}
