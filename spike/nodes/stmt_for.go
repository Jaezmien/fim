package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
)

type ForEveryStatementNode struct {
	StatementsNode

	VariableName string
	VariableType variable.VariableType
}
