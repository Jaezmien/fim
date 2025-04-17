package celestia

import (
	"fmt"
	"slices"
	"strconv"

	"git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"

	lunaErrors "git.jaezmien.com/Jaezmien/fim/luna/errors"
	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) (*variable.DynamicVariable, error) {
	for _, statement := range statements.Statements {
		if printNode, ok := statement.(*nodes.PrintNode); ok {
			value, err := i.EvaluateValueNode(printNode.Value, true)
			if err != nil {
				return nil, err
			}
			i.Writer.Write([]byte(value.GetValueString()))

			if printNode.NewLine {
				i.Writer.Write([]byte("\n"))
			}

			continue
		}
		if promptNode, ok := statement.(*nodes.PromptNode); ok {
			if !i.Variables.Has(promptNode.Identifier, true) {
				return nil, promptNode.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", promptNode.Identifier), i.source)
			}
			v := i.Variables.Get(promptNode.Identifier, true)
			if v.Constant {
				return nil, promptNode.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if v.DynamicVariable.GetType().IsArray() {
				return nil, promptNode.ToNode().CreateError("Expected variable to be of non-array type", i.source)
			}

			value, err := i.EvaluateValueNode(promptNode.Prompt, true)
			if err != nil {
				return nil, err
			}
			if value.GetType() != variable.STRING {
				return nil, promptNode.ToNode().CreateError("Expected prompt to be of type STRING", i.source)
			}

			response, err := i.Prompt(value.GetValueString())
			if err != nil {
				return nil, err
			}

			switch v.GetType() {
			case variable.STRING:
				v.DynamicVariable.SetValueString(response)
				break
			case variable.CHARACTER:
				value, ok := luna.AsCharacterValue(response)
				if !ok {
					return nil, promptNode.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid character value: %s", response), i.source)
				}
				v.DynamicVariable.SetValueCharacter(value)
				break
			case variable.BOOLEAN:
				value, ok := luna.AsBooleanValue(response)
				if !ok {
					return nil, promptNode.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid boolean value: %s", response), i.source)
				}
				v.DynamicVariable.SetValueBoolean(value)
				break
			case variable.NUMBER:
				value, err := strconv.ParseFloat(response, 64)
				if err != nil {
					return nil, promptNode.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid number value: %s", response), i.source)
				}
				v.DynamicVariable.SetValueNumber(value)
				break
			}

			continue
		}
		if variableNode, ok := statement.(*nodes.VariableDeclarationNode); ok {
			value, err := i.EvaluateValueNode(variableNode.Value, true)
			if err != nil {
				return nil, err
			}

			if value.GetType() == variable.UNKNOWN {
				if variableNode.ValueType.IsArray() {
					value = variable.NewDictionaryVariable(variableNode.ValueType)
				} else {
					defaultValue, ok := variableNode.ValueType.GetDefaultValue()
					if !ok {
						panic("Intepreter@EvaluateStatementsNode could not get default value.")
					}
					value = variable.FromValueType(defaultValue, variableNode.ValueType)
				}
			}

			if variableNode.ValueType != value.GetType() {
				return nil, variableNode.ToNode().CreateError(fmt.Sprintf("Expected type '%s', got '%s'", variableNode.ValueType, value.GetType()), i.source)
			}

			variable := &Variable{
				Name:            variableNode.Identifier,
				DynamicVariable: value,
				Constant:        variableNode.Constant,
			}

			i.Variables.PushVariable(variable, false)

			continue
		}
		if modifyNode, ok := statement.(*nodes.VariableModifyNode); ok {
			if !i.Variables.Has(modifyNode.Identifier, true) {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", modifyNode.Identifier), i.source)
			}
			v := i.Variables.Get(modifyNode.Identifier, true)

			if v.GetType().IsArray() {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Cannot modify an array."), i.source)
			}
			if v.Constant {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if modifyNode.ReinforcementType != variable.UNKNOWN {
				if v.GetType() != modifyNode.ReinforcementType {
					return nil, modifyNode.Value.ToNode().CreateError(fmt.Sprintf("Got reinforcement type '%s' when expecting type '%s'", modifyNode.ReinforcementType, v.GetType()), i.source)
				}
			}

			value, err := i.EvaluateValueNode(modifyNode.Value, true)
			if err != nil {
				return nil, err
			}

			if value.GetType() == variable.UNKNOWN {
				defaultValue, ok := v.GetType().GetDefaultValue()
				if !ok {
					panic("Intepreter@EvaluateStatementsNode could not get default value.")
				}
				value = variable.FromValueType(defaultValue, v.GetType())
			}

			if v.GetType() != value.GetType() && (v.DynamicVariable.GetType() != variable.STRING && !value.GetType().IsArray()) {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Expected type '%s', got '%s'.", v.GetType(), value.GetType()), i.source)
			}

			switch v.GetType() {
			case variable.STRING:
				v.SetValueString(value.GetValueString())
			case variable.CHARACTER:
				v.SetValueCharacter(value.GetValueCharacter())
			case variable.BOOLEAN:
				v.SetValueBoolean(value.GetValueBoolean())
			case variable.NUMBER:
				v.SetValueNumber(value.GetValueNumber())
			}

			continue
		}

		if modifyNode, ok := statement.(*nodes.ArrayModifyNode); ok {
			if !i.Variables.Has(modifyNode.Identifier, true) {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", modifyNode.Identifier), i.source)
			}
			v := i.Variables.Get(modifyNode.Identifier, true)

			if !v.GetType().IsArray() {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Invalid non-array variable."), i.source)
			}

			if modifyNode.ReinforcementType != variable.UNKNOWN {
				if v.GetType().AsBaseType() != modifyNode.ReinforcementType {
					return nil, modifyNode.Value.ToNode().CreateError(fmt.Sprintf("Got reinforcement type '%s' when expecting type '%s'", modifyNode.ReinforcementType, v.GetType()), i.source)
				}
			}

			index, err := i.EvaluateValueNode(modifyNode.Index, true)
			if err != nil {
				return nil, err
			}
			if index.GetType() != variable.NUMBER {
				return nil, modifyNode.Index.ToNode().CreateError(fmt.Sprintf("Expected a numeric index, got type %s", index.GetType()), i.source)
			}

			value, err := i.EvaluateValueNode(modifyNode.Value, true)
			if err != nil {
				return nil, err
			}

			if value.GetType() == variable.UNKNOWN {
				defaultValue, ok := v.GetType().AsBaseType().GetDefaultValue()
				if !ok {
					panic("Intepreter@EvaluateStatementsNode could not get default value.")
				}
				value = variable.FromValueType(defaultValue, v.GetType())
			}
			if value.GetType().IsArray() {
				return nil, modifyNode.Value.ToNode().CreateError(fmt.Sprintf("Cannot insert an array value"), i.source)
			}

			if v.GetType().AsBaseType() != value.GetType() {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Expected type '%s', got '%s'.", v.GetType().AsBaseType(), value.GetType()), i.source)
			}

			v.GetValueDictionary()[int(index.GetValueNumber())] = value

			continue
		}

		if ifNode, ok := statement.(*nodes.IfStatementNode); ok {
			for _, branch := range ifNode.Conditions {
				check := true

				if branch.Condition != nil {
					branchCheck, err := i.EvaluateValueNode(*branch.Condition, true)
					if err != nil {
						return nil, err
					}

					if branchCheck.GetType() != variable.BOOLEAN {
						return nil, branch.ToNode().CreateError(fmt.Sprintf("Expected condition to result in type %s, got %s", variable.BOOLEAN, branchCheck.GetType()), i.source)
					}

					check = branchCheck.GetValueBoolean()
				}

				if check {
					result, err := i.EvaluateStatementsNode(&branch.StatementsNode)
					return result, err
				}
			}

			continue
		}

		if unaryNode, ok := statement.(*nodes.UnaryExpressionNode); ok {
			if !i.Variables.Has(unaryNode.Identifier, true) {
				return nil, unaryNode.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", unaryNode.Identifier), i.source)
			}
			v := i.Variables.Get(unaryNode.Identifier, true)

			if v.GetType() != variable.NUMBER {
				return nil, unaryNode.ToNode().CreateError(fmt.Sprintf("Expected a number type, got %s.", v.GetType()), i.source)
			}
			if v.Constant {
				return nil, unaryNode.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if unaryNode.Increment {
				v.SetValueNumber(v.GetValueNumber() + 1)
			} else {
				v.SetValueNumber(v.GetValueNumber() - 1)
			}
			continue
		}

		if callNode, ok := statement.(*nodes.FunctionCallNode); ok {
			paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool { return p.Name == callNode.Identifier })
			if paragraphIndex == -1 {
				return nil, statement.ToNode().CreateError(fmt.Sprintf("Paragraph '%s' not found", callNode.Identifier), i.source)
			}

			paragraph := i.Paragraphs[paragraphIndex]

			parameters := make([]*variable.DynamicVariable, 0)
			for _, parameter := range callNode.Parameters {
				valueNode, err := i.EvaluateValueNode(parameter, true)
				if err != nil {
					return nil, err
				}
				parameters = append(parameters, valueNode)
			}

			_, err := paragraph.Execute(parameters...)
			if err != nil {
				return nil, err
			}

			continue
		}

		if returnNode, ok := statement.(*nodes.FunctionReturnNode); ok {
			value, err := i.EvaluateValueNode(returnNode.Value, true)

			return value, err
		}

		return nil, statement.ToNode().CreateError("Unsupported statement node.", i.source)
	}

	return nil, nil
}

func (i *Interpreter) EvaluateValueNode(n node.DynamicNode, local bool) (*variable.DynamicVariable, error) {
	if literalNode, ok := n.(*nodes.LiteralNode); ok {
		return literalNode.DynamicVariable, nil
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
			return variable.DynamicVariable, nil
		}

		if paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool { return p.Name == identifierNode.Identifier }); paragraphIndex != -1 {
			paragraph := i.Paragraphs[paragraphIndex]
			value, err := paragraph.Execute()
			return value, err
		}

		return nil, lunaErrors.NewParseError(fmt.Sprintf("Unknown identifier (%s)", identifierNode.Identifier), i.source, identifierNode.Start)
	}

	if callNode, ok := n.(*nodes.FunctionCallNode); ok {
		paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool { return p.Name == callNode.Identifier })
		if paragraphIndex != -1 {
			paragraph := i.Paragraphs[paragraphIndex]
			value, err := paragraph.Execute()
			return value, err
		}

		return nil, lunaErrors.NewParseError(fmt.Sprintf("Unknown paragraph (%s)", callNode.Identifier), i.source, callNode.Start)
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
			return value, nil
		case variable.BOOLEAN_ARRAY:
			value := v.GetValueDictionary()[indexAsInteger]
			if value == nil {
				defaultValue, _ := variable.BOOLEAN.GetDefaultValue()
				return variable.FromValueType(defaultValue, variable.BOOLEAN), nil
			}
			return value, nil
		case variable.NUMBER_ARRAY:
			value := v.GetValueDictionary()[indexAsInteger]
			if value == nil {
				defaultValue, _ := variable.NUMBER.GetDefaultValue()
				return variable.FromValueType(defaultValue, variable.NUMBER), nil
			}
			return value, nil

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
