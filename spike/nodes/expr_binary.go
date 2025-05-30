package nodes

import (
	"slices"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type BinaryExpressionType uint

const (
	BINARYTYPE_UNKNOWN BinaryExpressionType = iota
	BINARYTYPE_ARITHMETIC
	BINARYTYPE_RELATIONAL
)

type BinaryExpressionOperator uint

const (
	BINARYOPERATOR_UNKNOWN BinaryExpressionOperator = iota
	BINARYOPERATOR_ADD
	BINARYOPERATOR_SUB
	BINARYOPERATOR_MUL
	BINARYOPERATOR_DIV
	BINARYOPERATOR_MOD

	BINARYOPERATOR_AND
	BINARYOPERATOR_OR
	BINARYOPERATOR_GTE
	BINARYOPERATOR_LTE
	BINARYOPERATOR_GT
	BINARYOPERATOR_LT
	BINARYOPERATOR_NEQ
	BINARYOPERATOR_EQ
)

type BinaryExpressionNode struct {
	Node

	Left     DynamicNode
	Operator BinaryExpressionOperator
	Right    DynamicNode

	BinaryType BinaryExpressionType
}

func CreateExpression(tokens []*token.Token, tokenType token.TokenType, operator BinaryExpressionOperator, binaryType BinaryExpressionType) (*BinaryExpressionNode, error) {
	// XXX: In FiMSharp, I used .FindLastIndex for the BinaryExpressionNode.
	// Do I really need to do the same in this instance? I still have to check.
	index := slices.IndexFunc(tokens, func(t *token.Token) bool { return t.Type == tokenType })

	if index == -1 {
		return nil, nil
	}

	leftNode, err := CreateValueNode(tokens[:index], CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}
	rightNode, err := CreateValueNode(tokens[index+1:], CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}

	node := &BinaryExpressionNode{
		Left:       leftNode,
		Right:      rightNode,
		Operator:   operator,
		BinaryType: binaryType,
		Node: Node{
			Start:  leftNode.ToNode().Start,
			Length: rightNode.ToNode().Start + rightNode.ToNode().Length - leftNode.ToNode().Start,
		},
	}

	return node, nil
}
