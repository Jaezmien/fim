package twilight

import (
	"fmt"
	"testing"

	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"github.com/stretchr/testify/assert"
)

func AssertToken(t *testing.T, tok *token.Token, tokenType token.TokenType, value string, typeOf string) {
	assert.Equal(t, tokenType, tok.Type, fmt.Sprintf("Token type should be %s", typeOf))
	assert.Equal(t, value, tok.Value, fmt.Sprintf("Token value should be of %s", typeOf))
}

func TestTokenizer(t *testing.T) {
	t.Run("basic report", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
		Your faithful student, Twilight Sparkle.
		`

		tokens := Parse(source)

		AssertToken(t, tokens.Dequeue().Value, token.TokenType_ReportHeader, "Dear Princess Celestia:", "header")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Identifier, "Hello World", "identifier")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Punctuation, "!", "punctuation")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_ReportFooter, "Your faithful student,", "footer")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Identifier, "Twilight Sparkle", "identifier")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Punctuation, ".", "punctuation")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_EndOfFile, "", "end_of_file")
		assert.Equal(t, 0, tokens.Len(), "Tokens should be empty")
	})
}

func TestPostscript(t *testing.T) {
	t.Run("clean postscript", func(t *testing.T) {
		source :=
			`P.S. Hello!`

		tokens := Parse(source)

		AssertToken(t, tokens.Dequeue().Value, token.TokenType_EndOfFile, "", "end_of_file")
		assert.Equal(t, 0, tokens.Len(), "Tokens should be empty")
	})

	t.Run("clean post-postscript", func(t *testing.T) {
		source :=
			`P.S.S. Hello!`

		tokens := Parse(source)

		AssertToken(t, tokens.Dequeue().Value, token.TokenType_EndOfFile, "", "end_of_file")
		assert.Equal(t, 0, tokens.Len(), "Tokens should be empty")
	})
	t.Run("ignore invalid postscript", func(t *testing.T) {
		source :=
			`P.S.Hello!`

		tokens := Parse(source)

		assert.NotEqual(t, 0, tokens.Len(), "Tokens should not be empty")
	})

	t.Run("clean inside postscript", func(t *testing.T) {
		source :=
			`Hello
			P.S. Hello!
			World`

		tokens := Parse(source)

		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Identifier, "Hello", "")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Identifier, "World", "")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_EndOfFile, "", "end_of_file")
		assert.Equal(t, 0, tokens.Len(), "Tokens should be empty")
	})
}
