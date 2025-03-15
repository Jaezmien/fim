package nodes

import (
	"slices"

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

	Left     INode
	Operator BinaryExpressionOperator
	Right    INode

	BinaryType BinaryExpressionType
}

func (b *BinaryExpressionNode) Type() NodeType {
	return TYPE_BINARYEXPRESSION
}
func (b *BinaryExpressionNode) ToNode() Node {
	return Node{
		Start:  b.Start,
		Length: b.Length,
	}
}

func CreateExpression(tokens []*token.Token, tokenType token.TokenType, operator BinaryExpressionOperator, binaryType BinaryExpressionType) (*BinaryExpressionNode, bool) {
	index := slices.IndexFunc(tokens, func(t *token.Token) bool { return t.Type == tokenType; })

	if index == -1 {
		return nil, false
	}

	// XXX: Do I really need .FindLastIndex on this?
	leftNode := CreateValueNode(tokens[:index], CreateValueNodeOptions{})
	rightNode := CreateValueNode(tokens[index+1:], CreateValueNodeOptions{})

	node := &BinaryExpressionNode{
		Left: leftNode,
		Right: rightNode,
		Operator: operator,
		BinaryType: binaryType,
		Node: Node{
			Start: leftNode.ToNode().Start,
			Length: rightNode.ToNode().Start + rightNode.ToNode().Length - leftNode.ToNode().Start,
		},
	}

	return node, true
}
