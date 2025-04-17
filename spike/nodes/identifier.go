package nodes

import (
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type IdentifierNode struct {
	Node

	Identifier string
}

func (i *IdentifierNode) Type() NodeType {
	return TYPE_IDENTIFIER
}

type DictionaryIdentifierNode struct {
	Node

	Identifier string
	Index      DynamicNode
}

func (i *DictionaryIdentifierNode) Type() NodeType {
	return TYPE_IDENTIFIER_DICTIONARY
}
