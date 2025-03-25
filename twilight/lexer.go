package twilight

import (
	"slices"
	"strconv"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/twilight/parsers"
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func createTokens(partialTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	punctuations := [...]rune{'.', '!', '?', ':', ','}
	isRunePunctuation := func(r rune) bool {
		return utilities.ContainsRune(r, punctuations[:])
	}

	booleanStrings := [...]string{"yes", "true", "right", "correct", "no", "false", "wrong", "incorrect"}

	processTokenType := func(t *token.Token, condition bool, resultType token.TokenType) {
		if t.Type != token.TokenType_Identifier {
			return
		}

		if !condition {
			return
		}

		t.Type = resultType
	}

	for partialTokens.Len() > 0 {
		t := partialTokens.Dequeue().Value
		t.Type = token.TokenType_Identifier

		processTokenType(t, t.Length == 1 && isRunePunctuation(rune(t.Value[0])), token.TokenType_Punctuation)
		processTokenType(t, t.Length == 1 && t.Value == "\n", token.TokenType_NewLine)
		processTokenType(t, t.Length == 1 && t.Value == " ", token.TokenType_Whitespace)
		processTokenType(t, t.Length == 1 && strings.HasPrefix(t.Value, "(") && strings.HasSuffix(t.Value, ")"), token.TokenType_CommentParen)

		processTokenType(t, t.Length >= 1 && strings.HasPrefix(t.Value, "\"") && strings.HasSuffix(t.Value, "\""), token.TokenType_String)
		processTokenType(t, t.Length >= 1 && strings.HasPrefix(t.Value, "'") && strings.HasSuffix(t.Value, "'"), token.TokenType_Character)
		if _, err := strconv.ParseFloat(t.Value, 64); err == nil {
			processTokenType(t, t.Length >= 1, token.TokenType_Number)
		}
		processTokenType(t, t.Length >= 1 && slices.Contains(booleanStrings[:], t.Value), token.TokenType_Boolean)
		processTokenType(t, t.Value == "nothing", token.TokenType_Null)

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

type processMultiTokenResult = func(tokens *queue.Queue[*token.Token]) int

func mergeMultitokens(oldTokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	tokens := queue.New[*token.Token]()

	processMultiTokenType := func(condition processMultiTokenResult, resultType token.TokenType) {
		if oldTokens.Len() <= 0 {
			return
		}

		if oldTokens.First().Value.Type != token.TokenType_Identifier {
			return
		}

		amount := condition(oldTokens)
		if amount <= 0 {
			return
		}

		token := utilities.MergeTokens(oldTokens, amount)
		token.Type = resultType

		oldTokens.QueueFront(token)
	}

	for oldTokens.Len() > 0 {
		processMultiTokenType(parsers.IsReportHeader, token.TokenType_ReportHeader)
		processMultiTokenType(parsers.IsReportFooter, token.TokenType_ReportFooter)

		processMultiTokenType(parsers.IsFunctionHeaderMain, token.TokenType_FunctionMain)
		processMultiTokenType(parsers.IsFunctionHeader, token.TokenType_FunctionHeader)
		processMultiTokenType(parsers.IsFunctionFooter, token.TokenType_FunctionFooter)
		processMultiTokenType(parsers.IsFunctionParameter, token.TokenType_FunctionParameter)
		processMultiTokenType(parsers.IsFunctionReturn, token.TokenType_FunctionReturn)

		processMultiTokenType(parsers.IsPrintMethod, token.TokenType_Print)
		processMultiTokenType(parsers.IsPrintNewlineMethod, token.TokenType_PrintNewline)
		processMultiTokenType(parsers.IsReadMethod, token.TokenType_Prompt)
		processMultiTokenType(parsers.IsFunctionCallMethod, token.TokenType_FunctionCall)

		processMultiTokenType(parsers.IsVariableDeclaration, token.TokenType_Declaration)
		processMultiTokenType(parsers.IsVariableModifier, token.TokenType_Modify)

		processMultiTokenType(parsers.IsBooleanType, token.TokenType_TypeBoolean);
		processMultiTokenType(parsers.IsBooleanArrayType, token.TokenType_TypeBooleanArray);
		processMultiTokenType(parsers.IsNumberType, token.TokenType_TypeNumber);
		processMultiTokenType(parsers.IsNumberArrayType, token.TokenType_TypeNumberArray);
		processMultiTokenType(parsers.IsStringType, token.TokenType_TypeString);
		processMultiTokenType(parsers.IsStringArrayType, token.TokenType_TypeStringArray);
		processMultiTokenType(parsers.IsCharacterType, token.TokenType_TypeChar);

		processMultiTokenType(parsers.IsPostscript, token.TokenType_CommentPostScript)

		processMultiTokenType(parsers.IsInfixAddition, token.TokenType_OperatorAddInfix)
		processMultiTokenType(parsers.IsPrefixAddition, token.TokenType_OperatorAddPrefix)
		processMultiTokenType(parsers.IsInfixSubtraction, token.TokenType_OperatorSubInfix)
		processMultiTokenType(parsers.IsPrefixSubtraction, token.TokenType_OperatorSubPrefix)
		processMultiTokenType(parsers.IsInfixMultiplication, token.TokenType_OperatorMulInfix)
		processMultiTokenType(parsers.IsPrefixMultiplication, token.TokenType_OperatorMulPrefix)
		processMultiTokenType(parsers.IsInfixDivision, token.TokenType_OperatorDivInfix)
		processMultiTokenType(parsers.IsPrefixDivision, token.TokenType_OperatorDivPrefix)

		processMultiTokenType(parsers.IsLessThanEqualOperator, token.TokenType_OperatorLte)
		processMultiTokenType(parsers.IsGreaterThanEqualOperator, token.TokenType_OperatorGte)
		processMultiTokenType(parsers.IsGreaterThanOperator, token.TokenType_OperatorGt)
		processMultiTokenType(parsers.IsLessThanOperator, token.TokenType_OperatorLt)
		processMultiTokenType(parsers.IsNotEqualOperator, token.TokenType_OperatorNeq)
		processMultiTokenType(parsers.IsEqualOperator, token.TokenType_OperatorEq)

		processMultiTokenType(parsers.IsConstantKeyword, token.TokenType_KeywordConst)
		processMultiTokenType(parsers.IsAndKeyword, token.TokenType_KeywordAnd)
		processMultiTokenType(parsers.IsOrKeyword, token.TokenType_KeywordOr)
		processMultiTokenType(parsers.IsOfKeyword, token.TokenType_KeywordOf)

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
