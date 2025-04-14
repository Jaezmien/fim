package node

import "git.jaezmien.com/Jaezmien/fim/luna/errors"

type Node struct {
	Start  int
	Length int
}

func (n Node) CreateError(msg string, source string) error {
	return errors.NewParseError(msg, source, n.Start)
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
	TYPE_ARRAY_MODIFY

	TYPE_UNARYEXPRESSION
)

var nodeTypeFriendlyName = map[NodeType]string{
	TYPE_REPORT:          "REPORT",
	TYPE_FUNCTION:        "FUNCTION",
	TYPE_FUNCTION_CALL:   "FUNCTION(CALL)",
	TYPE_FUNCTION_RETURN: "FUNCTION(RETURN)",

	TYPE_STATEMENTS: "STATEMENTS",
	TYPE_PRINT:      "PRINT",
	TYPE_PROMPT:     "PROMPT",

	TYPE_LITERAL:               "LITERAL",
	TYPE_LITERAL_DICTIONARY:    "LITERAL_DICTIONARY",
	TYPE_IDENTIFIER:            "IDENTIFIER",
	TYPE_IDENTIFIER_DICTIONARY: "IDENTIFIER(DICTIONARY)",
	TYPE_BINARYEXPRESSION:      "BINARYEXPRESSION",

	TYPE_VARIABLE_DECLARATION: "VARIABLE_DECLARATION",
	TYPE_VARIABLE_MODIFY:      "VARIABLE_MODIFY",
	TYPE_ARRAY_MODIFY:         "ARRAY_MODIFY",
	 
	TYPE_UNARYEXPRESSION:      "UNARYEXPRESSION",
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
