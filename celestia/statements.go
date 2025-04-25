package celestia

import (
	"fmt"
	"math"
	"slices"
	"strconv"

	"git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"

	lunaErrors "git.jaezmien.com/Jaezmien/fim/luna/errors"
	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) (*variable.DynamicVariable, error) {
	newVariableCount := 0
	defer func() {
		i.Variables.PopVariableAmount(false, newVariableCount)
	}()

	for _, statement := range statements.Statements {
		switch n := statement.(type) {
		case *nodes.PrintNode:
			value, err := i.EvaluateValueNode(n.Value, true)
			if err != nil {
				return nil, err
			}
			i.Writer.Write([]byte(value.GetValueString()))

			if n.NewLine {
				i.Writer.Write([]byte("\n"))
			}
		case *nodes.PromptNode:
			if !i.Variables.Has(n.Identifier, true) {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", n.Identifier), i.source)
			}
			v := i.Variables.Get(n.Identifier, true)
			if v.Constant {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if v.DynamicVariable.GetType().IsArray() {
				return nil, n.ToNode().CreateError("Expected variable to be of non-array type", i.source)
			}

			value, err := i.EvaluateValueNode(n.Prompt, true)
			if err != nil {
				return nil, err
			}
			if value.GetType() != variable.STRING {
				return nil, n.ToNode().CreateError("Expected prompt to be of type STRING", i.source)
			}

			response, err := i.Prompt(value.GetValueString())
			if err != nil {
				return nil, err
			}

			switch v.GetType() {
			case variable.STRING:
				v.DynamicVariable.SetValueString(response)
			case variable.CHARACTER:
				value, ok := luna.AsCharacterValue(response)
				if !ok {
					return nil, n.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid character value: %s", response), i.source)
				}
				v.DynamicVariable.SetValueCharacter(value)
			case variable.BOOLEAN:
				value, ok := luna.AsBooleanValue(response)
				if !ok {
					return nil, n.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid boolean value: %s", response), i.source)
				}
				v.DynamicVariable.SetValueBoolean(value)
			case variable.NUMBER:
				value, err := strconv.ParseFloat(response, 64)
				if err != nil {
					return nil, n.Prompt.ToNode().CreateError(fmt.Sprintf("Invalid number value: %s", response), i.source)
				}
				v.DynamicVariable.SetValueNumber(value)
			}
		case *nodes.VariableDeclarationNode:
			if i.Variables.Get(n.Identifier, true) != nil {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' already exists.", n.Identifier), i.source)
			}

			value, err := i.EvaluateValueNode(n.Value, true)
			if err != nil {
				return nil, err
			}

			if value.GetType() == variable.UNKNOWN {
				if n.ValueType.IsArray() {
					value = variable.NewDictionaryVariable(n.ValueType)
				} else {
					defaultValue, ok := n.ValueType.GetDefaultValue()
					if !ok {
						panic("Intepreter@EvaluateStatementsNode could not get default value.")
					}
					value = variable.FromValueType(defaultValue, n.ValueType)
				}
			}

			if !n.ValueType.IsArray() {
				value = value.Clone()
			}

			if n.ValueType != value.GetType() {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Expected type '%s', got '%s'", n.ValueType, value.GetType()), i.source)
			}

			variable := &Variable{
				Name:            n.Identifier,
				DynamicVariable: value,
				Constant:        n.Constant,
			}

			i.Variables.PushVariable(variable, false)
			newVariableCount += 1

		case *nodes.VariableModifyNode:
			if !i.Variables.Has(n.Identifier, true) {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", n.Identifier), i.source)
			}
			v := i.Variables.Get(n.Identifier, true)

			if v.GetType().IsArray() {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Cannot modify an array."), i.source)
			}
			if v.Constant {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if n.ReinforcementType != variable.UNKNOWN {
				if v.GetType() != n.ReinforcementType {
					return nil, n.Value.ToNode().CreateError(fmt.Sprintf("Got reinforcement type '%s' when expecting type '%s'", n.ReinforcementType, v.GetType()), i.source)
				}
			}

			value, err := i.EvaluateValueNode(n.Value, true)
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
				return nil, n.ToNode().CreateError(fmt.Sprintf("Expected type '%s', got '%s'.", v.GetType(), value.GetType()), i.source)
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
		case *nodes.ArrayModifyNode:
			if !i.Variables.Has(n.Identifier, true) {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", n.Identifier), i.source)
			}
			v := i.Variables.Get(n.Identifier, true)

			if !v.GetType().IsArray() {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Invalid non-array variable."), i.source)
			}

			if n.ReinforcementType != variable.UNKNOWN {
				if v.GetType().AsBaseType() != n.ReinforcementType {
					return nil, n.Value.ToNode().CreateError(fmt.Sprintf("Got reinforcement type '%s' when expecting type '%s'", n.ReinforcementType, v.GetType()), i.source)
				}
			}

			index, err := i.EvaluateValueNode(n.Index, true)
			if err != nil {
				return nil, err
			}
			if index.GetType() != variable.NUMBER {
				return nil, n.Index.ToNode().CreateError(fmt.Sprintf("Expected a numeric index, got type %s", index.GetType()), i.source)
			}

			value, err := i.EvaluateValueNode(n.Value, true)
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
				return nil, n.Value.ToNode().CreateError(fmt.Sprintf("Cannot insert an array value"), i.source)
			}

			if v.GetType().AsBaseType() != value.GetType() {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Expected type '%s', got '%s'.", v.GetType().AsBaseType(), value.GetType()), i.source)
			}

			v.GetValueDictionary()[int(index.GetValueNumber())] = value
		case *nodes.IfStatementNode:
			for _, branch := range n.Conditions {
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

					if result != nil || err != nil {
						return result, err
					}

					break
				}
			}
		case *nodes.WhileStatementNode:
			for {
				branchCheck, err := i.EvaluateValueNode(*n.Condition, true)
				if err != nil {
					return nil, err
				}

				if branchCheck.GetType() != variable.BOOLEAN {
					return nil, n.ToNode().CreateError(fmt.Sprintf("Expected condition to result in type %s, got %s", variable.BOOLEAN, branchCheck.GetType()), i.source)
				}

				check := branchCheck.GetValueBoolean()

				if !check {
					break
				}

				result, err := i.EvaluateStatementsNode(&n.StatementsNode)

				if result != nil || err != nil {
					return result, err
				}
			}
		case *nodes.UnaryExpressionNode:
			if !i.Variables.Has(n.Identifier, true) {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", n.Identifier), i.source)
			}
			v := i.Variables.Get(n.Identifier, true)

			if v.GetType() != variable.NUMBER {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Expected a number type, got %s.", v.GetType()), i.source)
			}
			if v.Constant {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Cannot modify a constant variable."), i.source)
			}

			if n.Increment {
				v.SetValueNumber(v.GetValueNumber() + 1)
			} else {
				v.SetValueNumber(v.GetValueNumber() - 1)
			}
		case *nodes.FunctionCallNode:
			paragraphIndex := slices.IndexFunc(i.Paragraphs, func(p *Paragraph) bool { return p.Name == n.Identifier })
			if paragraphIndex == -1 {
				return nil, statement.ToNode().CreateError(fmt.Sprintf("Paragraph '%s' not found", n.Identifier), i.source)
			}

			paragraph := i.Paragraphs[paragraphIndex]

			parameters := make([]*variable.DynamicVariable, 0)
			for _, parameter := range n.Parameters {
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
		case *nodes.FunctionReturnNode:
			value, err := i.EvaluateValueNode(n.Value, true)
			return value, err
		default:
			return nil, statement.ToNode().CreateError("Unsupported statement node.", i.source)
		}
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
