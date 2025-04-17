package celestia

import (
	"bytes"
	"io"
	"testing"

	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/twilight"
	"github.com/stretchr/testify/assert"
)

func CreateReport(t *testing.T, source string) (*Interpreter, bool) {
	tokens := twilight.Parse(source)
	report, err := spike.CreateReport(tokens, source)
	if !assert.NoError(t, err, "handled by spike") {
		return nil, false
	}
	interpreter, err := NewInterpreter(report, source)
	if !assert.NoError(t, err, "handled pre-celestia") {
		return nil, false
	}

	return interpreter, true
}
func GetMainParagraph(t *testing.T, interpreter *Interpreter) (*Paragraph, bool) {
	var mainParagraph *Paragraph
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

type BasicReportOptions struct {
	Expects string
	Error   bool
	Prompt  func(prompt string) (string, error)
}

func ExecuteBasicReport(t *testing.T, source string, options BasicReportOptions) {
	interpreter, ok := CreateReport(t, source)
	if !ok {
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

	data, err := io.ReadAll(buffer)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Equal(t, options.Expects, string(data)) {
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
	report, err := spike.CreateReport(tokens, source)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "Hello World", report.Title, "Mismatch report name")
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
	})

	t.Run("should print without newline", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Newline Outputs!
			Today I learned how to output in the same line!
			I quickly said 1!
			I quickly said 2!
			I said 3!
			That's all about how to output in the same line.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "123\n"})
	})

	t.Run("should prompt", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Prompts!
			Today I learned how to acquire prompts!
			Did you know that Spike is a word?
			I asked Spike: "What do you want me to say? ".
			I said Spike!
			That's all about how to acquire prompts.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{
			Expects: "Hello!\n",
			Prompt: func(prompt string) (string, error) {
				return "Hello!", nil
			},
		})
	})

	t.Run("should convert response type", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Prompts!
			Today I learned how to acquire prompts!
			Did you know that Spike is a number?
			I asked Spike: "Give me a number! ".
			I said Spike!
			That's all about how to acquire prompts.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{
			Expects: "1\n",
			Prompt: func(prompt string) (string, error) {
				return "1", nil
			},
		})
	})

	t.Run("should error on invalid response type", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Prompts!
			Today I learned how to acquire prompts!
			Did you know that Spike is a number?
			I asked Spike: "Give me a number! ".
			I said Spike!
			That's all about how to acquire prompts.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{
			Error: true,
			Prompt: func(prompt string) (string, error) {
				return "Books!", nil
			},
		})
	})
}

