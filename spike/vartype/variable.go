package vartype

import (
	"slices"
	"strconv"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/spike/utilities"

	"git.jaezmien.com/Jaezmien/fim/spike/node"
)

type DynamicVariable struct {
	value any

	valueType VariableType
}

func FromValueType(value string, t VariableType) *DynamicVariable {
	switch t {
	case BOOLEAN:
		return NewBooleanVariable(slices.Contains([]string{"yes", "true", "right", "correct"}, value))
	case CHARACTER:
		return NewCharacterVariable(value)
	case NUMBER:
		value, ok := strconv.ParseFloat(value, 64)
		if ok != nil {
			panic(ok)
		}
		return NewNumberVariable(value)
	case STRING:
		return NewStringVariable(value)
	default:
		panic("DynamicVariable@FromDefaultValue called with variable type " + t.String())

	}
}

func NewStringVariable(value string) *DynamicVariable {
	return NewRawStringVariable(utilities.UnsanitizeString(value, true))
}
func NewRawStringVariable(value string) *DynamicVariable {
	return &DynamicVariable{
		value:     value,
		valueType: STRING,
	}
}
func NewCharacterVariable(value string) *DynamicVariable {
	value = value[1 : len(value)-1]

	if strings.HasPrefix(value, "\\") {
		switch value[1] {
		case '0':
			value = string(byte(0))
		case 'r':
			value = "\r"
		case 'n':
			value = "\n"
		case 't':
			value = "\t"
		default:
			value = string(value[1])
		}
	}

	return NewRawCharacterVariable(value)
}
func NewRawCharacterVariable(value string) *DynamicVariable {
	return &DynamicVariable{
		value:     value,
		valueType: CHARACTER,
	}
}
func NewNumberVariable(value float64) *DynamicVariable {
	return &DynamicVariable{
		value:     value,
		valueType: NUMBER,
	}
}
func NewBooleanVariable(value bool) *DynamicVariable {
	return &DynamicVariable{
		value:     value,
		valueType: BOOLEAN,
	}
}
func NewDictionaryVariable(t VariableType) *DynamicVariable {
	return &DynamicVariable{
		value:     make(map[int]*node.INode, 0),
		valueType: t,
	}
}
func NewUnknownVariable() *DynamicVariable {
	return &DynamicVariable{
		value:     nil,
		valueType: UNKNOWN,
	}
}

func (v *DynamicVariable) GetValueString() string {
	switch v.valueType {
	case STRING:
		return v.value.(string)
	case CHARACTER:
		return v.value.(string)
	case BOOLEAN:
		if v.value.(bool) {
			return "true"
		}
		return "false"
	case NUMBER:
		return strconv.FormatFloat(v.value.(float64), 'f', -1, 64)
	case UNKNOWN:
		return ""
	default:
		panic("Called DynamicVariable@GetValueString on an unhandled value type: " + v.valueType.String())
	}
}
func (v *DynamicVariable) SetValueString(value string) {
	if v.valueType != STRING {
		panic("Called DynamicVariable@SetValueString on a non-string variable")
	}
	v.value = value
}

func (v *DynamicVariable) GetValueCharacter() string {
	if v.valueType != CHARACTER {
		panic("Called DynamicVariable@GetValueCharacter on a non-character variable")
	}
	return v.value.(string)
}
func (v *DynamicVariable) SetValueCharacter(value string) {
	if v.valueType != CHARACTER {
		panic("Called DynamicVariable@SetValueCharacter on a non-character variable")
	}
	v.value = value
}

func (v *DynamicVariable) GetValueBoolean() bool {
	if v.valueType != BOOLEAN {
		panic("Called DynamicVariable@GetValueBoolean on a non-boolean variable")
	}
	return v.value.(bool)
}
func (v *DynamicVariable) SetValueBoolean(value bool) {
	if v.valueType != BOOLEAN {
		panic("Called DynamicVariable@SetValueBoolean on a non-boolean variable")
	}
	v.value = value
}

func (v *DynamicVariable) GetValueNumber() float64 {
	if v.valueType != NUMBER {
		panic("Called DynamicVariable@GetValueNumber on a non-number variable")
	}
	return v.value.(float64)
}
func (v *DynamicVariable) SetValueNumber(value float64) {
	if v.valueType != NUMBER {
		panic("Called DynamicVariable@SetValueNumber on a non-number variable")
	}
	v.value = value
}

func (v *DynamicVariable) GetValueDictionary() map[int]*node.INode {
	if !v.valueType.IsArray() {
		panic("Called DynamicVariable@GetValueDictionary on a non-dictionary variable")
	}

	return v.value.(map[int]*node.INode)
}

func (v *DynamicVariable) GetType() VariableType {
	return v.valueType
}
