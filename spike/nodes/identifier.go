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

func (i *IdentifierNode) ToNode() Node {
	return Node{
		Start:  i.Start,
		Length: i.Length,
	}
}

type DictionaryIdentifierNode struct {
	Node

	Identifier string
	Index INode
}

func (i *DictionaryIdentifierNode) Type() NodeType {
	return TYPE_IDENTIFIER_DICTIONARY
}

func (i *DictionaryIdentifierNode) ToNode() Node {
	return Node{
		Start:  i.Start,
		Length: i.Length,
	}
}
