package nodes

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
)

type CreateValueNodeOptions struct {
	possibleNullType *vartype.VariableType
	asArray bool
}

func CreateValueNode(tokens []*token.Token, options CreateValueNodeOptions) (INode, error) {
	if len(tokens) == 0 && options.possibleNullType != nil {
		if options.possibleNullType.IsArray() {
			panic("AST@CreateValueNode (Not implemented yet)")
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
			value:     defaultValue,
			ValueType: *options.possibleNullType,
		}

		return literalNode, nil
	}

	if len(tokens) == 0 {
		panic("AST@CreateValueNode called without any tokens")
	}

	if len(tokens) == 1 {
		t := tokens[0]

		if t.Type == token.TokenType_Identifier {
			node := &IdentifierNode{
				Node: Node{
					Start:  t.Start,
					Length: t.Length,
				},
				Identifier: t.Value,
			}

			return node, nil
		}

		defaultType := vartype.FromTokenType(t.Type)

		literalNode := &LiteralNode{
			Node: Node{
				Start:  0,
				Length: 0,
			},
		}

		if defaultType != vartype.UNKNOWN {
			literalNode.value = t.Value
			literalNode.ValueType = defaultType

			literalNode.Start = t.Start
			literalNode.Length = t.Length

			return literalNode, nil
		}
		if t.Type == token.TokenType_Null && options.possibleNullType != nil {
			literalNode.ValueType = *options.possibleNullType

			defaultValue, ok := literalNode.ValueType.GetDefaultValue()
			if !ok {
				panic("AST@CreateValueNode (possibly unknown type?)")
			}
			literalNode.value = defaultValue

			return literalNode, nil
		}
	}

	if len(tokens) > 1 {
		// BinaryExpressionNode
		expressions := []struct{
			tokenType token.TokenType
			operator BinaryExpressionOperator
			binaryType BinaryExpressionType
		}{
			// Arithmetic
			{
				tokenType: token.TokenType_OperatorMulInfix,
				operator: BINARYOPERATOR_MUL,
				binaryType: BINARYTYPE_ARITHMETIC,
			},
			{
				tokenType: token.TokenType_OperatorDivInfix,
				operator: BINARYOPERATOR_DIV,
				binaryType: BINARYTYPE_ARITHMETIC,
			},
			{
				tokenType: token.TokenType_OperatorAddInfix,
				operator: BINARYOPERATOR_ADD,
				binaryType: BINARYTYPE_ARITHMETIC,
			},
			{
				tokenType: token.TokenType_OperatorSubInfix,
				operator: BINARYOPERATOR_SUB,
				binaryType: BINARYTYPE_ARITHMETIC,
			},
			
			// Relational
			{
				tokenType: token.TokenType_OperatorGte,
				operator: BINARYOPERATOR_GTE,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType: token.TokenType_OperatorLte,
				operator: BINARYOPERATOR_LTE,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType: token.TokenType_OperatorGt,
				operator: BINARYOPERATOR_GT,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType: token.TokenType_OperatorLt,
				operator: BINARYOPERATOR_LT,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType: token.TokenType_OperatorNeq,
				operator: BINARYOPERATOR_NEQ,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType: token.TokenType_OperatorEq,
				operator: BINARYOPERATOR_EQ,
				binaryType: BINARYTYPE_RELATIONAL,
			},

			{
				tokenType: token.TokenType_KeywordAnd,
				operator: BINARYOPERATOR_AND,
				binaryType: BINARYTYPE_RELATIONAL,
			},
			{
				tokenType: token.TokenType_KeywordOr,
				operator: BINARYOPERATOR_OR,
				binaryType: BINARYTYPE_RELATIONAL,
			},
		}

		for _, expression := range expressions {
			expressionNode, _ := CreateExpression(tokens, expression.tokenType, expression.operator, expression.binaryType)

			if expressionNode != nil {
				return expressionNode, nil
			}
		}

		// LiteralDictionaryNode
		if slices.IndexFunc(tokens, func(t *token.Token) bool { return t.Type == token.TokenType_Punctuation && t.Value == "," }) != -1 {
			if options.possibleNullType == nil || !options.possibleNullType.IsArray() {
				panic("AST@CreateValueNode LiteralDictionaryNode with invalid options.possibleNullType")
			}

			baseType := options.possibleNullType.AsBaseType()
			var pairs map[int]INode = make(map[int]INode)

			lastSeenIndex := 0
			currentPairIndex := 1

			for {
				nextIndex := luna.IndexFunc(tokens, func(t *token.Token) bool { return t.Type == token.TokenType_Punctuation && t.Value == "," })
				if lastSeenIndex <= nextIndex {
					break
				}

				count := 0
				if nextIndex == -1 {
					count = len(tokens)
				} else {
					count = nextIndex
				}
				count -= lastSeenIndex

				value, err := CreateValueNode(tokens[lastSeenIndex:nextIndex], CreateValueNodeOptions{
					possibleNullType: &baseType,
				})
				if err != nil {
					return nil, err
				}

				pairs[currentPairIndex] = value

				lastSeenIndex = nextIndex + 1
				currentPairIndex += 1
			}
		}


	}

	unknownToken := strings.Builder{}
	for _, token := range tokens {
		unknownToken.WriteString(token.Value)
	}
	return nil, errors.New(fmt.Sprintf("Encountered unknown value token: '%s'", unknownToken.String()))
}
