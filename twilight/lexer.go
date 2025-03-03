package twilight

import (
	"slices"
	"strconv"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/twilight/parsers"
	"git.jaezmien.com/Jaezmien/fim/twilight/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

var punctuations = [...]rune{'.', '!', '?', ':', ','}

func isRunePunctuation(r rune) bool {
	return utilities.ContainsRune(r, punctuations[:])
}

var booleanStrings = [...]string{"yes", "true", "right", "correct", "no", "false", "wrong", "incorrect"}

func processTokenType(t *token.Token, condition bool, resultType token.TokenType) {
	if t.Type != token.TokenType_Literal {
		return
	}

	if !condition {
		return
	}

	t.Type = resultType
}

func createTokens(partialTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	for partialTokens.Len() > 0 {
		t := partialTokens.Dequeue().Value
		t.Type = token.TokenType_Literal

		processTokenType(t, t.Length == 1 && isRunePunctuation(rune(t.Value[0])), token.TokenType_Punctuation)
		processTokenType(t, t.Length == 1 && t.Value == "\n", token.TokenType_NewLine)
		processTokenType(t, t.Length == 1 && t.Value == " ", token.TokenType_Whitespace)
		processTokenType(t, t.Length == 1 && strings.HasPrefix(t.Value, "(") && strings.HasSuffix(t.Value, ")"), token.TokenType_CommentParen)

		processTokenType(t, t.Length > 1 && strings.HasPrefix(t.Value, "\"") && strings.HasSuffix(t.Value, "\""), token.TokenType_String)
		processTokenType(t, t.Length > 1 && strings.HasPrefix(t.Value, "'") && strings.HasSuffix(t.Value, "'"), token.TokenType_Character)
		if _, err := strconv.ParseFloat(t.Value, 64); err == nil {
			processTokenType(t, t.Length > 1, token.TokenType_Number)
		}
		processTokenType(t, t.Length > 1 && slices.Contains(booleanStrings[:], t.Value), token.TokenType_Character)
		processTokenType(t, t.Value == "nothing", token.TokenType_Null)

		if t.Length == 0 && t.Type == token.TokenType_Literal {
			continue
		}

		tokens.Queue(t)
	}

	tokens.Queue(&token.Token{
		Start:  tokens.Last().Value.Start,
		Length: 0,
		Value:  "",
		Type:   token.TokenType_EndOfFile,
	})

	return tokens
}

type processMultiTokenResult = func(tokens *queue.Queue[*token.Token]) int

func processMultiTokenType(tokens *queue.Queue[*token.Token], condition processMultiTokenResult, resultType token.TokenType) {
	if tokens.Len() <= 0 {
		return
	}

	if tokens.First().Value.Type != token.TokenType_Literal {
		return
	}

	amount := condition(tokens)
	if amount <= 0 {
		return
	}

	token := utilities.MergeTokens(tokens, amount)
	token.Type = resultType

	tokens.QueueFront(token)
}

func mergeMultitokens(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	for oldTokens.Len() > 0 {
		processMultiTokenType(oldTokens, parsers.IsReportHeader, token.TokenType_ReportHeader)
		processMultiTokenType(oldTokens, parsers.IsReportFooter, token.TokenType_ReportFooter)

		processMultiTokenType(oldTokens, parsers.IsFunctionHeaderMain, token.TokenType_FunctionMain)
		processMultiTokenType(oldTokens, parsers.IsFunctionHeader, token.TokenType_FunctionHeader)
		processMultiTokenType(oldTokens, parsers.IsFunctionFooter, token.TokenType_FunctionFooter)
		processMultiTokenType(oldTokens, parsers.IsFunctionParameter, token.TokenType_FunctionParameter)
		processMultiTokenType(oldTokens, parsers.IsFunctionReturn, token.TokenType_FunctionReturn)

		processMultiTokenType(oldTokens, parsers.IsPrintMethod, token.TokenType_Print)
		processMultiTokenType(oldTokens, parsers.IsPrintNewlineMethod, token.TokenType_PrintNewline)
		processMultiTokenType(oldTokens, parsers.IsReadMethod, token.TokenType_Prompt)
		processMultiTokenType(oldTokens, parsers.IsFunctionCallMethod, token.TokenType_FunctionCall)

		tokens.Queue(oldTokens.Dequeue().Value)
	}

	return tokens
}

func mergeLiterals(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	var literalToken *token.Token
	for oldTokens.Len() > 0 {
		t := oldTokens.Dequeue().Value

		if t.Type != token.TokenType_Literal {
			if literalToken != nil &&
				t.Type == token.TokenType_Whitespace &&
				oldTokens.Len() >= 1 &&
				oldTokens.First().Value.Type == token.TokenType_Literal {
				literalToken.Append(t)
				continue
			}

			if literalToken != nil {
				tokens.Queue(literalToken)
				literalToken = nil
			}
			tokens.Queue(t)
			continue
		}

		if literalToken != nil {
			literalToken.Append(t)
			continue
		}
		literalToken = t
	}

	// Flush remaining
	if literalToken != nil {
		tokens.Queue(literalToken)
		literalToken = nil
	}

	return tokens
}

func cleanTokens(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	for oldTokens.Len() > 0 {
		t := oldTokens.Dequeue().Value

		if t.Type == token.TokenType_NewLine {
			continue
		}
		if t.Type == token.TokenType_Whitespace {
			continue
		}
		if t.Type == token.TokenType_CommentParen {
			continue
		}
		if t.Type == token.TokenType_CommentPostScript {
			for oldTokens.Len() > 0 && oldTokens.First().Value.Type != token.TokenType_NewLine {
				oldTokens.Dequeue()
			}
			continue
		}

		tokens.Queue(t)
	}

	return tokens
}
