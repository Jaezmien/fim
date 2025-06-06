package token

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/luna/errors"
)

type TokenType uint

const (
	TokenType_Unknown TokenType = iota

	TokenType_Identifier
	TokenType_Punctuation
	TokenType_NewLine
	TokenType_Whitespace
	TokenType_EndOfFile

	TokenType_CommentParen
	TokenType_CommentPostScript

	TokenType_String
	TokenType_Character
	TokenType_Number
	TokenType_Boolean

	TokenType_Null

	TokenType_ReportHeader
	TokenType_ReportFooter

	TokenType_FunctionHeader
	TokenType_FunctionMain
	TokenType_FunctionFooter

	TokenType_FunctionReturn
	TokenType_FunctionParameter

	TokenType_Print
	TokenType_PrintNewline
	TokenType_Prompt
	TokenType_FunctionCall

	TokenType_Declaration
	TokenType_Modify

	TokenType_TypeString
	TokenType_TypeChar
	TokenType_TypeNumber
	TokenType_TypeBoolean

	TokenType_TypeStringArray
	TokenType_TypeNumberArray
	TokenType_TypeBooleanArray

	TokenType_OperatorEq
	TokenType_OperatorNeq
	TokenType_OperatorGt
	TokenType_OperatorGte
	TokenType_OperatorLt
	TokenType_OperatorLte

	TokenType_UnaryNot

	TokenType_OperatorAddInfix
	TokenType_OperatorAddPrefix
	TokenType_UnaryIncrementPrefix
	TokenType_UnaryIncrementPostfix

	TokenType_OperatorSubInfix
	TokenType_OperatorSubPrefix
	TokenType_UnaryDecrementPrefix
	TokenType_UnaryDecrementPostfix

	TokenType_OperatorMulInfix
	TokenType_OperatorMulPrefix

	TokenType_OperatorDivInfix
	TokenType_OperatorDivPrefix

	TokenType_OperatorModInfix
	TokenType_OperatorModPrefix

	TokenType_KeywordOr
	TokenType_KeywordAnd
	TokenType_KeywordConst
	TokenType_KeywordOf
	TokenType_KeywordIn
	TokenType_KeywordFrom
	TokenType_KeywordTo
	TokenType_KeywordThen
	TokenType_KeywordStatementEnd
	TokenType_KeywordReturn

	TokenType_IfClause
	TokenType_ElseClause
	TokenType_IfEndClause

	TokenType_WhileClause

	TokenType_ForEveryClause
)

