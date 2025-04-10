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

func (p *Paragraph) Execute(parameters ...*vartype.DynamicVariable) (*vartype.DynamicVariable, error) {
	p.Interpreter.Variables.PushScope()

	for idx, expecting := range p.FunctionNode.Parameters {
		if idx < len(parameters) {
			received := parameters[idx]

			if received.GetType() != expecting.VariableType {
				return nil, p.Interpreter.CreateErrorFromNode(p.FunctionNode.ToNode(), fmt.Sprintf("Expecting parameter type %s, got %s", expecting.VariableType, received.GetType()))
			}

			p.Interpreter.Variables.PushVariable(&Variable{
				Name: expecting.Name,
				DynamicVariable: received,
			}, false)
		} else {
			value, ok := expecting.VariableType.GetDefaultValue()
			if !ok {
				return nil, p.Interpreter.CreateErrorFromNode(p.FunctionNode.ToNode(), fmt.Sprintf("Could not get default value of %s (type %s)", expecting.Name, expecting.VariableType))
			}

			defaultVariable := vartype.FromValueType(value, expecting.VariableType)
			p.Interpreter.Variables.PushVariable(&Variable{
				Name: expecting.Name,
				DynamicVariable: defaultVariable,
			}, false)
		}
	}
	if len(p.FunctionNode.Parameters) > 0 {
		for i := range min(len(p.FunctionNode.Parameters), len(parameters)) {
			expecting := p.FunctionNode.Parameters[i]	
			received := parameters[i]

			if received.GetType() != expecting.VariableType {
				return nil, p.Interpreter.CreateErrorFromNode(p.FunctionNode.ToNode(), fmt.Sprintf("Expecting parameter type %s, got %s", expecting.VariableType, received.GetType()))
			}

			p.Interpreter.Variables.PushVariable(&Variable{
				Name: expecting.Name,
				DynamicVariable: received,
			}, false)
		}
	}

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
