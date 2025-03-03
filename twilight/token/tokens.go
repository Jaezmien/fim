package token

type TokenType uint

const (
	// 0
	TokenType_Unknown TokenType = iota

	// 1
	TokenType_Literal
	TokenType_Punctuation
	TokenType_NewLine
	TokenType_Whitespace
	TokenType_EndOfFile

	// 6
	TokenType_CommentParen
	TokenType_CommentPostScript

	// 8
	TokenType_String
	TokenType_Character
	TokenType_Number
	TokenType_Boolean

	// 12
	TokenType_Null

	// 13
	TokenType_ReportHeader
	TokenType_ReportFooter

	// 15
	TokenType_FunctionHeader
	TokenType_FunctionMain
	TokenType_FunctionFooter

	// 18
	TokenType_FunctionReturn
	TokenType_FunctionParameter

	// 20
	TokenType_Print
	TokenType_PrintNewline
	TokenType_Prompt
	TokenType_FunctionCall

	// 24
	TokenType_Declaration
	TokenType_Modify

	// 26
	TokenType_TypeString
	TokenType_TypeChar
	TokenType_TypeNumber
	TokenType_TypeBoolean

	// 30
	TokenType_TypeStringArray
	TokenType_TypeNumberArray
	TokenType_TypeBooleanArray

	// 33
	TokenType_OperatorEq
	TokenType_OperatorNeq
	TokenType_OperatorGt
	TokenType_OperatorGte
	TokenType_OperatorLt
	TokenType_OperatorLte

	// 39
	TokenType_UnaryNot

	// 40
	TokenType_OperatorAddInfix
	TokenType_OperatorAddPrefix
	TokenType_UnaryIncrement

	// 43
	TokenType_OperatorSubInfix
	TokenType_OperatorSubPrefix
	TokenType_UnaryDecrement

	// 46
	TokenType_OperatorMulInfix
	TokenType_OperatorMulPrefix

	// 48
	TokenType_OperatorDivInfix
	TokenType_OperatorDivPrefix

	// 50
	TokenType_KeywordOr
	TokenType_KeywordAnd
	TokenType_KeywordConst
	TokenType_KeywordOf
	TokenType_KeywordThen
	TokenType_KeywordStatementEnd
	TokenType_KeywordReturn

	// 57
	TokenType_IfClause
	TokenType_ElseClause
	TokenType_IfEndClause

	// 60
	TokenType_WhileClause	
)

type Token struct {
	Start  int
	Length int
	Type   TokenType
	Value string
}
func (t *Token) Append(token *Token) {
	t.Value += token.Value
	t.Length += token.Length
}