var tokenTypeFriendlyName = map[TokenType]string{
	TokenType_Unknown: "UNKNOWN",

	TokenType_Identifier:  "IDENTIFIER",
	TokenType_Punctuation: "PUNCTUATION",
	TokenType_NewLine:     "NEWLINE",
	TokenType_Whitespace:  "WHITESPACE",
	TokenType_EndOfFile:   "EOF",

	TokenType_CommentParen:      "COMMENT",
	TokenType_CommentPostScript: "COMMENT",

	TokenType_String:    "LITERAL(STRING)",
	TokenType_Character: "LITERAL(CHARACTER)",
	TokenType_Number:    "LITERAL(NUMBER)",
	TokenType_Boolean:   "LITERAL(BOOL)",

	TokenType_Null: "NULL",

	TokenType_ReportHeader: "REPORT(HEADER)",
	TokenType_ReportFooter: "REPORT(FOOTER)",

	TokenType_FunctionHeader: "FUNCTION(HEADER)",
	TokenType_FunctionMain:   "FUNCTION(MAIN)",
	TokenType_FunctionFooter: "FUNCTION(FOOTER)",

	TokenType_FunctionReturn:    "FUNCTION(RETURN)",
	TokenType_FunctionParameter: "FUNCTION(PARAMETER)",

	TokenType_Print:        "PRINT",
	TokenType_PrintNewline: "PRINT(NEWLINE)",
	TokenType_Prompt:       "PROMPT",
	TokenType_FunctionCall: "FUNCTION(CALL)",

	TokenType_Declaration: "VARIABLE(DECLARATION)",
	TokenType_Modify:      "VARIABLE(MODIFY)",

	TokenType_TypeString:  "TYPE(STRING)",
	TokenType_TypeChar:    "TYPE(CHARACTER)",
	TokenType_TypeNumber:  "TYPE(NUMBER)",
	TokenType_TypeBoolean: "TYPE(BOOLEAN)",

	TokenType_TypeStringArray:  "TYPE(STRING_ARRAY)",
	TokenType_TypeNumberArray:  "TYPE(NUMBER_ARRAY)",
	TokenType_TypeBooleanArray: "TYPE(BOOLEAN_ARRAY)",

	TokenType_OperatorEq:  "OPERATOR(EQ)",
	TokenType_OperatorNeq: "OPERATOR(NEQ)",
	TokenType_OperatorGt:  "OPERATOR(GT)",
	TokenType_OperatorGte: "OPERATOR(GTE)",
	TokenType_OperatorLt:  "OPERATOR(LT)",
	TokenType_OperatorLte: "OPERATOR(LTE)",

	TokenType_UnaryNot: "UNARY(NOT)",

	TokenType_OperatorAddInfix:      "OPERATOR(ADD_INFIX)",
	TokenType_OperatorAddPrefix:     "OPERATOR(ADD_PREFIX)",
	TokenType_UnaryIncrementPrefix:  "UNARY(INCREMENT(PREFIX))",
	TokenType_UnaryIncrementPostfix: "UNARY(INCREMENT(POSTFIX))",

	TokenType_OperatorSubInfix:      "OPERATOR(SUB_INFIX)",
	TokenType_OperatorSubPrefix:     "OPERATOR(SUB_PREFIX)",
	TokenType_UnaryDecrementPrefix:  "UNARY(DECREMENT(PREFIX))",
	TokenType_UnaryDecrementPostfix: "UNARY(DECREMENT(POSTFIX))",

	TokenType_OperatorMulInfix:  "OPERATOR(MUL_INFIX)",
	TokenType_OperatorMulPrefix: "OPEREATOR(MUL_PREFIX)",

	TokenType_OperatorDivInfix:  "OPERATOR(DIV_INFIX)",
	TokenType_OperatorDivPrefix: "OPERATOR(DIV_PREFIX)",

	TokenType_OperatorModInfix:  "OPERATOR(MOD_INFIX)",
	TokenType_OperatorModPrefix: "OPERATOR(MOD_PREFIX)",

	TokenType_KeywordOr:           "OR",
	TokenType_KeywordAnd:          "AND",
	TokenType_KeywordConst:        "CONST",
	TokenType_KeywordOf:           "OF",
	TokenType_KeywordThen:         "THEN",
	TokenType_KeywordIn:           "IN",
	TokenType_KeywordFrom:         "FROM",
	TokenType_KeywordTo:           "TO",
	TokenType_KeywordStatementEnd: "STATEMENT(END)",
	TokenType_KeywordReturn:       "RETURN",

	TokenType_IfClause:    "IF",
	TokenType_ElseClause:  "ELSE",
	TokenType_IfEndClause: "IF(END)",

	TokenType_WhileClause: "WHILE",

	TokenType_ForEveryClause: "FOREVERY",
}

func (t TokenType) String() string {
	return tokenTypeFriendlyName[t]
}
func (t TokenType) Message(format string) string {
	return fmt.Sprintf(format, t.String())
}

type Token struct {
	Start  int
	Length int
	Type   TokenType
	Value  string
}

func (t *Token) Append(token *Token) {
	t.Value += token.Value
	t.Length += token.Length
}

func (t *Token) CreateError(msg string, source string) error {
	return errors.NewParseError(msg, source, t.Start)
}
