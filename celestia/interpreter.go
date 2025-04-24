package celestia

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
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
		if funcNode, ok := n.(*nodes.FunctionNode); ok {
			paragraph := NewParagraph(interpreter, funcNode)

			for _, p := range interpreter.Paragraphs {
				if p.Name == paragraph.Name {
					return nil, funcNode.ToNode().CreateError(fmt.Sprintf("Paragraph '%s' already exists", p.Name), interpreter.source)
				}
			}

			interpreter.Paragraphs = append(interpreter.Paragraphs, paragraph)

			continue
		}

		if variableNode, ok := n.(*nodes.VariableDeclarationNode); ok {
			if interpreter.Variables.Get(variableNode.Identifier, true) != nil {
				return nil, variableNode.ToNode().CreateError(fmt.Sprintf("Variable '%s' already exists.", variableNode.Identifier), interpreter.source)
			}

			value, err := interpreter.EvaluateValueNode(variableNode.Value, false)
			if err != nil {
				return nil, err
			}

			if value.GetType() == variable.UNKNOWN {
				if variableNode.ValueType.IsArray() {
					value = variable.NewDictionaryVariable(variableNode.ValueType)
				} else {
					defaultValue, ok := variableNode.ValueType.GetDefaultValue()
					if !ok {
						panic("Intepreter@EvaluateStatementsNode could not get default value.")
					}
					value = variable.FromValueType(defaultValue, variableNode.ValueType)
				}
			}

			if !variableNode.ValueType.IsArray() {
				value = value.Clone()
			}

			variable := &Variable{
				Name:            variableNode.Identifier,
				DynamicVariable: value,
				Constant:        variableNode.Constant,
			}

			interpreter.Variables.PushVariable(variable, true)

			continue
		}

		return nil, n.ToNode().CreateError("Unsupported report body node", interpreter.source)
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
