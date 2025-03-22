package celestia

import (
	"errors"
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/celestia/utilities"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"

	spikeUtils "git.jaezmien.com/Jaezmien/fim/spike/utilities"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) (error) {
	for _, statement := range statements.Statements {
		if statement.Type() == nodes.TYPE_PRINT {
			printNode := statement.(*nodes.PrintNode)

			value, _, err := i.EvaluateValueNode(printNode.Value, true)
			if err != nil {
				return err
			}
			i.Writer.Write([]byte(value))

			if printNode.NewLine {
				i.Writer.Write([]byte("\n"))
			}

			continue
		}
		if statement.Type() == nodes.TYPE_VARIABLE_DECLARATION {
			variableNode := statement.(*nodes.VariableDeclarationNode)

			value, _, err := i.EvaluateValueNode(variableNode.Value, true)
			if err != nil {
				return err
			}

			variable := &Variable{
				Name: variableNode.Identifier,
				Value: value,
				ValueType: variableNode.ValueType,
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
			
			if variable.Constant {
				return i.CreateErrorFromNode(modifyNode.ToNode(), fmt.Sprintf("Cannot modify a constant variable."))
			}

			value, valueType, err := i.EvaluateValueNode(modifyNode.Value, true)
			if err != nil {
				return err
			}

			if variable.ValueType != valueType && (variable.ValueType != vartype.STRING && !valueType.IsArray()) {
				return i.CreateErrorFromNode(modifyNode.ToNode(), fmt.Sprintf("Expected type '%s', got '%s'.", variable.ValueType, valueType))
			}

			variable.Value = value

			continue
		}

		return i.CreateErrorFromNode(statement.ToNode(), fmt.Sprintf("Unsupported statement node: %s", statement.Type()))
	}

	return nil
}

func (i *Interpreter) EvaluateValueNode(node nodes.INode, local bool) (string, vartype.VariableType, error) {
	if node.Type() == nodes.TYPE_LITERAL {
		literalNode := node.(*nodes.LiteralNode)

		if literalNode.ValueType == vartype.BOOLEAN {
			return utilities.BoolAsString(literalNode.GetValueBoolean()), literalNode.ValueType, nil
		}
		if literalNode.ValueType == vartype.STRING {
			return literalNode.GetValueString(), literalNode.ValueType, nil
		}
		if literalNode.ValueType == vartype.CHARACTER {
			return literalNode.GetValueCharacter(), literalNode.ValueType, nil
		}
		if literalNode.ValueType == vartype.NUMBER {
			return utilities.FloatAsString(literalNode.GetValueNumber()), literalNode.ValueType, nil
		}
	}
	
	if node.Type() == nodes.TYPE_IDENTIFIER {
		identifierNode := node.(*nodes.IdentifierNode)

		if variable := i.Variables.Get(identifierNode.Identifier, local); variable != nil {
			return variable.Value, variable.ValueType, nil
		}

		// TODO: Check for paragraphs
		pair := spikeUtils.GetErrorIndexPair(i.source, identifierNode.Start)
		return "", vartype.UNKNOWN, errors.New(fmt.Sprintf("Unknown identifier at line %d:%d", pair.Line, pair.Column))
	}

	if node.Type() == nodes.TYPE_BINARYEXPRESSION {
		binaryNode := node.(*nodes.BinaryExpressionNode)

		left, leftType, err := i.EvaluateValueNode(binaryNode.Left, local)
		if err != nil {
			return "", vartype.UNKNOWN, err
		}
		right, rightType, err := i.EvaluateValueNode(binaryNode.Right, local)
		if err != nil {
			return "", vartype.UNKNOWN, err
		}

		if binaryNode.Operator == nodes.BINARYOPERATOR_ADD {
			if leftType == vartype.STRING || rightType == vartype.STRING {
				return (left + right), vartype.STRING, nil
			}
		}

		switch binaryNode.Operator {
		case nodes.BINARYOPERATOR_ADD:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat + rightFloat), vartype.NUMBER, nil
		case nodes.BINARYOPERATOR_SUB:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat - rightFloat), vartype.NUMBER, nil
		case nodes.BINARYOPERATOR_MUL:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat * rightFloat), vartype.NUMBER, nil
		case nodes.BINARYOPERATOR_DIV:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat / rightFloat), vartype.NUMBER, nil

		case nodes.BINARYOPERATOR_AND:
			leftBool := utilities.StringAsBool(left)
			rightBool := utilities.StringAsBool(right)
			return utilities.BoolAsString(leftBool && rightBool), vartype.BOOLEAN, nil
		case nodes.BINARYOPERATOR_OR:
			leftBool := utilities.StringAsBool(left)
			rightBool := utilities.StringAsBool(right)
			return utilities.BoolAsString(leftBool || rightBool), vartype.BOOLEAN, nil

		case nodes.BINARYOPERATOR_GTE:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool >= rightBool), vartype.BOOLEAN, nil
		case nodes.BINARYOPERATOR_LTE:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool <= rightBool), vartype.BOOLEAN, nil
		case nodes.BINARYOPERATOR_GT:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool > rightBool), vartype.BOOLEAN, nil
		case nodes.BINARYOPERATOR_LT:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool < rightBool), vartype.BOOLEAN, nil

		case nodes.BINARYOPERATOR_NEQ:
			return utilities.BoolAsString(left != right), vartype.BOOLEAN, nil
		case nodes.BINARYOPERATOR_EQ:
			return utilities.BoolAsString(left == right), vartype.BOOLEAN, nil
		}
	}

	pair := spikeUtils.GetErrorIndexPair(i.source, node.ToNode().Start)
	return "", vartype.UNKNOWN, errors.New(fmt.Sprintf("Unsupported value node at line %d:%d", pair.Line, pair.Column))
}
