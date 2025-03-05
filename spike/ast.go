package spike

import (
	"errors"
	"fmt"
	"slices"

	"git.jaezmien.com/Jaezmien/fim/spike/utilities"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type AST struct {
	tokens     []*token.Token
	tokenIndex int

	source string
}

func NewAST(tokens []*token.Token, source string) *AST {
	return &AST{
		tokens:     tokens,
		tokenIndex: 0,
		source:     source,
	}
}

func (a *AST) PeekAt(index int) *token.Token {
	if index < 0 || index >= len(a.tokens) {
		return nil
	}
	return a.tokens[index]
}
func (a *AST) Peek() *token.Token {
	return a.PeekAt(a.tokenIndex)
}
func (a *AST) PeekNext() *token.Token {
	return a.PeekAt(a.tokenIndex + 1)
}
func (a *AST) PeekPrevious() *token.Token {
	return a.PeekAt(a.tokenIndex - 1)
}
func (a *AST) PeekIndex() int {
	return a.tokenIndex
}

func (a *AST) Next() {
	a.tokenIndex += 1
}
func (a *AST) MoveTo(index int) {
	a.tokenIndex = index
}

func (a *AST) EndOfFile() bool {
	current := a.Peek()
	if current == nil {
		panic("AST@EndOfFile called with nil token")
	}
	return current.Type == token.TokenType_EndOfFile
}

func (a *AST) CreateErrorFromIndex(index int, errorMessage string) error {
	pair := utilities.GetErrorIndexPair(a.source, index)
	return errors.New(fmt.Sprintf("[line: %d] %s", pair.Line, errorMessage))
}
func (a *AST) CreateErrorFromToken(t *token.Token, errorMessage string) error {
	return a.CreateErrorFromIndex(t.Start, errorMessage)
}

func (a *AST) GetSourceText(start int, length int) string {
	return a.source[start : start+length]
}

func (a *AST) CheckType(tokenTypes ...token.TokenType) bool {
	current := a.Peek()
	if current == nil {
		panic("AST@CheckType called with nil token")
	}
	return slices.Contains(tokenTypes, current.Type)
}

func (a *AST) Contains(tokenType token.TokenType) bool {
	for idx := a.PeekIndex(); idx < len(a.tokens); idx++ {
		if a.PeekAt(idx).Type == tokenType {
			return true
		}
	}
	return false
}
func (a *AST) ContainsWithStop(tokenType token.TokenType, stopTokens ...token.TokenType) bool {
	for idx := a.PeekIndex(); idx < len(a.tokens); idx++ {
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

	if current == nil {
		panic("AST@ConsumeFunc called with nil token")
	}

	if !predicate(current) {
		return nil, a.CreateErrorFromToken(current, errorMessage)
	}

	if a.PeekNext() == nil {
		return nil, a.CreateErrorFromToken(current, "Reached END_OF_FILE")
	}
	a.Next()

	return a.PeekPrevious(), nil
}
func (a *AST) ConsumeToken(tokenType token.TokenType, errorMessage string) (*token.Token, error) {
	return a.ConsumeFunc(func(t *token.Token) bool {
		return t.Type == tokenType
	}, errorMessage)
}

func (a *AST) ConsumeFuncUntilMatch(predicate func(*token.Token) bool, errorMessage string) ([]*token.Token, error) {
	tokens := make([]*token.Token, 0)

	for {
		current := a.Peek()
		if current == nil {
			panic("AST@ConsumeFuncUntilMatch called with nil token")
		}

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
func (a *AST) ConsumeTokenUntilMatch(tokenType token.TokenType, errorMessage string) ([]*token.Token, error) {
	return a.ConsumeFuncUntilMatch(func(t *token.Token) bool {
		return t.Type == tokenType
	}, errorMessage)
}
