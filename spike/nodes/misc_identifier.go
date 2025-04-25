package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	. "git.jaezmien.com/Jaezmien/fim/spike/node"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type IdentifierNode struct {
	Node

	Identifier string
}

type DictionaryIdentifierNode struct {
	Node

	Identifier string
	Index      DynamicNode
}

func ParseIdentifierNode(ast *ast.AST) (DynamicNode, error) {
	if ast.ContainsWithStop(token.TokenType_KeywordOf, token.TokenType_EndOfFile, token.TokenType_Punctuation) {
		node := &DictionaryIdentifierNode{}

		indexToken, err := ast.ConsumeUntilTokenMatch(token.TokenType_KeywordOf, token.TokenType_KeywordOf.Message("Expected %s"))
		if err != nil {
			return nil, err
		}
		indexNode, err := CreateValueNode(indexToken, CreateValueNodeOptions{})
		if err != nil {
			return nil, err
		}
		node.Index = indexNode

		_, err = ast.ConsumeToken(token.TokenType_KeywordOf, token.TokenType_KeywordOf.Message("Expected %s"))
		if err != nil {
			return nil, err
		}

		identifierToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
		if err != nil {
			return nil, err
		}
		node.Identifier = identifierToken.Value

		node.Start = indexNode.ToNode().Start
		node.Length = identifierToken.Start + identifierToken.Length - indexNode.ToNode().Start

		return node, nil
	} else {
		node := &IdentifierNode{}

		identifierToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
		if err != nil {
			return nil, err
		}
		node.Identifier = identifierToken.Value

		node.Start = identifierToken.Start
		node.Length = identifierToken.Length

		return node, nil
	}
}
