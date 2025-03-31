package ast

import (
	"slices"

	luna "git.jaezmien.com/Jaezmien/fim/luna/utilities"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type AST struct {
	Tokens     []*token.Token
	TokenIndex int

	Source string
}

func NewAST(tokens []*token.Token, source string) *AST {
	return &AST{
		Tokens:     tokens,
		TokenIndex: 0,
		Source:     source,
	}
}

func (a *AST) PeekAt(index int) *token.Token {
	if index < 0 || index >= len(a.Tokens) {
		return &token.Token{
			Start:  0,
			Length: 0,
			Type:   token.TokenType_Unknown,
		}
	}
	return a.Tokens[index]
}
func (a *AST) Peek() *token.Token {
	return a.PeekAt(a.TokenIndex)
}
func (a *AST) PeekNext() *token.Token {
	return a.PeekAt(a.TokenIndex + 1)
}
func (a *AST) PeekPrevious() *token.Token {
	return a.PeekAt(a.TokenIndex - 1)
}
func (a *AST) PeekIndex() int {
	return a.TokenIndex
}

func (a *AST) Next() {
	a.TokenIndex += 1
}
func (a *AST) MoveTo(index int) {
	a.TokenIndex = index
}

func (a *AST) EndOfFile() bool {
	current := a.Peek()
	return current.Type == token.TokenType_EndOfFile
}

func (a *AST) CreateErrorFromIndex(index int, errorMessage string) error {
	return luna.CreateErrorFromIndex(a.Source, index, errorMessage)
}
func (a *AST) CreateErrorFromToken(t *token.Token, errorMessage string) error {
	return a.CreateErrorFromIndex(t.Start, errorMessage)
}

func (a *AST) GetSourceText(start int, length int) string {
	return a.Source[start : start+length]
}

func (a *AST) CheckType(tokenTypes ...token.TokenType) bool {
	current := a.Peek()
	return slices.Contains(tokenTypes, current.Type)
}

func (a *AST) Contains(tokenType token.TokenType) bool {
	for idx := a.PeekIndex(); idx < len(a.Tokens); idx++ {
		if a.PeekAt(idx).Type == tokenType {
			return true
		}
	}
	return false
}
func (a *AST) ContainsWithStop(tokenType token.TokenType, stopTokens ...token.TokenType) bool {
	for idx := a.PeekIndex(); idx < len(a.Tokens); idx++ {
		current := a.PeekAt(idx)

		if slices.Contains(stopTokens, current.Type) {
			break
		}
		if current.Type == tokenType {
			return true
		}
	}
	return false
}

func (a *AST) ConsumeFunc(predicate func(*token.Token) bool, errorMessage string) (*token.Token, error) {
	current := a.Peek()

	if !predicate(current) {
		return nil, a.CreateErrorFromToken(current, errorMessage)
	}

	a.Next()

	return a.PeekPrevious(), nil
}
func (a *AST) ConsumeToken(tokenType token.TokenType, errorMessage string) (*token.Token, error) {
	return a.ConsumeFunc(func(t *token.Token) bool {
		return t.Type == tokenType
	}, errorMessage)
}

func (a *AST) ConsumeUntilFuncMatch(predicate func(*token.Token) bool, errorMessage string) ([]*token.Token, error) {
	tokens := make([]*token.Token, 0)

	for {
		current := a.Peek()

		if a.EndOfFile() {
			return nil, a.CreateErrorFromToken(current, errorMessage)
		}

		if predicate(current) {
			break
		}

		tokens = append(tokens, current)
		a.Next()
	}

	return tokens, nil
}
func (a *AST) ConsumeUntilTokenMatch(tokenType token.TokenType, errorMessage string) ([]*token.Token, error) {
	return a.ConsumeUntilFuncMatch(func(t *token.Token) bool {
		return t.Type == tokenType
	}, errorMessage)
}
