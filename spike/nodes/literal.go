package nodes

import (
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
)

type LiteralNode struct {
	Node

	*variable.DynamicVariable
}

func (l *LiteralNode) Type() NodeType {
	return TYPE_LITERAL
}

// -- //

type LiteralDictionaryNode struct {
	Node

	ArrayType variable.VariableType
	Values map[int]DynamicNode
}

func (l *LiteralDictionaryNode) Type() NodeType {
	return TYPE_LITERAL_DICTIONARY
}
