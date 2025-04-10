package celestia

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"

	lunaErrors "git.jaezmien.com/Jaezmien/fim/luna/errors"
	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) (*vartype.DynamicVariable, error) {
	for _, statement := range statements.Statements {
		if statement.Type() == node.TYPE_PRINT {
			printNode := statement.(*nodes.PrintNode)

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
		if statement.Type() == node.TYPE_PROMPT {
			promptNode := statement.(*nodes.PromptNode)

			if !i.Variables.Has(promptNode.Identifier, true) {
				return nil, promptNode.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", promptNode.Identifier), i.source)
			}
			variable := i.Variables.Get(promptNode.Identifier, true)
			if variable.Constant {
				return nil, promptNode.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if variable.DynamicVariable.GetType().IsArray() {
				return nil, promptNode.ToNode().CreateError("Expected variable to be of non-array type", i.source)
			}

			value, err := i.EvaluateValueNode(promptNode.Prompt, true)
			if err != nil {
				return nil, err
			}
			if value.GetType() != vartype.STRING {
				return nil, promptNode.ToNode().CreateError("Expected prompt to be of type STRING", i.source)
			}

			response, err := i.Prompt(value.GetValueString())
			if err != nil {
				return nil, err
			}

			switch variable.GetType() {
			case vartype.STRING:
				variable.DynamicVariable.SetValueString(response)
				break
			case vartype.CHARACTER:
				value, ok := luna.AsCharacterValue(response)
				if !ok {
					return nil, promptNode.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid character value: %s", response), i.source)
				}
				variable.DynamicVariable.SetValueCharacter(value)
				break
			case vartype.BOOLEAN:
				value, ok := luna.AsBooleanValue(response)
				if !ok {
					return nil, promptNode.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid boolean value: %s", response), i.source)
				}
				variable.DynamicVariable.SetValueBoolean(value)
				break
			case vartype.NUMBER:
				value, err := strconv.ParseFloat(response, 64)
				if err != nil {
					return nil, promptNode.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid number value: %s", response), i.source)
				}
				variable.DynamicVariable.SetValueNumber(value)
				break
			}

			continue
		}
		if statement.Type() == node.TYPE_VARIABLE_DECLARATION {
			variableNode := statement.(*nodes.VariableDeclarationNode)

			value, err := i.EvaluateValueNode(variableNode.Value, true)
			if err != nil {
				return nil, err
			}

			if value.GetType() == vartype.UNKNOWN {
				if variableNode.ValueType.IsArray() {
					value = vartype.NewDictionaryVariable(variableNode.ValueType)
				} else {
					defaultValue, ok := variableNode.ValueType.GetDefaultValue()
					if !ok {
						panic("Intepreter@EvaluateStatementsNode could not get default value.")
					}
					value = vartype.FromValueType(defaultValue, variableNode.ValueType)
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
		if statement.Type() == node.TYPE_VARIABLE_MODIFY {
			modifyNode := statement.(*nodes.VariableModifyNode)

			if !i.Variables.Has(modifyNode.Identifier, true) {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", modifyNode.Identifier), i.source)
			}
			variable := i.Variables.Get(modifyNode.Identifier, true)

			if variable.GetType().IsArray() {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Cannot modify an array."), i.source)
			}
			if variable.Constant {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if modifyNode.ReinforcementType != vartype.UNKNOWN {
				if variable.GetType() != modifyNode.ReinforcementType {
					return nil, modifyNode.Value.ToNode().CreateError(fmt.Sprintf("Got reinforcement type '%s' when expecting type '%s'", modifyNode.ReinforcementType, variable.GetType()), i.source)
				}
			}

			value, err := i.EvaluateValueNode(modifyNode.Value, true)
			if err != nil {
				return nil, err
			}

			if value.GetType() == vartype.UNKNOWN {
				defaultValue, ok := variable.GetType().GetDefaultValue()
				if !ok {
					panic("Intepreter@EvaluateStatementsNode could not get default value.")
				}
				value = vartype.FromValueType(defaultValue, variable.GetType())
			}

			if variable.GetType() != value.GetType() && (variable.DynamicVariable.GetType() != vartype.STRING && !value.GetType().IsArray()) {
				return nil, modifyNode.ToNode().CreateError(fmt.Sprintf("Expected type '%s', got '%s'.", variable.GetType(), value.GetType()), i.source)
			}

			switch variable.GetType() {
			case vartype.STRING:
				variable.SetValueString(value.GetValueString())
			case vartype.CHARACTER:
				variable.SetValueCharacter(value.GetValueCharacter())
			case vartype.BOOLEAN:
				variable.SetValueBoolean(value.GetValueBoolean())
			case vartype.NUMBER:
				variable.SetValueNumber(value.GetValueNumber())
			}

			continue
		}
		if statement.Type() == node.TYPE_FUNCTION_CALL {
			callNode := statement.(*nodes.FunctionCallNode)

			paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool { return p.Name == callNode.Identifier })
			if paragraphIndex == -1 {
				return nil, statement.ToNode().CreateError(fmt.Sprintf("Paragraph '%s' not found", callNode.Identifier), i.source)
			}

			paragraph := i.Paragraphs[paragraphIndex]

			parameters := make([]*vartype.DynamicVariable, 0)
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

		if statement.Type() == node.TYPE_FUNCTION_RETURN {
			returnNode := statement.(*nodes.FunctionReturnNode)

			value, err := i.EvaluateValueNode(returnNode.Value, true)

			return value, err
		}

		return nil, statement.ToNode().CreateError(fmt.Sprintf("Unsupported statement node: %s", statement.Type()), i.source)
	}

	return nil, nil
}

func (i *Interpreter) EvaluateValueNode(n node.INode, local bool) (*vartype.DynamicVariable, error) {
	if n.Type() == node.TYPE_LITERAL {
		literalNode := n.(*nodes.LiteralNode)

		return literalNode.DynamicVariable, nil
	}

	if n.Type() == node.TYPE_LITERAL_DICTIONARY {
		literalNode := n.(*nodes.LiteralDictionaryNode)

		return literalNode.DynamicVariable, nil
	}

	if n.Type() == node.TYPE_IDENTIFIER {
		identifierNode := n.(*nodes.IdentifierNode)

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

	if n.Type() == node.TYPE_FUNCTION_CALL {
		callNode := n.(*nodes.FunctionCallNode)

		paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool { return p.Name == callNode.Identifier })
		if paragraphIndex != -1 {
			paragraph := i.Paragraphs[paragraphIndex]
			value, err := paragraph.Execute()
			return value, err
		}

		return nil, lunaErrors.NewParseError(fmt.Sprintf("Unknown paragraph (%s)", callNode.Identifier), i.source, callNode.Start)
	}

	if n.Type() == node.TYPE_IDENTIFIER_DICTIONARY {
		identifierNode := n.(*nodes.DictionaryIdentifierNode)

		variable := i.Variables.Get(identifierNode.Identifier, local)
		if variable == nil {
			return nil, lunaErrors.NewParseError(fmt.Sprintf("Unknown identifier (%s)", identifierNode.Identifier), i.source, identifierNode.Start)
		}
		if !variable.GetType().IsArray() && variable.GetType() != vartype.STRING {
			return nil, lunaErrors.NewParseError(fmt.Sprintf("Invalid non-dicionary identifier (%s)", identifierNode.Identifier), i.source, identifierNode.Start)
		}

		index, _ := i.EvaluateValueNode(identifierNode.Index, local)
		if index.GetType() != vartype.NUMBER {
			return nil, lunaErrors.NewParseError(fmt.Sprintf("Expected numeric index, got type %s", index.GetType()), i.source, identifierNode.Index.ToNode().Start)
		}

		indexAsInteger := int(index.GetValueNumber())

		switch variable.GetType() {
		case vartype.STRING:
			value := variable.GetValueString()[indexAsInteger-1]
			return vartype.NewRawCharacterVariable(string(value)), nil
		case vartype.STRING_ARRAY:
			preValue := variable.GetValueDictionary()[indexAsInteger]
			if preValue == nil {
				defaultValue, _ := vartype.STRING.GetDefaultValue()
				return vartype.FromValueType(defaultValue, vartype.STRING), nil
			}
			value, err := i.EvaluateValueNode(*preValue, local)
			if err != nil {
				return nil, err
			}
			if value.GetType() != vartype.STRING {
				return nil, errors.New(fmt.Sprintf("Expected string"))
			}
			return value, nil
		case vartype.BOOLEAN_ARRAY:
			preValue := variable.GetValueDictionary()[indexAsInteger]
			if preValue == nil {
				defaultValue, _ := vartype.BOOLEAN.GetDefaultValue()
				return vartype.FromValueType(defaultValue, vartype.BOOLEAN), nil
			}
			value, err := i.EvaluateValueNode(*preValue, local)
			if err != nil {
				return nil, err
			}
			if value.GetType() != vartype.BOOLEAN {
				return nil, errors.New(fmt.Sprintf("Expected boolean"))
			}
			return value, nil
		case vartype.NUMBER_ARRAY:
			preValue := variable.GetValueDictionary()[indexAsInteger]
			if preValue == nil {
				defaultValue, _ := vartype.NUMBER.GetDefaultValue()
				return vartype.FromValueType(defaultValue, vartype.NUMBER), nil
			}
			value, err := i.EvaluateValueNode(*preValue, local)
			if err != nil {
				return nil, err
			}
			if value.GetType() != vartype.NUMBER {
				return nil, errors.New(fmt.Sprintf("Expected number"))
			}
			return value, nil

		}
	}

	if n.Type() == node.TYPE_BINARYEXPRESSION {
		binaryNode := n.(*nodes.BinaryExpressionNode)

		left, err := i.EvaluateValueNode(binaryNode.Left, local)
		if err != nil {
			return nil, err
		}
		right, err := i.EvaluateValueNode(binaryNode.Right, local)
		if err != nil {
			return nil, err
		}

		if binaryNode.Operator == nodes.BINARYOPERATOR_ADD {
			if left.GetType() == vartype.STRING || right.GetType() == vartype.STRING {
				variable := vartype.NewRawStringVariable(left.GetValueString() + right.GetValueString())
				return variable, nil
			}
		}

		// TODO: Add type checks

		switch binaryNode.Operator {
		case nodes.BINARYOPERATOR_ADD:
			return vartype.NewNumberVariable(left.GetValueNumber() + right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_SUB:
			return vartype.NewNumberVariable(left.GetValueNumber() - right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_MUL:
			return vartype.NewNumberVariable(left.GetValueNumber() * right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_DIV:
			return vartype.NewNumberVariable(left.GetValueNumber() / right.GetValueNumber()), nil

		case nodes.BINARYOPERATOR_AND:
			return vartype.NewBooleanVariable(left.GetValueBoolean() && right.GetValueBoolean()), nil
		case nodes.BINARYOPERATOR_OR:
			return vartype.NewBooleanVariable(left.GetValueBoolean() || right.GetValueBoolean()), nil

		case nodes.BINARYOPERATOR_GTE:
			return vartype.NewBooleanVariable(left.GetValueNumber() >= right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_LTE:
			return vartype.NewBooleanVariable(left.GetValueNumber() <= right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_GT:
			return vartype.NewBooleanVariable(left.GetValueNumber() > right.GetValueNumber()), nil
		case nodes.BINARYOPERATOR_LT:
			return vartype.NewBooleanVariable(left.GetValueNumber() < right.GetValueNumber()), nil

		case nodes.BINARYOPERATOR_NEQ:
			return vartype.NewBooleanVariable(left.GetValueString() != right.GetValueString()), nil
		case nodes.BINARYOPERATOR_EQ:
			return vartype.NewBooleanVariable(left.GetValueString() == right.GetValueString()), nil
		}
	}

	return nil, lunaErrors.NewParseError(fmt.Sprintf("Unsupported value node"), i.source, n.ToNode().Start)
}
