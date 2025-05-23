package twilight

import (
	"slices"
	"strconv"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
	"git.jaezmien.com/Jaezmien/fim/twilight/parsers"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

// Assign basic TokenTypes to partial tokens.
//
// Note: This will also insert an End Of File token to the end of the queue.
func createTokens(partialTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	punctuations := [...]rune{'.', '!', '?', ':', ','}

	tokenTypeProcessors := []struct {
		condition func(t *token.Token) bool
		result    token.TokenType
	}{
		{condition: func(t *token.Token) bool {
			return t.Length == 1 && slices.Contains(punctuations[:], rune(t.Value[0]))
		}, result: token.TokenType_Punctuation},
		{condition: func(t *token.Token) bool { return t.Length == 1 && t.Value == "\n" }, result: token.TokenType_NewLine},
		{condition: func(t *token.Token) bool { return t.Length == 1 && t.Value == " " }, result: token.TokenType_Whitespace},
		{condition: func(t *token.Token) bool {
			return t.Length >= 1 && strings.HasPrefix(t.Value, "(") && strings.HasSuffix(t.Value, ")")
		}, result: token.TokenType_CommentParen},
		{condition: func(t *token.Token) bool {
			return t.Length >= 1 && strings.HasPrefix(t.Value, "\"") && strings.HasSuffix(t.Value, "\"")
		}, result: token.TokenType_String},
		{condition: func(t *token.Token) bool {
			return t.Length >= 1 && strings.HasPrefix(t.Value, "'") && strings.HasSuffix(t.Value, "'")
		}, result: token.TokenType_Character},
		{condition: func(t *token.Token) bool {
			_, err := strconv.ParseFloat(t.Value, 64)
			return t.Length >= 1 && err == nil
		}, result: token.TokenType_Number},
		{condition: func(t *token.Token) bool {
			_, ok := luna.AsBooleanValue(t.Value)
			return t.Length >= 1 && ok
		}, result: token.TokenType_Boolean},
		{condition: func(t *token.Token) bool { return t.Value == "nothing" }, result: token.TokenType_Null},
	}

	for partialTokens.Len() > 0 {
		t := partialTokens.Dequeue().Value
		t.Type = token.TokenType_Identifier

		for _, processor := range tokenTypeProcessors {
			if !processor.condition(t) {
				continue
			}

			t.Type = processor.result
			break
		}

		if t.Length == 0 && t.Type == token.TokenType_Identifier {
			continue
		}

		tokens.Queue(t)
	}

	lastToken := tokens.Last()
	startIndex := 0
	if lastToken != nil {
		startIndex = lastToken.Value.Start
	}

	tokens.Queue(&token.Token{
		Start:  startIndex,
		Length: 0,
		Value:  "",
		Type:   token.TokenType_EndOfFile,
	})

	return tokens
}

// Merge full tokens that span across multiple tokens.
func mergeMultiTokens(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	multiTokenProcessors := []struct {
		condition func(tokens *queue.Queue[*token.Token]) int
		result    token.TokenType
	}{
		{condition: parsers.CheckReportHeader, result: token.TokenType_ReportHeader},
		{condition: parsers.CheckReportFooter, result: token.TokenType_ReportFooter},

		{condition: parsers.CheckFunctionHeaderMain, result: token.TokenType_FunctionMain},
		{condition: parsers.CheckFunctionHeader, result: token.TokenType_FunctionHeader},
		{condition: parsers.CheckFunctionFooter, result: token.TokenType_FunctionFooter},
		{condition: parsers.CheckFunctionParameter, result: token.TokenType_FunctionParameter},
		{condition: parsers.CheckFunctionReturn, result: token.TokenType_FunctionReturn},

		{condition: parsers.CheckPrintMethod, result: token.TokenType_Print},
		{condition: parsers.CheckPrintNewlineMethod, result: token.TokenType_PrintNewline},
		{condition: parsers.CheckReadMethod, result: token.TokenType_Prompt},
		{condition: parsers.CheckFunctionCallMethod, result: token.TokenType_FunctionCall},

		{condition: parsers.CheckVariableDeclaration, result: token.TokenType_Declaration},
		{condition: parsers.CheckVariableModifier, result: token.TokenType_Modify},

		{condition: parsers.CheckBooleanType, result: token.TokenType_TypeBoolean},
		{condition: parsers.CheckBooleanArrayType, result: token.TokenType_TypeBooleanArray},
		{condition: parsers.CheckNumberType, result: token.TokenType_TypeNumber},
		{condition: parsers.CheckNumberArrayType, result: token.TokenType_TypeNumberArray},
		{condition: parsers.CheckStringType, result: token.TokenType_TypeString},
		{condition: parsers.CheckStringArrayType, result: token.TokenType_TypeStringArray},
		{condition: parsers.CheckCharacterType, result: token.TokenType_TypeChar},

		{condition: parsers.CheckPostscript, result: token.TokenType_CommentPostScript},

		{condition: parsers.CheckIfKeyword, result: token.TokenType_IfClause},
		{condition: parsers.CheckElseKeyword, result: token.TokenType_ElseClause},
		{condition: parsers.CheckIfEndKeyword, result: token.TokenType_IfEndClause},

		{condition: parsers.CheckWhileKeyword, result: token.TokenType_WhileClause},
		{condition: parsers.CheckForEveryKeyword, result: token.TokenType_ForEveryClause},
		{condition: parsers.CheckStatementEndKeyword, result: token.TokenType_KeywordStatementEnd},

		{condition: parsers.CheckInfixAddition, result: token.TokenType_OperatorAddInfix},
		{condition: parsers.CheckPrefixAddition, result: token.TokenType_OperatorAddPrefix},
		{condition: parsers.CheckInfixSubtraction, result: token.TokenType_OperatorSubInfix},
		{condition: parsers.CheckPrefixSubtraction, result: token.TokenType_OperatorSubPrefix},
		{condition: parsers.CheckInfixMultiplication, result: token.TokenType_OperatorMulInfix},
		{condition: parsers.CheckPrefixMultiplication, result: token.TokenType_OperatorMulPrefix},
		{condition: parsers.CheckInfixDivision, result: token.TokenType_OperatorDivInfix},
		{condition: parsers.CheckPrefixDivision, result: token.TokenType_OperatorDivPrefix},

		{condition: parsers.CheckPrefixModulo, result: token.TokenType_OperatorModPrefix},
		{condition: parsers.CheckInfixModulo, result: token.TokenType_OperatorModInfix},

		{condition: parsers.CheckLessThanEqualOperator, result: token.TokenType_OperatorLte},
		{condition: parsers.CheckGreaterThanEqualOperator, result: token.TokenType_OperatorGte},
		{condition: parsers.CheckGreaterThanOperator, result: token.TokenType_OperatorGt},
		{condition: parsers.CheckLessThanOperator, result: token.TokenType_OperatorLt},
		{condition: parsers.CheckNotEqualOperator, result: token.TokenType_OperatorNeq},
		{condition: parsers.CheckEqualOperator, result: token.TokenType_OperatorEq},

		{condition: parsers.CheckUnaryPrefixIncrement, result: token.TokenType_UnaryIncrementPrefix},
		{condition: parsers.CheckUnaryPostfixIncrement, result: token.TokenType_UnaryIncrementPostfix},
		{condition: parsers.CheckUnaryPrefixDecrement, result: token.TokenType_UnaryDecrementPrefix},
		{condition: parsers.CheckUnaryPostfixDecrement, result: token.TokenType_UnaryDecrementPostfix},

		{condition: parsers.CheckReturnKeyword, result: token.TokenType_KeywordReturn},
		{condition: parsers.CheckConstantKeyword, result: token.TokenType_KeywordConst},
		{condition: parsers.CheckThenKeyword, result: token.TokenType_KeywordThen},
		{condition: parsers.CheckAndKeyword, result: token.TokenType_KeywordAnd},
		{condition: parsers.CheckOrKeyword, result: token.TokenType_KeywordOr},
		{condition: parsers.CheckOfKeyword, result: token.TokenType_KeywordOf},
		{condition: parsers.CheckInKeyword, result: token.TokenType_KeywordIn},
		{condition: parsers.CheckFromKeyword, result: token.TokenType_KeywordFrom},
		{condition: parsers.CheckToKeyword, result: token.TokenType_KeywordTo},
	}

	for oldTokens.Len() > 0 {
		for _, processor := range multiTokenProcessors {
			if oldTokens.Len() <= 0 {
				break
			}

			if oldTokens.First().Value.Type != token.TokenType_Identifier {
				continue
			}

			amount := processor.condition(oldTokens)
			if amount <= 0 {
				continue
			}

			token := utilities.MergeTokens(oldTokens, amount)
			token.Type = processor.result

			oldTokens.QueueFront(token)

		}

		tokens.Queue(oldTokens.Dequeue().Value)
	}

	return tokens
}

// This will turn specific tokens into identifier tokens if conditions has not been met
func smartIdentifierTokens(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	isForEvery := false

	for oldTokens.Len() > 0 {
		t := oldTokens.Dequeue().Value

		if t.Type == token.TokenType_ForEveryClause {
			isForEvery = true
		}

		if !isForEvery {
			if t.Type == token.TokenType_KeywordIn {
				t.Type = token.TokenType_Identifier
			}
			if t.Type == token.TokenType_KeywordFrom {
				t.Type = token.TokenType_Identifier
			}
			if t.Type == token.TokenType_KeywordTo {
				t.Type = token.TokenType_Identifier
			}
		}

		if t.Type == token.TokenType_Punctuation || t.Type == token.TokenType_NewLine {
			isForEvery = false
		}

		tokens.Queue(t)
	}

	return tokens
}

// Combine any identifier tokens that are queued against one another into just one partial
// token instead.
func mergeIdentifiers(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	var identifierToken *token.Token
	for oldTokens.Len() > 0 {
		t := oldTokens.Dequeue().Value

		if t.Type != token.TokenType_Identifier {
			if identifierToken != nil &&
				t.Type == token.TokenType_Whitespace &&
				oldTokens.Len() >= 1 &&
				oldTokens.First().Value.Type == token.TokenType_Identifier {
				identifierToken.Append(t)
				continue
			}

			if identifierToken != nil {
				tokens.Queue(identifierToken)
				identifierToken = nil
			}
			tokens.Queue(t)
			continue
		}

		if identifierToken != nil {
			identifierToken.Append(t)
			continue
		}
		identifierToken = t
	}

	// Flush remaining
	if identifierToken != nil {
		tokens.Queue(identifierToken)
		identifierToken = nil
	}

	return tokens
}

// Remove any unnecessary tokens from the token queue
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
			for oldTokens.Len() > 0 &&
				(oldTokens.First().Value.Type != token.TokenType_NewLine &&
					oldTokens.First().Value.Type != token.TokenType_EndOfFile) {
				oldTokens.Dequeue()
			}
			continue
		}

		tokens.Queue(t)
	}

	return tokens
}
