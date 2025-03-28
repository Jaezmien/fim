package twilight

import (
	"testing"

	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"github.com/stretchr/testify/assert"
)

func CheckTokens(t *testing.T, tokens []*token.Token, checks []struct{ tokenType token.TokenType; expectedValue string }) {
	if !assert.Equal(t, len(checks), len(tokens), "Mismatch token count") {
		return
	}
	
	for idx, check := range checks {
		assert.Equal(t, check.tokenType, tokens[idx].Type)
		assert.Equal(t, check.expectedValue, tokens[idx].Value)
	}
}

func TestTokenizer(t *testing.T) {
	t.Run("basic report", func(t *testing.T) {
		source :=
			`Dear Princess Celestia: Hello World!
		Your faithful student, Twilight Sparkle.
		`

		tokens := Parse(source)

		checks := []struct{
			tokenType token.TokenType
			expectedValue string
		}{
			{ tokenType: token.TokenType_ReportHeader, expectedValue: "Dear Princess Celestia:" },
			{ tokenType: token.TokenType_Identifier, expectedValue: "Hello World" },
			{ tokenType: token.TokenType_Punctuation, expectedValue: "!" },
			{ tokenType: token.TokenType_ReportFooter, expectedValue: "Your faithful student," },
			{ tokenType: token.TokenType_Identifier, expectedValue: "Twilight Sparkle" },
			{ tokenType: token.TokenType_Punctuation, expectedValue: "." },
			{ tokenType: token.TokenType_EndOfFile, expectedValue: "" },
		}

		CheckTokens(t, tokens, checks)
	})
}

func TestPostscript(t *testing.T) {
	t.Run("clean postscript", func(t *testing.T) {
		source :=
			`P.S. Hello!`

		tokens := Parse(source)

		checks := []struct{
			tokenType token.TokenType
			expectedValue string
		}{
			{ tokenType: token.TokenType_EndOfFile, expectedValue: "" },
		}

		CheckTokens(t, tokens, checks)
	})

	t.Run("clean post-postscript", func(t *testing.T) {
		source :=
			`P.S.S. Hello!`

		tokens := Parse(source)

		checks := []struct{
			tokenType token.TokenType
			expectedValue string
		}{
			{ tokenType: token.TokenType_EndOfFile, expectedValue: "" },
		}

		CheckTokens(t, tokens, checks)
	})
	t.Run("ignore invalid postscript", func(t *testing.T) {
		source :=
			`P.S.Hello!`

		tokens := Parse(source)

		checks := []struct{
			tokenType token.TokenType
			expectedValue string
		}{
			{ tokenType: token.TokenType_Identifier, expectedValue: "P" },
			{ tokenType: token.TokenType_Punctuation, expectedValue: "." },
			{ tokenType: token.TokenType_Identifier, expectedValue: "S" },
			{ tokenType: token.TokenType_Punctuation, expectedValue: "." },
			{ tokenType: token.TokenType_Identifier, expectedValue: "Hello" },
			{ tokenType: token.TokenType_Punctuation, expectedValue: "!" },
			{ tokenType: token.TokenType_EndOfFile, expectedValue: "" },
		}

		CheckTokens(t, tokens, checks)
	})

	t.Run("clean inside postscript", func(t *testing.T) {
		source :=
			`Hello
			P.S. Hello!
			World`

		tokens := Parse(source)

		checks := []struct{
			tokenType token.TokenType
			expectedValue string
		}{
			{ tokenType: token.TokenType_Identifier, expectedValue: "Hello" },
			{ tokenType: token.TokenType_Identifier, expectedValue: "World" },
			{ tokenType: token.TokenType_EndOfFile, expectedValue: "" },
		}

		CheckTokens(t, tokens, checks)
	})
}
