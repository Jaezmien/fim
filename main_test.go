package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"git.jaezmien.com/Jaezmien/fim/celestia"
	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/twilight"
	"github.com/stretchr/testify/assert"
)

type BasicReportOptions struct {
	Expects       string
	IgnoreExpects bool
	CompileOnly   bool
	Error         bool
	Prompt        func(prompt string) (string, error)
}

func CreateReport(t *testing.T, source string, options BasicReportOptions) (*celestia.Interpreter, bool) {
	tokens := twilight.Parse(source)
	report, err := spike.CreateReport(tokens, source)

	if err != nil {
		if options.Error {
			return nil, false
		}

		return nil, assert.NoError(t, err, "handled by spike")
	}

	interpreter, err := celestia.NewInterpreter(report, source)

	if err != nil {
		if options.Error {
			return nil, false
		}

		return nil, assert.NoError(t, err, "handled pre-celestia")
	}

	return interpreter, true
}
func GetMainParagraph(t *testing.T, interpreter *celestia.Interpreter) (*celestia.Paragraph, bool) {
	var mainParagraph *celestia.Paragraph
	for _, paragraph := range interpreter.Paragraphs {
		if paragraph.Main {
			mainParagraph = paragraph
			break
		}
	}
	if !assert.NotNil(t, mainParagraph) {
		return nil, false
	}

	return mainParagraph, true
}

func ExecuteBasicReport(t *testing.T, source string, options BasicReportOptions) {
	interpreter, ok := CreateReport(t, source, options)
	if !ok {
		return
	}

	if options.CompileOnly {
		return
	}

	buffer := &bytes.Buffer{}
	interpreter.Writer = buffer
	if options.Prompt != nil {
		interpreter.Prompt = options.Prompt
	}

	mainParagraph, ok := GetMainParagraph(t, interpreter)
	if !ok {
		return
	}

	_, err := mainParagraph.Execute()
	if options.Error && !assert.Error(t, err, "handled by celestia") {
		return
	}
	if !options.Error && !assert.NoError(t, err, "handled by celestia") {
		return
	}

	if options.IgnoreExpects {
		return
	}

	data, err := io.ReadAll(buffer)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Equal(t, options.Expects, string(data)) {
		return
	}
}

func TestReports(t *testing.T) {
	reports := []struct {
		BasicReportOptions
		Name string
	}{
		{
			Name: "array.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "Banana Cake\nGala\n",
			},
		},
		{
			Name: "brainfuck.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "Hello World!",
			},
		},
		{
			Name: "bubblesort.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "1\n2\n3\n4\n5\n7\n7\n",
			},
		},
		{
			Name: "cider.fim",
			BasicReportOptions: BasicReportOptions{
				IgnoreExpects: true,
			},
		},
		{
			Name: "deadfish.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "Hello world",
			},
		},
		{
			Name: "digitalroot.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "9\n",
			},
		},
		{
			Name: "disan.fim",
			BasicReportOptions: BasicReportOptions{
				IgnoreExpects: true,
				Prompt: func(prompt string) (string, error) {
					return "10", nil
				},
			},
		},
		{
			Name: "e.fim",
			BasicReportOptions: BasicReportOptions{
				CompileOnly: true,
			},
		},
		{
			Name: "factorial.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "6\n24\n120\n720\n",
			},
		},
		{
			Name: "fibonacci.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "34\n",
			},
		},
		{
			Name: "fizzbuzz.fim",
			BasicReportOptions: BasicReportOptions{
				IgnoreExpects: true,
			},
		},
		{
			Name: "hello.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "Hello World!\n",
			},
		},
		{
			Name: "input.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "Applejack said: Awesome!\n",
				Prompt: func(prompt string) (string, error) {
					return "Awesome!", nil
				},
			},
		},
		{
			Name: "insertionsort.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "1\n2\n3\n4\n5\n7\n7\n",
			},
		},
		{
			Name: "mississippis.fim",
			BasicReportOptions: BasicReportOptions{
				IgnoreExpects: true,
			},
		},
		{
			Name: "parameters.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "x\n1\ny\n0\n",
			},
		},
		{
			Name: "quicksort.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "1\n2\n3\n4\n5\n7\n7\n",
			},
		},
		{
			Name: "recursion.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "5\n4\n3\n2\n1\n",
			},
		},
		{
			Name: "squareroot.fim",
			BasicReportOptions: BasicReportOptions{
				Prompt: func(prompt string) (string, error) {
					return "50", nil
				},
				Expects: "7.071067984011346\n",
			},
		},
		{
			Name: "string_index.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "T\nw\n",
			},
		},
		{
			Name: "sum.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "5051\n",
			},
		},
		{
			Name: "truth_machine.fim",
			BasicReportOptions: BasicReportOptions{
				Expects: "0\n",
				Prompt: func(prompt string) (string, error) {
					return "0", nil
				},
			},
		},
	}

	for _, report := range reports {
		data, err := os.ReadFile("./samples/" + report.Name)
		if !assert.NoError(t, err) {
			continue
		}
		if !assert.NotEmpty(t, data) {
			continue
		}

		source := string(data)

		t.Logf("Testing report '%s'...", report.Name)
		ExecuteBasicReport(t, source, report.BasicReportOptions)
	}
}
