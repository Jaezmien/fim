package celestia

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/vartype"
)

type Paragraph struct {
	Interpreter  *Interpreter
	FunctionNode *nodes.FunctionNode

	Name string
	Main bool
}

func NewParagraph(interpreter *Interpreter, node *nodes.FunctionNode) *Paragraph {
	p := &Paragraph{
		Interpreter:  interpreter,
		FunctionNode: node,
		Main:         node.Main,
		Name:         node.Name,
	}

	return p
}

func (p *Paragraph) Execute() (*vartype.DynamicVariable, error) {
	p.Interpreter.Variables.PushScope()
	value, err := p.Interpreter.EvaluateStatementsNode(p.FunctionNode.Body)
	p.Interpreter.Variables.PopScope()

	if value != nil && p.FunctionNode.ReturnType == vartype.UNKNOWN {
		return nil, p.Interpreter.CreateErrorFromNode(p.FunctionNode.Node, fmt.Sprintf("Paragraph '%s' with no return type returned a value", p.Name))
	}
	if value != nil && value.GetType() != p.FunctionNode.ReturnType {
		return nil, p.Interpreter.CreateErrorFromNode(p.FunctionNode.Node, fmt.Sprintf("Paragraph '%s' expected value with type '%s', received '%s'", p.Name, p.FunctionNode.ReturnType, value.GetType()))
	}

	return value, err
}