func TestBasicReports(t *testing.T) {
	t.Run("should print nothing", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing empty!
			Today I learned how to print empty!
			I said nothing!
			That's all about how to print empty. 
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "\n"})
	})

	t.Run("should print string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing strings!
			Today I learned how to print a string!
			I said "Hello World"!
			That's all about how to print a string.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello World\n"})
	})

	t.Run("should print character", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing a character!
			Today I learned how to print a single char!
			I said 'a'!
			That's all about how to print a single char.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "a\n"})
	})

	t.Run("should print boolean", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing a boolean!
			Today I learned how to print a boolean!
			I said correct!
			That's all about how to print a boolean.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "true\n"})
	})

	t.Run("should print number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Printing numbers!
			Today I learned how to print a numeric value!
			I said 1!
			That's all about how to print a numeric value.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "2\n"})
	})
	t.Run("should subtract number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Expressions!
			Today I learned how to evaluate expressions!
			I said 1 minus 1!
			That's all about how to evaluate expressions.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "0\n"})
	})
	t.Run("should concatenate string", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Concatenations!
			Today I learned how to concatenate strings!
			I said "Hello" plus " " plus "World"!
			That's all about how to concatenate strings.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello World\n"})
	})
	t.Run("should concatenate multiple types", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Concatenations!
			Today I learned how to concatenate different types!
			I said "Hello " plus 1!
			That's all about how to concatenate different types.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello 1\n"})
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
	})
	t.Run("should create an empty variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Empties!
			Today I learned how to output empty variables!
			Did you know that Spike is a word?
			I said Spike!
			Did you know that Owlowiscious is an argument?
			I said Owlowiscious!
			Did you know that Gummy is a number?
			I said Gummy!
			Did you know that Tank is a letter?
			I said Tank!
			That's all about how to output empty variables.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "\nfalse\n0\n\x00\n"})
	})
	t.Run("should create an explicit empty variable", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Empties!
			Today I learned how to output an empty variable!
			Did you know that Spike is the word nothing?
			I said Spike!
			Did you know that Owlowiscious is the argument nothing?
			I said Owlowiscious!
			Did you know that Gummy is the number nothing?
			I said Gummy!
			Did you know that Tank is the letter nothing?
			I said Tank!
			That's all about how to output an empty variable.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "\nfalse\n0\n\x00\n"})
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "2\n"})
	})
	t.Run("should fail on invalid value type", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Variable Type!
			Today I learned how to throw an error!
			Did you know that Spike is the number "Hello"?
			I said Spike!
			That's all about how to throw an error.
			Your faithful student, Twilight Sparkle.
			`

		tokens := twilight.Parse(source)
		report, err := spike.CreateReport(tokens, source)
		if !assert.NoError(t, err) {
			return
		}
		interpreter, err := NewInterpreter(report, source)
		if !assert.NoError(t, err) {
			return
		}

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

		_, err = mainParagraph.Execute()
		if !assert.Error(t, err) {
			return
		}
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
	})
	t.Run("should set to default value", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Variables!
			Today I learned how to reset a variable!
			Did you know that Spike is the number 2?
			Spike became the word nothing.
			I said Spike!
			That's all about how to reset a variable.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "0\n"})
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
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

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "true\n"})
	})
}

func TestUnary(t *testing.T) {
	t.Run("should increment number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Unaries!
			Today I learned how to increment a value!
			Did you know that Spike is the number 1?
			Spike got one more.
			There was one more Spike.
			I said Spike!
			That's all about how to increment a value.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "3\n"})
	})
	t.Run("should decrement number", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Unaries!
			Today I learned how to decrement a value!
			Did you know that Spike is the number 3?
			Spike got one less.
			There was one less Spike.
			I said Spike!
			That's all about how to decrement a value.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "1\n"})
	})
	t.Run("should only work for a value", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Unaries!
			Today I learned how to decrement a value!
			Did you know that Spike is the word "Apples"?
			Spike got one more.
			I said Spike!
			That's all about how to decrement a value.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Error: true})
	})
}

func TestArray(t *testing.T) {
	t.Run("should create an empty array", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Arrays!
			Today I learned how to print arrays!
			Did you know that Apples has many words?
			I said 1 of Apples!
			That's all about how to print arrays.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "\n"})
	})
	t.Run("should print when out of bounds", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Arrays!
			Today I learned how to print arrays!
			Did you know that Apples has many words?
			I said 100 of Apples!
			That's all about how to print arrays.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "\n"})
	})
	t.Run("should print", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Arrays!
			Today I learned how to print arrays!
			Did you know that Apples has the words "Gala", "Red Delicious", "Mcintosh", "Honeycrisp"?
			I said 1 of Apples!
			I said 2 of Apples!
			I said 3 of Apples!
			I said 4 of Apples!
			That's all about how to print arrays.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Gala\nRed Delicious\nMcintosh\nHoneycrisp\n"})
	})
	t.Run("should print nothing on out of range", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Arrays!
			Today I learned how to print arrays!
			Did you know that Apples has the words "Gala"?
			I said 1 of Apples!
			I said 2 of Apples!
			That's all about how to print arrays.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Gala\n\n"})
	})
	t.Run("should modify at index", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Arrays!
			Today I learned how to modify arrays!
			Did you know that Apples has the words "Gala", "Red Delicious", "Mcintosh", "Honeycrisp"?
			I said 1 of Apples!
			1 of Apples is "Gala!".
			I said 1 of Apples!

			Did you know that Applebloom is the number 1?
			Applebloom of Apples is "Gala".
			I said Applebloom of Apples.
			That's all about how to modify arrays.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Gala\nGala!\nGala\n"})
	})
}

func TestFunctions(t *testing.T) {
	t.Run("should run function", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Functions!
			I learned how to say hello world!
			I said "Hello World"!
			That's all about how to say hello world.
			Today I learned how to run a function!
			I remembered how to say hello world.
			That's all about how to run a function.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello World\n"})
	})
	t.Run("should return a value", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Returns!
			I learned how to ask to get a word!
			Then you get "Hello World"!
			That's all about how to ask.
			Today I learned how to run a function!
			I said how to ask.
			That's all about how to run a function.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello World\n"})
	})
	t.Run("should run even with return", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Returns!
			I learned how to ask to get a word!
			Then you get "Hello World"!
			That's all about how to ask.
			Today I learned how to run a function!
			I remembered how to ask.
			That's all about how to run a function.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{})
	})
	t.Run("should accept a value", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Returns!
			I learned how to give a text using the word text!
			I said "Hello " plus text.
			That's all about how to give a text.
			Today I learned how to run a function!
			I remembered how to give a text using the word "World".
			That's all about how to run a function.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello World\n"})
	})
	t.Run("should handle multiple values", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Returns!
			I learned how to give a text using the word x, the word y!
			I said x plus " " plus y.
			That's all about how to give a text.
			Today I learned how to run a function!
			I remembered how to give a text using the word "Hello", "World".
			That's all about how to run a function.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello World\n"})
	})
	t.Run("should handle a default value", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Returns!
			I learned how to give a text using the word x, the number y!
			I said x.
			I said y.
			That's all about how to give a text.
			Today I learned how to run a function!
			I remembered how to give a text using the word "Hello".
			That's all about how to run a function.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello\n0\n"})
	})
}

func TestIfStatements(t *testing.T) {
	t.Run("should run if statement", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Statements!
			Today I learned how to branch statements!
			Did you know that Spike is the number 1?
			If Spike is equal to 1 then,
				I said "Hello World".
			That's what I would do.
			That's all about how to branch statements.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello World\n"})
	})
	t.Run("should ignore if statement", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Statements!
			Today I learned how to branch statements!
			Did you know that Spike is the number 2?
			If Spike is equal to 1,
				I said "Hello World".
			That's what I would do.
			That's all about how to branch statements.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: ""})
	})
	t.Run("should fallback to else statement", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Statements!
			Today I learned how to branch statements!
			Did you know that Spike is the number 2?
			If Spike is equal to 1,
				I said "Nope! Not this one.".
			Otherwise,
				I said "Hello Equestria".
			That's what I would do.
			That's all about how to branch statements.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello Equestria\n"})
	})
	t.Run("should run if else statement", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Statements!
			Today I learned how to branch statements!
			Did you know that Spike is the number 2?
			If Spike is equal to 1,
				I said "Nope! Not this one.".
			Otherwise Spike is equal to 2, 
				I said "Hello Ponyville".
			Otherwise,
				I said "Well that isn't right!".
			That's what I would do.
			That's all about how to branch statements.
			Your faithful student, Twilight Sparkle.
			`

		ExecuteBasicReport(t, source, BasicReportOptions{Expects: "Hello Ponyville\n"})
	})
}
