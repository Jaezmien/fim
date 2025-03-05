package nodes

import (
	"slices"
	"strconv"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/spike/utilities"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type LiteralNode struct {
	Node

	value     string
	ValueType vartype.VariableType
}

func (l *LiteralNode) Type() NodeType {
	return TYPE_LITERAL
}

func (l *LiteralNode) SetValue(value string) {
	l.value = value
}

func (l *LiteralNode) GetValueString() string {
	if l.ValueType != vartype.STRING {
		panic("Called LiteralNode@ValueString on a non-string literal")
	}

	return utilities.UnsanitizeString(l.value, true)
}
func (l *LiteralNode) GetValueCharacter() string {
	if l.ValueType != vartype.CHARACTER {
		panic("Called LiteralNode@ValueCharacter on a non-character literal")
	}

	value := l.value[1 : len(l.value)-1]

	if strings.HasPrefix(value, "\\") {
		switch value[1] {
		case '0':
			return string(byte(0))
		case 'r':
			return "\r"
		case 'n':
			return "\n"
		case 't':
			return "\t"
		default:
			return string(value[1])
		}
	}
	return value
}
func (l *LiteralNode) GetValueBoolean() bool {
	if l.ValueType != vartype.BOOLEAN {
		panic("Called LiteralNode@ValueBoolean on a non-boolean literal")
	}

	return slices.Contains([]string{"yes", "true", "right", "correct"}, l.value)
}
func (l *LiteralNode) GetValueNumber() float64 {
	if l.ValueType != vartype.NUMBER {
		panic("Called LiteralNode@ValueNumber on a non-number literal")
	}

	value, ok := strconv.ParseFloat(l.value, 64)
	if ok != nil {
		panic(ok)
	}

	return value
}

func ParseLiteralNode(ast *spike.AST, expectedEndType ...token.TokenType) (*LiteralNode, error) {
	node := &LiteralNode{}

	return node, nil
}
