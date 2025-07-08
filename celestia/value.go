package celestia

import (
	"fmt"
	"math"
	"slices"

	"git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"

	lunaErrors "git.jaezmien.com/Jaezmien/fim/luna/errors"
)

func (i *Interpreter) EvaluateValueNode(n node.DynamicNode, local bool) (*variable.DynamicVariable, error) {
	if literalNode, ok := n.(*nodes.LiteralNode); ok {
		return literalNode.DynamicVariable.Clone(), nil
	}

	if literalNode, ok := n.(*nodes.LiteralDictionaryNode); ok {
		dictionary := variable.NewDictionaryVariable(literalNode.ArrayType)
		for idx, value := range literalNode.Values {
			evaluatedValue, err := i.EvaluateValueNode(value, local)
			if err != nil {
				return nil, err
			}

			dictionary.GetValueDictionary()[idx] = evaluatedValue
		}

		return dictionary, nil
	}

	if identifierNode, ok := n.(*nodes.IdentifierNode); ok {
		if variable := i.Variables.Get(identifierNode.Identifier, local); variable != nil {
			return variable.DynamicVariable.Clone(), nil
		}

		if paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool {
			return p.Name == identifierNode.Identifier
		}); paragraphIndex != -1 {
			paragraph := i.Paragraphs[paragraphIndex]
			value, err := paragraph.Execute()
			return value, err
		}

		return nil, lunaErrors.NewParseError(fmt.Sprintf("Unknown identifier (%s)", identifierNode.Identifier), i.source, identifierNode.Start)
	}

	if callNode, ok := n.(*nodes.FunctionCallNode); ok {
		paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool { return p.Name == callNode.Identifier })
		if paragraphIndex == -1 {
			return nil, lunaErrors.NewParseError(fmt.Sprintf("Unknown paragraph (%s)", callNode.Identifier), i.source, callNode.Start)
		}

		paragraph := i.Paragraphs[paragraphIndex]

		if paragraph.FunctionNode.ReturnType == variable.UNKNOWN {
			return nil, callNode.CreateError("Tried calling a function that doesn't return a value", i.source)
		}

		parameters := make([]*variable.DynamicVariable, 0)
		for _, param := range callNode.Parameters {
			value, err := i.EvaluateValueNode(param, local)
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, value)
		}

		value, err := paragraph.Execute(parameters...)
		return value, err
	}

	if identifierNode, ok := n.(*nodes.DictionaryIdentifierNode); ok {
		v := i.Variables.Get(identifierNode.Identifier, local)
		if v == nil {
			return nil, lunaErrors.NewParseError(fmt.Sprintf("Unknown identifier (%s)", identifierNode.Identifier), i.source, identifierNode.Start)
		}
		if !v.GetType().IsArray() && v.GetType() != variable.STRING {
			return nil, lunaErrors.NewParseError(fmt.Sprintf("Invalid non-dicionary identifier (%s)", identifierNode.Identifier), i.source, identifierNode.Start)
		}

		index, _ := i.EvaluateValueNode(identifierNode.Index, local)
		if index.GetType() != variable.NUMBER {
			return nil, lunaErrors.NewParseError(fmt.Sprintf("Expected numeric index, got type %s", index.GetType()), i.source, identifierNode.Index.ToNode().Start)
		}

		indexAsInteger := int(index.GetValueNumber())

		switch v.GetType() {
		case variable.STRING:
			value := v.GetValueString()[indexAsInteger-1]
			return variable.NewRawCharacterVariable(string(value)), nil
		case variable.STRING_ARRAY:
			value := v.GetValueDictionary()[indexAsInteger]
			if value == nil {
				defaultValue, _ := variable.STRING.GetDefaultValue()
				return variable.FromValueType(defaultValue, variable.STRING), nil
			}
			return value.Clone(), nil
		case variable.BOOLEAN_ARRAY:
			value := v.GetValueDictionary()[indexAsInteger]
			if value == nil {
				defaultValue, _ := variable.BOOLEAN.GetDefaultValue()
				return variable.FromValueType(defaultValue, variable.BOOLEAN), nil
			}
			return value.Clone(), nil
		case variable.NUMBER_ARRAY:
			value := v.GetValueDictionary()[indexAsInteger]
			if value == nil {
				defaultValue, _ := variable.NUMBER.GetDefaultValue()
				return variable.FromValueType(defaultValue, variable.NUMBER), nil
			}
			return value.Clone(), nil

		}
	}

	if binaryNode, ok := n.(*nodes.BinaryExpressionNode); ok {
		left, err := i.EvaluateValueNode(binaryNode.Left, local)
		if err != nil {
			return nil, err
		}
		right, err := i.EvaluateValueNode(binaryNode.Right, local)
		if err != nil {
			return nil, err
		}

		if binaryNode.Operator == nodes.BINARYOPERATOR_ADD {
			if left.GetType() == variable.STRING || right.GetType() == variable.STRING {
				variable := variable.NewRawStringVariable(left.GetValueString() + right.GetValueString())
				return variable, nil
			}
		}

		// TODO: Add type checks for the following operators below
		// Right now we're just letting it panic if the GetValue[Type]
		// doesn't match.

		switch binaryNode.Operator {
		case nodes.BINARYOPERATOR_ADD:
			return variable.NewNumberVariable(left.GetValueNumber() + right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_SUB:
			return variable.NewNumberVariable(left.GetValueNumber() - right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_MUL:
			return variable.NewNumberVariable(left.GetValueNumber() * right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_DIV:
			return variable.NewNumberVariable(left.GetValueNumber() / right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_MOD:
			return variable.NewNumberVariable(math.Mod(left.GetValueNumber(), right.GetValueNumber())), nil

		case nodes.BINARYOPERATOR_AND:
			return variable.NewBooleanVariable(left.GetValueBoolean() && right.GetValueBoolean()), nil
		case nodes.BINARYOPERATOR_OR:
			return variable.NewBooleanVariable(left.GetValueBoolean() || right.GetValueBoolean()), nil

		case nodes.BINARYOPERATOR_GTE:
			return variable.NewBooleanVariable(left.GetValueNumber() >= right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_LTE:
			return variable.NewBooleanVariable(left.GetValueNumber() <= right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_GT:
			return variable.NewBooleanVariable(left.GetValueNumber() > right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_LT:
			return variable.NewBooleanVariable(left.GetValueNumber() < right.GetValueNumber()), nil

		case nodes.BINARYOPERATOR_NEQ:
			return variable.NewBooleanVariable(left.GetValueString() != right.GetValueString()), nil
		case nodes.BINARYOPERATOR_EQ:
			return variable.NewBooleanVariable(left.GetValueString() == right.GetValueString()), nil
		}
	}

	return nil, lunaErrors.NewParseError(fmt.Sprintf("Unsupported value node"), i.source, n.ToNode().Start)
}
