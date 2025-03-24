package celestia

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
)

type Interpreter struct {
	Writer io.Writer
	ErrorWriter io.Writer

	Prompt func(prompt string) (string, error)

	reportNode *nodes.ReportNode
	source     string

	Variables *VariableManager
	Paragraphs []*Paragraph
}

func NewInterpreter(reportNode *nodes.ReportNode, source string) (*Interpreter, error) {
	interpreter := &Interpreter{
		Writer:     os.Stdout,
		ErrorWriter:     os.Stderr,
		reportNode: reportNode,
		source:     source,
		Paragraphs: make([]*Paragraph, 0),
		Variables: NewVariableManager(),
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

	for _, node := range interpreter.reportNode.Body {
		if node.Type() == nodes.TYPE_FUNCTION {
			funcNode := node.(*nodes.FunctionNode)
			paragraph := NewParagraph(interpreter, funcNode)

			for _, p := range interpreter.Paragraphs {
				if p.Name == paragraph.Name {
					return nil, interpreter.CreateErrorFromNode(funcNode.ToNode(), fmt.Sprintf("Paragraph '%s' already exists", p.Name))
				}
			}

			interpreter.Paragraphs = append(interpreter.Paragraphs, paragraph)

			continue
		}
		if node.Type() == nodes.TYPE_VARIABLE_DECLARATION {
			variableNode := node.(*nodes.VariableDeclarationNode)

			value, err := interpreter.EvaluateValueNode(variableNode.Value, false)
			if err != nil {
				return nil, err
			}

			variable := &Variable{
				Name: variableNode.Identifier,
				DynamicVariable: value,
				Constant: variableNode.Constant,
			}
			
			interpreter.Variables.PushVariable(variable, true)

			continue
		}

		return nil, interpreter.CreateErrorFromNode(node.ToNode(), fmt.Sprintf("Unsupported node type: %s", node.Type()))
	}

	return interpreter, nil
}

func (i *Interpreter) ReportName() string {
	return i.reportNode.Name
}
func (i *Interpreter) ReportAuthor() string {
	return i.reportNode.Author
}

func (i *Interpreter) CreateErrorFromIndex(index int, errorMessage string) error {
	pair := luna.GetErrorIndexPair(i.source, index)
	return errors.New(fmt.Sprintf("[line: %d:%d] %s", pair.Line, pair.Column, errorMessage))
}
func (i *Interpreter) CreateErrorFromNode(n nodes.Node, errorMessage string) error {
	return i.CreateErrorFromIndex(n.Start, errorMessage)
}
