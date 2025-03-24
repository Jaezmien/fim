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
var variableTypeFriendlyName = map[VariableType]string {
	UNKNOWN: "",
	BOOLEAN: "BOOLEAN",
	CHARACTER: "CHARACTER",
	NUMBER: "NUMBER",
	STRING: "STRING",
	BOOLEAN_ARRAY: "ARRAY(BOOLEAN)",
	NUMBER_ARRAY: "ARRAY(NUMBER)",
	STRING_ARRAY: "ARRAY(STRING)",
}
func (t VariableType) String() string {
	return variableTypeFriendlyName[t]
}

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

// Returns the base type of array types
func (t VariableType) AsBaseType() VariableType {
	switch t {
	case BOOLEAN_ARRAY:
		return BOOLEAN
	case NUMBER_ARRAY:
		return NUMBER
	case STRING_ARRAY:
		return STRING
	default:
		return UNKNOWN
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
func FromTokenTypeHint(t token.TokenType) VariableType {
	switch t {
	case token.TokenType_TypeBoolean:
		return BOOLEAN
	case token.TokenType_TypeNumber:
		return NUMBER
	case token.TokenType_TypeChar:
		return CHARACTER
	case token.TokenType_TypeString:
		return STRING
	case token.TokenType_TypeBooleanArray:
		return BOOLEAN_ARRAY
	case token.TokenType_TypeNumberArray:
		return NUMBER_ARRAY
	case token.TokenType_TypeStringArray:
		return STRING_ARRAY
	default:
		return UNKNOWN
	}
}
