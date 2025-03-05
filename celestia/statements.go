package celestia

import (
	"strconv"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
)

func (i *Interpreter) EvaluateStatementsNode(statements *nodes.StatementsNode) {
	for _, statement := range statements.Statements {
		if statement.Type() == nodes.TYPE_PRINT {
			printNode := statement.(*nodes.PrintNode)

			value := i.EvaluateValueNode(printNode.Value)
			i.Writer.Write([]byte(value))

			if printNode.NewLine {
				i.Writer.Write([]byte("\n"))
			}
		}
	}
}

func (i *Interpreter) EvaluateValueNode(node nodes.INode) string {
	if node.Type() == nodes.TYPE_LITERAL {
		literalNode := node.(*nodes.LiteralNode)

		if literalNode.ValueType == vartype.BOOLEAN {
			if literalNode.GetValueBoolean() {
				return "true"
			} else {
				return "false"
			}
		}
		if literalNode.ValueType == vartype.STRING {
			return literalNode.GetValueString()
		}
		if literalNode.ValueType == vartype.CHARACTER {
			return literalNode.GetValueCharacter()
		}
		if literalNode.ValueType == vartype.NUMBER {
			return strconv.FormatFloat(literalNode.GetValueNumber(), 'f', -1, 64)
		}
	}

	panic("Unsupported value node")
}
