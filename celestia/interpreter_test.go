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
			`Dear Princess Celestia: Outputs!
			Today I learned how to output something!
			I said 1!
			That's all about how to output something.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})

	t.Run("should print without newline", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Newline Outputs!
			Today I learned how to output in the same line!
			I quickly said 1!
			That's all about how to output in the same line.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1")
	})
}

func TestBasicReports(t *testing.T) {
	t.Run("should print string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing strings!
			Today I learned how to print a string!
			I said "Hello World"!
			That's all about how to print a string.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "Hello World\n")
	})

	t.Run("should print character", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing a character!
			Today I learned how to print a single char!
			I said 'a'!
			That's all about how to print a single char.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "a\n")
	})

	t.Run("should print boolean", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing a boolean!
			Today I learned how to print a boolean!
			I said correct!
			That's all about how to print a boolean.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "true\n")
	})

	t.Run("should print number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing numbers!
			Today I learned how to print a numeric value!
			I said 1!
			That's all about how to print a numeric value.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
}

func TestExpressionNode(t *testing.T) {
	t.Run("should add number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Expressions!
			Today I learned how to evaluate expressions!
			I said 1 plus 1!
			That's all about how to evaluate expressions.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "2\n")
	})
	t.Run("should subtract number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Expressions!
			Today I learned how to evaluate expressions!
			I said 1 minus 1!
			That's all about how to evaluate expressions.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "0\n")
	})
	t.Run("should concatenate string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Concatenations!
			Today I learned how to concatenate strings!
			I said "Hello" plus " " plus "World"!
			That's all about how to concatenate strings.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "Hello World\n")
	})
	t.Run("should concatenate multiple types", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Concatenations!
			Today I learned how to concatenate different types!
			I said "Hello " plus 1!
			That's all about how to concatenate different types.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "Hello 1\n")
	})
}

func TestDeclaration(t *testing.T) {
	t.Run("should create global variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Globals!
			Did you know that Spike is the number 1?
			Today I learned how to print a value!
			I said Spike!
			That's all about how to print a value.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
	t.Run("should create local variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Globals!
			Today I learned how to print a value!
			Did you know that Spike is the number 1?
			I said Spike!
			That's all about how to print a value.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
	t.Run("should create an empty variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Empties!
			Today I learned how to output an empty variable!
			Did you know that Spike is a word?
			I said Spike!
			That's all about how to output an empty variable.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "\n")
	})
	t.Run("should create variable from another variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Variables!
			Today I learned how to create variables!
			Did you know that Spike is the number 1?
			Did you know that Owlowiscious is the number Spike plus 1?
			I said Owlowiscious!
			That's all about how to create variables.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "2\n")
	})
}

func TestModify(t *testing.T) {
	t.Run("should modify local variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Variables!
			Today I learned how to modify variables!
			Did you know that Spike is the number 2?
			Spike became 1.
			I said Spike!
			That's all about how to modify variables.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
	t.Run("should modify global variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Variables!
			Did you know that Spike is the number 2?
			Today I learned how to modify variables!
			Spike became 1.
			I said Spike!
			That's all about how to modify variables.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
	t.Run("should convert number to string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Converting!
			Today I learned how to convert types!
			Did you know that Spike is the word "Hello"?
			Spike became 1.
			I said Spike!
			That's all about how to convert types.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "1\n")
	})
	t.Run("should convert boolean to string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Converting!
			Today I learned how to convert types!
			Did you know that Spike is the word "Hello"?
			Spike became correct.
			I said Spike!
			That's all about how to convert types.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, "true\n")
	})
}
