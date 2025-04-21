package nodes

import (
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
)

type LiteralNode struct {
	Node

	*variable.DynamicVariable
}

func NewLiteralNode(start int, length int, v *variable.DynamicVariable) *LiteralNode {
	return &LiteralNode{
		Node: *NewNode(start, length),
		DynamicVariable: v,
	}
}

// -- //

type LiteralDictionaryNode struct {
	Node

	ArrayType variable.VariableType
	Values    map[int]DynamicNode
}
