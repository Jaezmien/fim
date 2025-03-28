package twilight

import (
	"slices"
	"strconv"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/parsers"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func createTokens(partialTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	punctuations := [...]rune{'.', '!', '?', ':', ','}
	booleanStrings := [...]string{"yes", "true", "right", "correct", "no", "false", "wrong", "incorrect"}

	tokenTypeProcessors := []struct {
		condition func(t *token.Token) bool
		result    token.TokenType
	}{
		{condition: func(t *token.Token) bool {
			return t.Length == 1 && utilities.ContainsRune(rune(t.Value[0]), punctuations[:])
		}, result: token.TokenType_Punctuation},
		{condition: func(t *token.Token) bool { return t.Length == 1 && t.Value == "\n" }, result: token.TokenType_NewLine},
		{condition: func(t *token.Token) bool { return t.Length == 1 && t.Value == " " }, result: token.TokenType_Whitespace},
		{condition: func(t *token.Token) bool {
			return t.Length == 1 && strings.HasPrefix(t.Value, "(") && strings.HasSuffix(t.Value, ")")
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
		{condition: func(t *token.Token) bool { return t.Length >= 1 && slices.Contains(booleanStrings[:], t.Value) }, result: token.TokenType_Boolean},
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

func mergeMultitokens(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	multiTokenProcessors := []struct{
		condition func(tokens *queue.Queue[*token.Token]) int
		result token.TokenType
	} {
		{condition: parsers.IsReportHeader, result: token.TokenType_ReportHeader},
		{condition: parsers.IsReportFooter, result: token.TokenType_ReportFooter},

		{condition: parsers.IsFunctionHeaderMain, result: token.TokenType_FunctionMain},
		{condition: parsers.IsFunctionHeader, result: token.TokenType_FunctionHeader},
		{condition: parsers.IsFunctionFooter, result: token.TokenType_FunctionFooter},
		{condition: parsers.IsFunctionParameter, result: token.TokenType_FunctionParameter},
		{condition: parsers.IsFunctionReturn, result: token.TokenType_FunctionReturn},

		{condition: parsers.IsPrintMethod, result: token.TokenType_Print},
		{condition: parsers.IsPrintNewlineMethod, result: token.TokenType_PrintNewline},
		{condition: parsers.IsReadMethod, result: token.TokenType_Prompt},
		{condition: parsers.IsFunctionCallMethod, result: token.TokenType_FunctionCall},

		{condition: parsers.IsVariableDeclaration, result: token.TokenType_Declaration},
		{condition: parsers.IsVariableModifier, result: token.TokenType_Modify},

		{condition: parsers.IsBooleanType, result: token.TokenType_TypeBoolean},
		{condition: parsers.IsBooleanArrayType, result: token.TokenType_TypeBooleanArray},
		{condition: parsers.IsNumberType, result: token.TokenType_TypeNumber},
		{condition: parsers.IsNumberArrayType, result: token.TokenType_TypeNumberArray},
		{condition: parsers.IsStringType, result: token.TokenType_TypeString},
		{condition: parsers.IsStringArrayType, result: token.TokenType_TypeStringArray},
		{condition: parsers.IsCharacterType, result: token.TokenType_TypeChar},

		{condition: parsers.IsPostscript, result: token.TokenType_CommentPostScript},

		{condition: parsers.IsInfixAddition, result: token.TokenType_OperatorAddInfix},
		{condition: parsers.IsPrefixAddition, result: token.TokenType_OperatorAddPrefix},
		{condition: parsers.IsInfixSubtraction, result: token.TokenType_OperatorSubInfix},
		{condition: parsers.IsPrefixSubtraction, result: token.TokenType_OperatorSubPrefix},
		{condition: parsers.IsInfixMultiplication, result: token.TokenType_OperatorMulInfix},
		{condition: parsers.IsPrefixMultiplication, result: token.TokenType_OperatorMulPrefix},
		{condition: parsers.IsInfixDivision, result: token.TokenType_OperatorDivInfix},
		{condition: parsers.IsPrefixDivision, result: token.TokenType_OperatorDivPrefix},

		{condition: parsers.IsLessThanEqualOperator, result: token.TokenType_OperatorLte},
		{condition: parsers.IsGreaterThanEqualOperator, result: token.TokenType_OperatorGte},
		{condition: parsers.IsGreaterThanOperator, result: token.TokenType_OperatorGt},
		{condition: parsers.IsLessThanOperator, result: token.TokenType_OperatorLt},
		{condition: parsers.IsNotEqualOperator, result: token.TokenType_OperatorNeq},
		{condition: parsers.IsEqualOperator, result: token.TokenType_OperatorEq},

		{condition: parsers.IsConstantKeyword, result: token.TokenType_KeywordConst},
		{condition: parsers.IsAndKeyword, result: token.TokenType_KeywordAnd},
		{condition: parsers.IsOrKeyword, result: token.TokenType_KeywordOr},
		{condition: parsers.IsOfKeyword, result: token.TokenType_KeywordOf},
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

func mergeLiterals(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	var literalToken *token.Token
	for oldTokens.Len() > 0 {
		t := oldTokens.Dequeue().Value

		if t.Type != token.TokenType_Identifier {
			if literalToken != nil &&
				t.Type == token.TokenType_Whitespace &&
				oldTokens.Len() >= 1 &&
				oldTokens.First().Value.Type == token.TokenType_Identifier {
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
