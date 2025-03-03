package token

type TokenType uint

const (
	TokenType_Unknown TokenType = iota

	TokenType_Literal
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
	TokenType_PrintInline
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
	TokenType_UnaryIncrement

	TokenType_OperatorSubInfix
	TokenType_OperatorSubPrefix
	TokenType_UnaryDecrement

	TokenType_OperatorMulInfix
	TokenType_OperatorMulPrefix

	TokenType_OperatorDivInfix
	TokenType_OperatorDivPrefix

	TokenType_KeywordOr
	TokenType_KeywordAnd
	TokenType_KeywordConst
	TokenType_KeywordOf
	TokenType_KeywordThen
	TokenType_KeywordStatementEnd
	TokenType_KeywordReturn

	TokenType_IfClause
	TokenType_ElseClause
	TokenType_IfEndClause

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
