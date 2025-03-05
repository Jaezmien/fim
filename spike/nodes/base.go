package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type Node struct {
	Start  int
	Length int
}

type NodeType uint

const (
	TYPE_REPORT NodeType = iota
	TYPE_FUNCTION

	TYPE_STATEMENTS
	TYPE_PRINT

	TYPE_LITERAL
	TYPE_IDENTIFIER
)

type INode interface {
	Type() NodeType
}

func NewNode(start int, length int) *Node {
	return &Node{
		Start:  start,
		Length: length,
	}
}

type CreateValueNodeOptions struct {
	possibleNullType *vartype.VariableType
}

func CreateValueNode(tokens []*token.Token, options CreateValueNodeOptions) INode {
	if len(tokens) == 0 && options.possibleNullType != nil {
		if options.possibleNullType.IsArray() {
			panic("Not implemented yet")
		}

		defaultValue, ok := options.possibleNullType.GetDefaultValue()
		if !ok {
			panic("possibly unknown type?")
		}

		literalNode := &LiteralNode{
			Node: Node{
				Start:  0,
				Length: 0,
			},
			value:     defaultValue,
			ValueType: *options.possibleNullType,
		}

		return literalNode
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

			return node
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

			return literalNode
		}
		if t.Type == token.TokenType_Null && options.possibleNullType != nil {
			literalNode.ValueType = *options.possibleNullType

			defaultValue, ok := literalNode.ValueType.GetDefaultValue()
			if !ok {
				panic("possibly unknown type?")
			}
			literalNode.value = defaultValue

			return literalNode
		}
	}

	panic("ast@CreateValueNode TODO/UNKNOWN")
}
