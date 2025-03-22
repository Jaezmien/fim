package celestia

import (
	"bytes"
	"io"
	"testing"

	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/twilight"
	"github.com/stretchr/testify/assert"
)

func ExecuteBasicReport(t *testing.T, source string, expected string) {
	tokens := twilight.Parse(source)
	report, err := spike.CreateReport(tokens.Flatten(), source)
	if !assert.NoError(t, err) {
		return
	}
	interpreter, err := NewInterpreter(report, source)
	if !assert.NoError(t, err) {
		return
	}

	buffer := &bytes.Buffer{}
	interpreter.Writer = buffer

	var mainParagraph *Paragraph
	for _, paragraph := range interpreter.Paragraphs {
		if paragraph.Main {
			mainParagraph = paragraph
			break
		}
	}
	if !assert.NotNil(t, mainParagraph) {
		return
	}

	err = mainParagraph.Execute()
	if !assert.NoError(t, err) {
		return
	}

	data, err := io.ReadAll(buffer)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Equal(t, expected, string(data)) {
		return
	}
}

func TestReport(t *testing.T) {
	source := `Dear Princess Celestia: Hello World!
		Today I learned how to say hello world!
		I said "Hello World"!
		That's all about how to say hello world.
		Your faithful student, Twilight Sparkle.
		`

	tokens := twilight.Parse(source)
	report, err := spike.CreateReport(tokens.Flatten(), source)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "Hello World", report.Name, "Mismatch report name")
	assert.Equal(t, "Twilight Sparkle", report.Author, "Mismatch report author")

	interpreter, err := NewInterpreter(report, source)
	if !assert.NoError(t, err) {
		return
	}

	buffer := &bytes.Buffer{}
	interpreter.Writer = buffer

	var mainParagraph *Paragraph
	for _, paragraph := range interpreter.Paragraphs {
		if paragraph.Main {
			mainParagraph = paragraph
			break
		}
	}
	if !assert.NotNil(t, mainParagraph) {
		return
	}
	assert.Equal(t, "how to say hello world", mainParagraph.FunctionNode.Name, "Mismatch function name")

	mainParagraph.Execute()
	data, err := io.ReadAll(buffer)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "Hello World\n", string(data))
}

func TestIO(t *testing.T) {
	t.Run("should print", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said 1!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})

	t.Run("should print with newline", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I quickly said 1!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1")
	})
}

func TestBasicReports(t *testing.T) {
	t.Run("should print string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said "Hello World"!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "Hello World\n")
	})

	t.Run("should print character", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said 'a'!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "a\n")
	})

	t.Run("should print boolean", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said correct!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "true\n")
	})

	t.Run("should print number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said 1!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
}

func TestExpressionNode(t *testing.T) {
	t.Run("should add number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said 1 plus 1!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "2\n")
	})
	t.Run("should subtract number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said 1 minus 1!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "0\n")
	})
	t.Run("should concatenate string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said "Hello" plus " " plus "World"!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "Hello World\n")
	})
	t.Run("should concatenate number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			I said "Hello " plus 1!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "Hello 1\n")
	})
}

func TestDeclaration(t *testing.T) {
	t.Run("should create global variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Did you know that Spike is the number 1?
			Today I learned how to say hello world!
			I said Spike!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
	t.Run("should create local variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			Did you know that Spike is the number 1?
			I said Spike!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
	t.Run("should create variable from another variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
			Today I learned how to say hello world!
			Did you know that Spike is the number 1?
			Did you know that Owlowiscious is the number Spike plus 1?
			I said Owlowiscious!
			That's all about how to say hello world.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "2\n")
	})
}
