package nodes

import (
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
)

type LiteralNode struct {
	Node

	*vartype.DynamicVariable
}

func (l *LiteralNode) Type() NodeType {
	return TYPE_LITERAL
}
func (l *LiteralNode) ToNode() Node {
	return Node{
		Start:  l.Start,
		Length: l.Length,
	}
}

// -- //

type LiteralDictionaryNode struct {
	Node

	ArrayType vartype.VariableType
	Values map[int]INode
}

func (l *LiteralDictionaryNode) Type() NodeType {
	return TYPE_LITERAL_DICTIONARY
}
func (l *LiteralDictionaryNode) ToNode() Node {
	return Node{
		Start:  l.Start,
		Length: l.Length,
	}
}
