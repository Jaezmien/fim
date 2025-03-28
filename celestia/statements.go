package celestia

import (
	"errors"
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"

	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) error {
	for _, statement := range statements.Statements {
		if statement.Type() == node.TYPE_PRINT {
			printNode := statement.(*nodes.PrintNode)

			value, err := i.EvaluateValueNode(printNode.Value, true)
			if err != nil {
				return err
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
				return i.CreateErrorFromNode(promptNode.ToNode(), fmt.Sprintf("Variable '%s' does not exist.", promptNode.Identifier))
			}
			variable := i.Variables.Get(promptNode.Identifier, true)
			if variable.Constant {
				return i.CreateErrorFromNode(promptNode.ToNode(), fmt.Sprintf("Cannot modify a constant variable."))
			}
			if variable.DynamicVariable.GetType() != vartype.STRING {
				return i.CreateErrorFromNode(promptNode.ToNode(), "Expected variable to be of type STRING")
			}

			value, err := i.EvaluateValueNode(promptNode.Prompt, true)
			if err != nil {
				return err
			}
			if value.GetType() != vartype.STRING {
				return i.CreateErrorFromNode(promptNode.ToNode(), "Expected prompt to be of type STRING")
			}

			response, err := i.Prompt(value.GetValueString())
			if err != nil {
				return err
			}
			variable.DynamicVariable.SetValueString(response)

			continue
		}
		if statement.Type() == node.TYPE_VARIABLE_DECLARATION {
			variableNode := statement.(*nodes.VariableDeclarationNode)

			value, err := i.EvaluateValueNode(variableNode.Value, true)
			if err != nil {
				return err
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
				return i.CreateErrorFromNode(variableNode.ToNode(), fmt.Sprintf("Expected type '%s', got '%s'", variableNode.ValueType, value.GetType()))
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
				return i.CreateErrorFromNode(modifyNode.ToNode(), fmt.Sprintf("Variable '%s' does not exist.", modifyNode.Identifier))
			}
			variable := i.Variables.Get(modifyNode.Identifier, true)

			if variable.GetType().IsArray() {
				return i.CreateErrorFromNode(modifyNode.ToNode(), fmt.Sprintf("Cannot modify an array."))
			}
			if variable.Constant {
				return i.CreateErrorFromNode(modifyNode.ToNode(), fmt.Sprintf("Cannot modify a constant variable."))
			}

			if modifyNode.ReinforcementType != vartype.UNKNOWN {
				if variable.GetType() != modifyNode.ReinforcementType {
					return i.CreateErrorFromNode(
						modifyNode.Value.ToNode(),
						fmt.Sprintf("Got reinforcement type '%s' when expecting type '%s'", modifyNode.ReinforcementType, variable.GetType()),
					)
				}
			}

			value, err := i.EvaluateValueNode(modifyNode.Value, true)
			if err != nil {
				return err
			}

			if value.GetType() == vartype.UNKNOWN {
				defaultValue, ok := variable.GetType().GetDefaultValue()
				if !ok {
					panic("Intepreter@EvaluateStatementsNode could not get default value.")
				}
				value = vartype.FromValueType(defaultValue, variable.GetType())
			}

			if variable.GetType() != value.GetType() && (variable.DynamicVariable.GetType() != vartype.STRING && !value.GetType().IsArray()) {
				return i.CreateErrorFromNode(modifyNode.ToNode(), fmt.Sprintf("Expected type '%s', got '%s'.", variable.GetType(), value.GetType()))
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

		return i.CreateErrorFromNode(statement.ToNode(), fmt.Sprintf("Unsupported statement node: %s", statement.Type()))
	}

	return nil
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

		// TODO: Check for paragraphs
		pair := luna.GetErrorIndexPair(i.source, identifierNode.Start)
		return nil, errors.New(fmt.Sprintf("Unknown identifier at line %d:%d", pair.Line, pair.Column))
	}

	if n.Type() == node.TYPE_IDENTIFIER_DICTIONARY {
		identifierNode := n.(*nodes.DictionaryIdentifierNode)

		variable := i.Variables.Get(identifierNode.Identifier, local)
		if variable == nil {
			pair := luna.GetErrorIndexPair(i.source, identifierNode.Start)
			return nil, errors.New(fmt.Sprintf("Unknown identifier at line %d:%d", pair.Line, pair.Column))
		}
		if !variable.GetType().IsArray() && variable.GetType() != vartype.STRING {
			pair := luna.GetErrorIndexPair(i.source, identifierNode.Start)
			return nil, errors.New(fmt.Sprintf("Invalid non-dictionary identifier at line %d:%d", pair.Line, pair.Column))
		}

		index, _ := i.EvaluateValueNode(identifierNode.Index, local)
		if index.GetType() != vartype.NUMBER {
			pair := luna.GetErrorIndexPair(i.source, identifierNode.Index.ToNode().Start)
			return nil, errors.New(fmt.Sprintf("Expected numeric index at line %d:%d", pair.Line, pair.Column))
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

	pair := luna.GetErrorIndexPair(i.source, n.ToNode().Start)
	return nil, errors.New(fmt.Sprintf("Unsupported value node at line %d:%d", pair.Line, pair.Column))
}
