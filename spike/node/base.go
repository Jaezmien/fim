package node

import "git.jaezmien.com/Jaezmien/fim/luna/errors"

type Node struct {
	Start  int
	Length int
}

func (n Node) ToNode() Node {
	return *NewNode(n.Start, n.Length)
}

func (n Node) CreateError(msg string, source string) error {
	return errors.NewParseError(msg, source, n.Start)
}

type DynamicNode interface {
	ToNode() Node
}

func NewNode(start int, length int) *Node {
	return &Node{
		Start:  start,
		Length: length,
	}
}
