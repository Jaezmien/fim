package celestia

import (
	"errors"
	"fmt"
	"io"
	"os"

	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/spike/utilities"
)

type Interpreter struct {
	Writer io.Writer
	ErrorWriter io.Writer
	Reader io.Reader

	reportNode *nodes.ReportNode
	source     string

	Paragraphs []*Paragraph
}

func NewInterpreter(reportNode *nodes.ReportNode, source string) (*Interpreter, error) {
	interpreter := &Interpreter{
		Writer:     os.Stdout,
		ErrorWriter:     os.Stderr,
		Reader:     os.Stdin,
		reportNode: reportNode,
		source:     source,
		Paragraphs: make([]*Paragraph, 0),
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

		panic(fmt.Sprintf("Unsupported node type: %d", node.Type()))
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
	pair := utilities.GetErrorIndexPair(i.source, index)
	return errors.New(fmt.Sprintf("[line: %d:%d] %s", pair.Line, pair.Column, errorMessage))
}
func (i *Interpreter) CreateErrorFromNode(n nodes.Node, errorMessage string) error {
	return i.CreateErrorFromIndex(n.Start, errorMessage)
}
