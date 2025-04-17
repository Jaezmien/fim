package celestia

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
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

func (p *Paragraph) Execute(parameters ...*variable.DynamicVariable) (*variable.DynamicVariable, error) {
	p.Interpreter.Variables.PushScope()

	for idx, expecting := range p.FunctionNode.Parameters {
		if idx < len(parameters) {
			received := parameters[idx]

			if received.GetType() != expecting.VariableType {
				return nil, p.FunctionNode.ToNode().CreateError(fmt.Sprintf("Expecting parameter type %s, got %s", expecting.VariableType, received.GetType()), p.Interpreter.source)
			}

			p.Interpreter.Variables.PushVariable(&Variable{
				Name:            expecting.Name,
				DynamicVariable: received,
			}, false)
		} else {
			value, ok := expecting.VariableType.GetDefaultValue()
			if !ok {
				return nil, p.FunctionNode.ToNode().CreateError(fmt.Sprintf("Could not get default value of %s (type %s)", expecting.Name, expecting.VariableType), p.Interpreter.source)
			}

			defaultVariable := variable.FromValueType(value, expecting.VariableType)
			p.Interpreter.Variables.PushVariable(&Variable{
				Name:            expecting.Name,
				DynamicVariable: defaultVariable,
			}, false)
		}
	}
	if len(p.FunctionNode.Parameters) > 0 {
		for i := range min(len(p.FunctionNode.Parameters), len(parameters)) {
			expecting := p.FunctionNode.Parameters[i]
			received := parameters[i]

			if received.GetType() != expecting.VariableType {
				return nil, p.FunctionNode.ToNode().CreateError(fmt.Sprintf("Expecting parameter type %s, got %s", expecting.VariableType, received.GetType()), p.Interpreter.source)
			}

			p.Interpreter.Variables.PushVariable(&Variable{
				Name:            expecting.Name,
				DynamicVariable: received,
			}, false)
		}
	}

	value, err := p.Interpreter.EvaluateStatementsNode(p.FunctionNode.Body)
	p.Interpreter.Variables.PopScope()

	if value != nil && p.FunctionNode.ReturnType == variable.UNKNOWN {
		return nil, p.FunctionNode.Node.CreateError(fmt.Sprintf("Paragraph '%s' with no return type returned a value", p.Name), p.Interpreter.source)
	}
	if value != nil && value.GetType() != p.FunctionNode.ReturnType {
		return nil, p.FunctionNode.Node.CreateError(fmt.Sprintf("Paragraph '%s' expected return value of type '%s', received '%s'", p.Name, p.FunctionNode.ReturnType, value.GetType()), p.Interpreter.source)
	}

	return value, err
}
