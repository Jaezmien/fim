package nodes

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
