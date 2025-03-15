package celestia

import (
	"git.jaezmien.com/Jaezmien/fim/celestia/utilities"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) {
	for _, statement := range statements.Statements {
		if statement.Type() == nodes.TYPE_PRINT {
			printNode := statement.(*nodes.PrintNode)

			value, _ := i.EvaluateValueNode(printNode.Value)
			i.Writer.Write([]byte(value))

			if printNode.NewLine {
				i.Writer.Write([]byte("\n"))
			}
		}
	}
}

func (i *Interpreter) EvaluateValueNode(node nodes.INode) (string, vartype.VariableType) {
	if node.Type() == nodes.TYPE_LITERAL {
		literalNode := node.(*nodes.LiteralNode)

		if literalNode.ValueType == vartype.BOOLEAN {
			return utilities.BoolAsString(literalNode.GetValueBoolean()), literalNode.ValueType
		}
		if literalNode.ValueType == vartype.STRING {
			return literalNode.GetValueString(), literalNode.ValueType
		}
		if literalNode.ValueType == vartype.CHARACTER {
			return literalNode.GetValueCharacter(), literalNode.ValueType
		}
		if literalNode.ValueType == vartype.NUMBER {
			return utilities.FloatAsString(literalNode.GetValueNumber()), literalNode.ValueType
		}
	}

	if node.Type() == nodes.TYPE_BINARYEXPRESSION {
		binaryNode := node.(*nodes.BinaryExpressionNode)

		left, leftType := i.EvaluateValueNode(binaryNode.Left)
		right, rightType := i.EvaluateValueNode(binaryNode.Right)

		if binaryNode.Operator == nodes.BINARYOPERATOR_ADD {
			if leftType == vartype.STRING || rightType == vartype.STRING {
				return (left + right), vartype.STRING
			}
		}

		switch binaryNode.Operator {
		case nodes.BINARYOPERATOR_ADD:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat + rightFloat), vartype.NUMBER
		case nodes.BINARYOPERATOR_SUB:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat - rightFloat), vartype.NUMBER
		case nodes.BINARYOPERATOR_MUL:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat * rightFloat), vartype.NUMBER
		case nodes.BINARYOPERATOR_DIV:
			leftFloat := utilities.StringAsFloat(left)
			rightFloat := utilities.StringAsFloat(right)
			return utilities.FloatAsString(leftFloat / rightFloat), vartype.NUMBER

		case nodes.BINARYOPERATOR_AND:
			leftBool := utilities.StringAsBool(left)
			rightBool := utilities.StringAsBool(right)
			return utilities.BoolAsString(leftBool && rightBool), vartype.BOOLEAN
		case nodes.BINARYOPERATOR_OR:
			leftBool := utilities.StringAsBool(left)
			rightBool := utilities.StringAsBool(right)
			return utilities.BoolAsString(leftBool || rightBool), vartype.BOOLEAN

		case nodes.BINARYOPERATOR_GTE:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool >= rightBool), vartype.BOOLEAN
		case nodes.BINARYOPERATOR_LTE:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool <= rightBool), vartype.BOOLEAN
		case nodes.BINARYOPERATOR_GT:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool > rightBool), vartype.BOOLEAN
		case nodes.BINARYOPERATOR_LT:
			leftBool := utilities.StringAsFloat(left)
			rightBool := utilities.StringAsFloat(right)
			return utilities.BoolAsString(leftBool < rightBool), vartype.BOOLEAN

		case nodes.BINARYOPERATOR_NEQ:
			return utilities.BoolAsString(left != right), vartype.BOOLEAN
		case nodes.BINARYOPERATOR_EQ:
			return utilities.BoolAsString(left == right), vartype.BOOLEAN
		}
	}

	panic("Unsupported value node")
}
