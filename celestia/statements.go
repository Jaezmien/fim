package celestia

import (
	"fmt"
	"slices"
	"strconv"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"

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
		case *nodes.ForEveryArrayStatementNode:
			if !i.Variables.Has(n.Identifier, true) {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' does not exist.", n.Identifier), i.source)
			}

			v := i.Variables.Get(n.Identifier, true)

			if v.GetType().IsArray() {
				if v.GetType().AsBaseType() != n.VariableType {
					return nil, n.ToNode().CreateError(fmt.Sprintf("Expected loop variable to be type %s, got %s", v.GetType().AsBaseType(), n.VariableType), i.source)
				}
			} else if v.GetType() == variable.STRING {
				if variable.CHARACTER != n.VariableType {
					return nil, n.ToNode().CreateError(fmt.Sprintf("Expected loop variable to be type %s, got %s", variable.CHARACTER, n.VariableType), i.source)
				}
			} else {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Expected an array variable, got type %s", v.GetType()), i.source)
			}

			if i.Variables.Has(n.VariableName, true) {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' already exists.", n.VariableName), i.source)
			}

			if v.GetType() == variable.STRING {
				for _, c := range v.GetValueString() {
					variable := &Variable{
						Name:            n.VariableName,
						DynamicVariable: variable.NewRawCharacterVariable(string(c)),
						Constant:        true,
					}

					i.Variables.PushVariable(variable, false)
					result, err := i.EvaluateStatementsNode(&n.StatementsNode)
					i.Variables.PopVariable(false)

					if result != nil || err != nil {
						return result, err
					}
				}
				break
			}

			keys := make([]int, 0, len(v.GetValueDictionary()))
			for k := range v.GetValueDictionary() {
				keys = append(keys, k)
			}
			slices.Sort(keys)

			switch v.GetType() {
			case variable.STRING_ARRAY:
				for _, k := range keys {
					v := v.GetValueDictionary()[k]

					variable := &Variable{
						Name:            n.VariableName,
						DynamicVariable: variable.NewRawStringVariable(v.GetValueString()),
						Constant:        true,
					}

					i.Variables.PushVariable(variable, false)
					result, err := i.EvaluateStatementsNode(&n.StatementsNode)
					i.Variables.PopVariable(false)

					if result != nil || err != nil {
						return result, err
					}
				}
			case variable.BOOLEAN_ARRAY:
				for _, k := range keys {
					v := v.GetValueDictionary()[k]

					variable := &Variable{
						Name:            n.VariableName,
						DynamicVariable: variable.NewBooleanVariable(v.GetValueBoolean()),
						Constant:        true,
					}

					i.Variables.PushVariable(variable, false)
					result, err := i.EvaluateStatementsNode(&n.StatementsNode)
					i.Variables.PopVariable(false)

					if result != nil || err != nil {
						return result, err
					}
				}
			case variable.NUMBER_ARRAY:
				for _, k := range keys {
					v := v.GetValueDictionary()[k]

					variable := &Variable{
						Name:            n.VariableName,
						DynamicVariable: variable.NewNumberVariable(v.GetValueNumber()),
						Constant:        true,
					}

					i.Variables.PushVariable(variable, false)
					result, err := i.EvaluateStatementsNode(&n.StatementsNode)
					i.Variables.PopVariable(false)

					if result != nil || err != nil {
						return result, err
					}
				}
			}

		case *nodes.ForEveryRangeStatementNode:
			if i.Variables.Has(n.VariableName, true) {
				return nil, n.ToNode().CreateError(fmt.Sprintf("Variable '%s' already exists.", n.VariableName), i.source)
			}

			fromRange, err := i.EvaluateValueNode(n.RangeStart, true)
			if err != nil {
				return nil, err
			}
			if fromRange.GetType() != variable.NUMBER {
				return nil, n.RangeStart.ToNode().CreateError(fmt.Sprintf("Expected a number type, got %s", fromRange.GetType()), i.source)
			}

			toRange, err := i.EvaluateValueNode(n.RangeEnd, true)
			if err != nil {
				return nil, err
			}
			if toRange.GetType() != variable.NUMBER {
				return nil, n.RangeEnd.ToNode().CreateError(fmt.Sprintf("Expected a number type, got %s", toRange.GetType()), i.source)
			}

			startValue := fromRange.GetValueNumber()
			endValue := toRange.GetValueNumber()

			if startValue == endValue {
				break
			}

			currentValue := startValue

			isForwards := true
			if endValue < startValue {
				isForwards = false
			}

			for {
				if isForwards && currentValue > endValue {
					break
				} else if !isForwards && currentValue < endValue {
					break
				}

				variable := &Variable{
					Name:            n.VariableName,
					DynamicVariable: variable.NewNumberVariable(currentValue),
					Constant:        true,
				}

				i.Variables.PushVariable(variable, false)
				result, err := i.EvaluateStatementsNode(&n.StatementsNode)
				i.Variables.PopVariable(false)

				if result != nil || err != nil {
					return result, err
				}
				
				if isForwards {
					currentValue += 1.0
				} else {
					currentValue -= 1.0
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

