package celestia

import (
	"errors"
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"

	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) (error) {
	for _, statement := range statements.Statements {
		if statement.Type() == nodes.TYPE_PRINT {
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
		if statement.Type() == nodes.TYPE_PROMPT {
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
		if statement.Type() == nodes.TYPE_VARIABLE_DECLARATION {
			variableNode := statement.(*nodes.VariableDeclarationNode)

			value, err := i.EvaluateValueNode(variableNode.Value, true)
			if err != nil {
				return err
			}

			variable := &Variable{
				Name: variableNode.Identifier,
				DynamicVariable: value,
				Constant: variableNode.Constant,
			}
			
			i.Variables.PushVariable(variable, false)

			continue
		}
		if statement.Type() == nodes.TYPE_VARIABLE_MODIFY {
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

			value, err := i.EvaluateValueNode(modifyNode.Value, true)
			if err != nil {
				return err
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

func (i *Interpreter) EvaluateValueNode(node nodes.INode, local bool) (*vartype.DynamicVariable, error) {
	if node.Type() == nodes.TYPE_LITERAL {
		literalNode := node.(*nodes.LiteralNode)

		return literalNode.DynamicVariable, nil
	}
	
	if node.Type() == nodes.TYPE_IDENTIFIER {
		identifierNode := node.(*nodes.IdentifierNode)

		if variable := i.Variables.Get(identifierNode.Identifier, local); variable != nil {
			return variable.DynamicVariable, nil
		}

		// TODO: Check for paragraphs
		pair := luna.GetErrorIndexPair(i.source, identifierNode.Start)
		return nil, errors.New(fmt.Sprintf("Unknown identifier at line %d:%d", pair.Line, pair.Column))
	}

	if node.Type() == nodes.TYPE_BINARYEXPRESSION {
		binaryNode := node.(*nodes.BinaryExpressionNode)

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

	pair := luna.GetErrorIndexPair(i.source, node.ToNode().Start)
	return nil, errors.New(fmt.Sprintf("Unsupported value node at line %d:%d", pair.Line, pair.Column))
}
