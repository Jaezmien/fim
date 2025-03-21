package celestia

import "git.jaezmien.com/Jaezmien/fim/spike/nodes"

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
	}

	return p
}

func (p *Paragraph) Execute() (error) {
	err := p.Interpreter.EvaluateStatementsNode(p.FunctionNode.Body)

	return err
}
