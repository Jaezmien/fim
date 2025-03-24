package vartype

import (
	"slices"
	"strconv"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/spike/utilities"
)

type DynamicVariable struct {
	valueString string
	valueNumber float64
	valueBoolean bool
	valueCharacter string
	valueDictionary map[int]DynamicVariable

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
		valueString: value,
		valueType: STRING,
	}
}
func NewCharacterVariable(value string) *DynamicVariable {
	value = value[1: len(value)-1] 

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
		valueCharacter: value,
		valueType: CHARACTER,
	}
}
func NewNumberVariable(value float64) *DynamicVariable {
	return &DynamicVariable{
		valueNumber: value,
		valueType: NUMBER,
	}
}
func NewBooleanVariable(value bool) *DynamicVariable {
	return &DynamicVariable{
		valueBoolean: value,
		valueType: BOOLEAN,
	}
}
func NewDictionaryVariable(t VariableType) *DynamicVariable {
	return &DynamicVariable{
		valueDictionary: make(map[int]DynamicVariable),
		valueType: t,
	}
}

func (v *DynamicVariable) GetValueString() string {
	switch v.valueType {
	case STRING:
		return v.valueString
	case CHARACTER:
		return v.valueCharacter
	case BOOLEAN:
		if v.valueBoolean {
			return "true"
		}
		return "false"
	case NUMBER:
		return strconv.FormatFloat(v.valueNumber, 'f', -1, 64)
	default:
		if !v.valueType.IsArray() {
			panic("Called DynamicVariable@GetValueString on an unhandled value type: " + v.valueType.String())
		}

		s := make([]string, len(v.valueDictionary))
		for _, element := range v.valueDictionary {
			s = append(s, element.GetValueString())
		}
		return strings.Join(s, ",")
	}
}
func (v *DynamicVariable) SetValueString(value string) {
	if v.valueType != STRING {
		panic("Called DynamicVariable@SetValueString on a non-string variable")
	}
	v.valueString = value
}

func (v *DynamicVariable) GetValueCharacter() string {
	if v.valueType != CHARACTER {
		panic("Called DynamicVariable@GetValueCharacter on a non-character variable")
	}
	return v.valueCharacter
}
func (v *DynamicVariable) SetValueCharacter(value string) {
	if v.valueType != CHARACTER {
		panic("Called DynamicVariable@SetValueCharacter on a non-character variable")
	}
	v.valueCharacter = value
}

func (v *DynamicVariable) GetValueBoolean() bool {
	if v.valueType != BOOLEAN {
		panic("Called DynamicVariable@GetValueBoolean on a non-boolean variable")
	}
	return v.valueBoolean
}
func (v *DynamicVariable) SetValueBoolean(value bool) {
	if v.valueType != BOOLEAN {
		panic("Called DynamicVariable@SetValueBoolean on a non-boolean variable")
	}
	v.valueBoolean = value
}

func (v *DynamicVariable) GetValueNumber() float64 {
	if v.valueType != NUMBER {
		panic("Called DynamicVariable@GetValueNumber on a non-number variable")
	}
	return v.valueNumber
}
func (v *DynamicVariable) SetValueNumber(value float64) {
	if v.valueType != NUMBER {
		panic("Called DynamicVariable@SetValueNumber on a non-number variable")
	}
	v.valueNumber = value
}

func (v *DynamicVariable) GetValueDictionary() map[int]DynamicVariable {
	if v.valueType != NUMBER {
		panic("Called DynamicVariable@GetValueDictionary on a non-dictionary variable")
	}

	return v.valueDictionary
}

func (v *DynamicVariable) GetType() VariableType {
	return v.valueType
}
