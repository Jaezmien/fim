package node

type Node struct {
	Start  int
	Length int
}

type NodeType uint

const (
	TYPE_REPORT NodeType = iota

	TYPE_FUNCTION
	TYPE_FUNCTION_CALL
	TYPE_FUNCTION_RETURN

	TYPE_STATEMENTS
	TYPE_PRINT
	TYPE_PROMPT

	TYPE_LITERAL
	TYPE_LITERAL_DICTIONARY
	TYPE_IDENTIFIER
	TYPE_IDENTIFIER_DICTIONARY
	TYPE_BINARYEXPRESSION

	TYPE_VARIABLE_DECLARATION
	TYPE_VARIABLE_MODIFY
)

var nodeTypeFriendlyName = map[NodeType]string{
	TYPE_REPORT:               "REPORT",
	TYPE_FUNCTION:             "FUNCTION",
	TYPE_STATEMENTS:           "STATEMENTS",
	TYPE_PRINT:                "PRINT",
	TYPE_LITERAL:              "LITERAL",
	TYPE_LITERAL_DICTIONARY:   "LITERAL_DICTIONARY",
	TYPE_IDENTIFIER:           "IDENTIFIER",
	TYPE_BINARYEXPRESSION:     "BINARYEXPRESSION",
	TYPE_VARIABLE_DECLARATION: "VARIABLE_DECLARATION",
	TYPE_VARIABLE_MODIFY:      "VARIABLE_MODIFY",
}

func (t NodeType) String() string {
	return nodeTypeFriendlyName[t]
}

type INode interface {
	Type() NodeType
	ToNode() Node
}

func NewNode(start int, length int) *Node {
	return &Node{
		Start:  start,
		Length: length,
	}
}
