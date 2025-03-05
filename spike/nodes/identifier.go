package nodes

type IdentifierNode struct {
	Node

	Identifier string
}

func (i *IdentifierNode) Type() NodeType {
	return TYPE_IDENTIFIER
}
