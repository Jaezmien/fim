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
