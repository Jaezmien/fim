package nodes

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/spike/variable"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"

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
	if len(tokens) == 0 && options.possibleNullType != nil {
		if options.possibleNullType != nil && options.possibleNullType.IsArray() {
			literalNode := &LiteralNode{
				Node: Node{
					Start:  0,
					Length: 0,
				},
				DynamicVariable: variable.NewDictionaryVariable(*options.possibleNullType),
			}

			return literalNode, nil
		}

		defaultValue, ok := options.possibleNullType.GetDefaultValue()
		if !ok {
			panic("AST@CreateValueNode (len 0, possibly unknown type?)")
		}

		literalNode := &LiteralNode{
			Node: Node{
				Start:  0,
				Length: 0,
			},
			DynamicVariable: variable.FromValueType(defaultValue, *options.possibleNullType),
		}

		return literalNode, nil
	}

	if len(tokens) == 0 {
		panic("AST@CreateValueNode called without any tokens")
	}

	if len(tokens) == 1 {
		t := tokens[0]

		if t.Value == "nothing" {
			nullNode := &LiteralNode{
				Node: Node{
					Start:  0,
					Length: 0,
				},
				DynamicVariable: variable.NewUnknownVariable(),
			}

			return nullNode, nil
		}

		if t.Type == token.TokenType_Identifier {
			node := &IdentifierNode{
				Node: Node{
					Start:  t.Start,
					Length: t.Length,
				},
				Identifier: t.Value,
			}

			if options.possibleNullType != nil && options.possibleNullType.IsArray() {
				arrayNode := wrapAsDictionaryNode(node, *options.possibleNullType, t.Start, t.Length)
				return arrayNode, nil
			}

			return node, nil
		}

		defaultType := variable.FromTokenType(t.Type)

		literalNode := &LiteralNode{
			Node: Node{
				Start:  0,
				Length: 0,
			},
		}

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
				panic("AST@CreateValueNode (possibly unknown type?)")
			}

			literalNode.DynamicVariable = variable.FromValueType(
				defaultValue,
				*options.possibleNullType,
			)

			return literalNode, nil
		}
	}

	if len(tokens) > 1 {
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
			expressionNode, err := CreateExpression(tokens, expression.tokenType, expression.operator, expression.binaryType)
			if err != nil {
				return nil, err
			}

			if expressionNode != nil {
				return expressionNode, nil
			}
		}

		// LiteralDictionaryNode
		if luna.IndexFunc(tokens, func(t *token.Token) bool { return t.Type == token.TokenType_Punctuation && t.Value == "," }, 0) != -1 {
			if options.possibleNullType == nil || !options.possibleNullType.IsArray() {
				panic("AST@CreateValueNode LiteralDictionaryNode with invalid options.possibleNullType")
			}

			baseType := options.possibleNullType.AsBaseType()
			values := make(map[int]DynamicNode, 0)

			lastSeenIndex := 0
			currentPairIndex := 1

			for {
				nextIndex := luna.IndexFunc(tokens, func(t *token.Token) bool { return t.Type == token.TokenType_Punctuation && t.Value == "," }, lastSeenIndex)

				count := 0
				if nextIndex == -1 {
					count = len(tokens)
				} else {
					count = nextIndex
				}
				count -= lastSeenIndex

				value, err := CreateValueNode(tokens[lastSeenIndex:lastSeenIndex+count], CreateValueNodeOptions{
					possibleNullType: &baseType,
				})
				if err != nil {
					return nil, err
				}

				values[currentPairIndex] = value

				if nextIndex == -1 {
					break
				}
				lastSeenIndex = nextIndex + 1
				currentPairIndex += 1
			}

			startToken := tokens[0]
			endToken := tokens[len(tokens)-1]

			dictionaryNode := &LiteralDictionaryNode{
				Node:      *NewNode(startToken.Start, endToken.Start+endToken.Length-startToken.Start),
				Values:    values,
				ArrayType: *options.possibleNullType,
			}

			return dictionaryNode, nil
		}

		// DictionaryIdentifierNode
		dinOfIndex := slices.IndexFunc(tokens, func(t *token.Token) bool { return t.Type == token.TokenType_KeywordOf })
		if dinOfIndex != -1 && dinOfIndex < len(tokens)-1 {
			indexTokens := tokens[:dinOfIndex]
			identifierTokens := tokens[dinOfIndex+1:]

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

			startToken := tokens[0]
			endToken := tokens[len(tokens)-1]

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
