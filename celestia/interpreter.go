package celestia

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
)

type Interpreter struct {
	Writer      io.Writer
	ErrorWriter io.Writer

	Prompt func(prompt string) (string, error)

	reportNode *nodes.ReportNode
	source     string

	Variables  *VariableManager
	Paragraphs []*Paragraph
}

// Create a new interpreter based on the ReportNode
func NewInterpreter(reportNode *nodes.ReportNode, source string) (*Interpreter, error) {
	interpreter := &Interpreter{
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
		reportNode:  reportNode,
		source:      source,
		Paragraphs:  make([]*Paragraph, 0),
		Variables:   NewVariableManager(),
	}

	interpreter.Prompt = func(prompt string) (string, error) {
		interpreter.Writer.Write([]byte(prompt))

		scanner := bufio.NewScanner(os.Stdin)

		var response string
		for scanner.Scan() {
			response = scanner.Text()
			break
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}

		return response, nil
	}

	for _, n := range interpreter.reportNode.Body {
		if n.Type() == node.TYPE_FUNCTION {
			funcNode := n.(*nodes.FunctionNode)
			paragraph := NewParagraph(interpreter, funcNode)

			for _, p := range interpreter.Paragraphs {
				if p.Name == paragraph.Name {
					return nil, funcNode.ToNode().CreateError(fmt.Sprintf("Paragraph '%s' already exists", p.Name), interpreter.source)
				}
			}

			interpreter.Paragraphs = append(interpreter.Paragraphs, paragraph)

			continue
		}
		if n.Type() == node.TYPE_VARIABLE_DECLARATION {
			variableNode := n.(*nodes.VariableDeclarationNode)

			value, err := interpreter.EvaluateValueNode(variableNode.Value, false)
			if err != nil {
				return nil, err
			}

			variable := &Variable{
				Name:            variableNode.Identifier,
				DynamicVariable: value,
				Constant:        variableNode.Constant,
			}

			interpreter.Variables.PushVariable(variable, true)

			continue
		}

		return nil, n.ToNode().CreateError(fmt.Sprintf("Unsupported node type: %s", n.Type()), interpreter.source)
	}

	return interpreter, nil
}

// Return the report's title
func (i *Interpreter) ReportTitle() string {
	return i.reportNode.Title
}

// Return the report's author
func (i *Interpreter) ReportAuthor() string {
	return i.reportNode.Author
}
