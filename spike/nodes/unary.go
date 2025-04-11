package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type UnaryExpressionNode struct {
	Node

	Identifier string
	Increment bool
}

func (b *UnaryExpressionNode) Type() NodeType {
	return TYPE_UNARYEXPRESSION
}
func (b *UnaryExpressionNode) ToNode() Node {
	return Node{
		Start:  b.Start,
		Length: b.Length,
	}
}

func ParsePrefixUnary(ast *ast.AST) (*UnaryExpressionNode, error) {
	unaryNode := &UnaryExpressionNode{}

	startToken, err := ast.ConsumeFunc(func(t *token.Token) bool {
		return t.Type == token.TokenType_UnaryIncrementPrefix || t.Type == token.TokenType_UnaryDecrementPrefix
	}, "Expected unary prefix token")
	unaryNode.Increment = startToken.Type == token.TokenType_UnaryIncrementPrefix

	identifierNode, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	unaryNode.Identifier = identifierNode.Value

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	unaryNode.Start = startToken.Start
	unaryNode.Length = endToken.Start + endToken.Length - startToken.Start

	return unaryNode, nil
}

func ParsePostfixUnary(ast *ast.AST) (*UnaryExpressionNode, error) {
	unaryNode := &UnaryExpressionNode{}

	identifierNode, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	unaryNode.Identifier = identifierNode.Value

	postfixToken, err := ast.ConsumeFunc(func(t *token.Token) bool {
		return t.Type == token.TokenType_UnaryIncrementPostfix || t.Type == token.TokenType_UnaryDecrementPostfix
	}, "Expected unary postfix token")
	unaryNode.Increment = postfixToken.Type == token.TokenType_UnaryIncrementPostfix

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	unaryNode.Start = identifierNode.Start
	unaryNode.Length = endToken.Start + endToken.Length - identifierNode.Start

	return unaryNode, nil
}
