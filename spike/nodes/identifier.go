package nodes

import (
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type IdentifierNode struct {
	Node

	Identifier string
}

type DictionaryIdentifierNode struct {
	Node

	Identifier string
	Index      DynamicNode
}
