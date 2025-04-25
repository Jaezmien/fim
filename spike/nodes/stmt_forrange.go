package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type ForEveryRangeStatementNode struct {
	ForEveryStatementNode

	RangeStart DynamicNode
	RangeEnd   DynamicNode
}

func ParseForEveryRangeStatementNode(ast *ast.AST) (*ForEveryRangeStatementNode, error) {
	node := &ForEveryRangeStatementNode{}

	startToken, err := ast.ConsumeToken(token.TokenType_ForEveryClause, token.TokenType_ForEveryClause.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	typeToken := ast.Consume()
	possibleType := variable.FromTokenTypeHint(typeToken.Type)
	if possibleType == variable.UNKNOWN || possibleType.IsArray() {
		return nil, typeToken.CreateError("Expected variable type", ast.Source)
	}
	node.VariableType = possibleType

	identifierToken, err := ast.ConsumeToken(token.TokenType_Identifier, token.TokenType_Identifier.Message("Expected %s"))
	if err != nil {
		return nil, err
	}
	node.VariableName = identifierToken.Value

	_, err = ast.ConsumeToken(token.TokenType_KeywordFrom, token.TokenType_KeywordFrom.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	fromTokens, err := ast.ConsumeUntilTokenMatch(token.TokenType_KeywordTo, "Expected 'to' token")
	if err != nil {
		return nil, err
	}
	fromNode, err := CreateValueNode(fromTokens, CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}
	node.RangeStart = fromNode

	_, err = ast.ConsumeToken(token.TokenType_KeywordTo, token.TokenType_KeywordTo.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	toTokens, err := ast.ConsumeUntilTokenMatch(token.TokenType_Punctuation, "Expected ")
	if err != nil {
		return nil, err
	}
	toNode, err := CreateValueNode(toTokens, CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}
	node.RangeEnd = toNode

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
