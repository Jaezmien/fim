package twilight

import (
	"fmt"
	"testing"

	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"github.com/stretchr/testify/assert"
)

func AssertToken(t *testing.T, tok *token.Token, tokenType token.TokenType, value string, typeOf string) {
	assert.Equal(t, tok.Type, tokenType, fmt.Sprintf("Token type should be %s", typeOf))
	assert.Equal(t, tok.Value, value, fmt.Sprintf("Token value should be of %s", typeOf))
}

func TestTokenizer(t *testing.T) {
	t.Run("basic report", func(t *testing.T) {
		source :=
		`Dear Princess Celestia: Hello World!
		Your faithful student, Twilight Sparkle.
		`

		tokens := Parse(source)

		AssertToken(t, tokens.Dequeue().Value, token.TokenType_ReportHeader, "Dear Princess Celestia:", "header")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Identifier, "Hello World", "literal")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Punctuation, "!", "punctuation")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_ReportFooter, "Your faithful student,", "footer")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Identifier, "Twilight Sparkle", "literal")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_Punctuation, ".", "punctuation")
		AssertToken(t, tokens.Dequeue().Value, token.TokenType_EndOfFile, "", "end_of_file")
		assert.Equal(t, tokens.Len(), 0, "Tokens should be empty")
	})
}
