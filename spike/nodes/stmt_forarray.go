package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type ForEveryArrayStatementNode struct {
	ForEveryStatementNode

	Identifier string
}

func ParseForEveryArrayStatementNode(ast *ast.AST) (*ForEveryArrayStatementNode, error) {
	node := &ForEveryArrayStatementNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_ForEveryClause, token.TokenType_ForEveryClause.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	typeToken := ast.Consume()
	possibleType := variable.FromTokenTypeHint(typeToken.Type)
	if possibleType == variable.UNKNOWN || possibleType.IsArray() {
		return nil, typeToken.CreateError("Expected non-array variable type", ast.Source)
	}
	node.VariableType = possibleType

	variableNameToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	node.VariableName = variableNameToken.Value

	_, err = ast.ConsumeToken(token.TokenType_KeywordIn, token.TokenType_KeywordIn.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	identifierToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	node.Identifier = identifierToken.Value

	_, err = ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	statements, err := ParseStatementsNode(ast, token.TokenType_KeywordStatementEnd)
	if err != nil {
		return nil, err
	}
	node.StatementsNode = *statements

	_, err = ast.ConsumeToken(token.TokenType_KeywordStatementEnd, token.TokenType_KeywordStatementEnd.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	endToken, err := ast.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	node.Start = startToken.Start
	node.Length = endToken.Start + endToken.Length - startToken.Start

	return node, nil
}
