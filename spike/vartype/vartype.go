package vartype

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type VariableType uint

const (
	UNKNOWN VariableType = iota

	BOOLEAN
	CHARACTER
	NUMBER
	STRING

	BOOLEAN_ARRAY
	NUMBER_ARRAY
	STRING_ARRAY
)

func (t VariableType) IsArray() bool {
	switch t {
	case BOOLEAN_ARRAY:
		return true
	case NUMBER_ARRAY:
		return true
	case STRING_ARRAY:
		return true
	default:
		return false
	}
}

func (t VariableType) GetDefaultValue() (string, bool) {
	switch t {
	case BOOLEAN:
		return "false", true
	case CHARACTER:
		return fmt.Sprintf("'%s'", string(rune(0))), true
	case NUMBER:
		return "0", true
	case STRING:
		return "\"\"", true
	default:
		return "", false

	}
}

func FromTokenType(t token.TokenType) VariableType {
	switch t {
	case token.TokenType_Boolean:
		return BOOLEAN
	case token.TokenType_Number:
		return NUMBER
	case token.TokenType_Character:
		return CHARACTER
	case token.TokenType_String:
		return STRING
	default:
		return UNKNOWN
	}
}
