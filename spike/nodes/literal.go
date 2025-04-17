package nodes

import (
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
)

type LiteralNode struct {
	Node

	*variable.DynamicVariable
}

// -- //

type LiteralDictionaryNode struct {
	Node

	ArrayType variable.VariableType
	Values map[int]DynamicNode
}
