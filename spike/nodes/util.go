package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

// Consumes until a punctuation is found, unless if it detects a function parameter keyword
// wherein it consumes punctuations *except* commas.
//
// isArray is used when we know that we're combing through an array values.
// (e.g. Declaring a variable, and we already know that this is an array type.)
//
// See: var_declare.go
func ConsumeUntilPunctuation(ast *ast.AST, isArray bool) (valueTokens []*token.Token, err error) {
	if isArray || ast.ContainsWithStop(token.TokenType_FunctionParameter, token.TokenType_EndOfFile, token.TokenType_Punctuation) {
		valueTokens, err = ast.ConsumeUntilFuncMatch(func(t *token.Token) bool {
			return t.Type == token.TokenType_Punctuation && t.Value != ","
		}, token.TokenType_Punctuation.Message("Could not find %s"))
	} else {
		valueTokens, err = ast.ConsumeUntilTokenMatch(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Could not find %s"))
	}

	return valueTokens, err
}
