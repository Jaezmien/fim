package nodes

type Node struct {
	Start int
	Length int
}

func NewNode(start int, length int) *Node {
	return &Node {
		Start: start,
		Length: length,
	}
}
