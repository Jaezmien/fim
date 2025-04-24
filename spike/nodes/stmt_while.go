package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

type WhileStatementNode struct {
	ConditionStatementNode
}

func ParseWhileStatementNode(curAST *ast.AST, expectedEndType ...token.TokenType) (*WhileStatementNode, error) {
	node := &WhileStatementNode{}

	startToken, err := curAST.ConsumeToken(token.TokenType_WhileClause, token.TokenType_WhileClause.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	conditionTokens, err := curAST.ConsumeUntilFuncMatch(func(t *token.Token) bool {
		return t.Type == token.TokenType_Punctuation || t.Type == token.TokenType_KeywordThen
	}, "Expected token for statement condition ending")
	if err != nil {
		return nil, err
	}
	conditionNode, err := CreateValueNode(conditionTokens, CreateValueNodeOptions{})
	if err != nil {
		return nil, err
	}
	node.Condition = &conditionNode

	if curAST.CheckType(token.TokenType_KeywordThen) {
		curAST.Consume()
	}

	_, err = curAST.ConsumeFunc(func(t *token.Token) bool {
		return t.Type == token.TokenType_Punctuation
	}, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	statements, err := ParseStatementsNode(curAST, token.TokenType_KeywordStatementEnd)
	if err != nil {
		return nil, err
	}
	node.StatementsNode = *statements

	_, err = curAST.ConsumeToken(token.TokenType_KeywordStatementEnd, token.TokenType_KeywordStatementEnd.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	endToken, err := curAST.ConsumeToken(token.TokenType_Punctuation, token.TokenType_Punctuation.Message("Expected %s"))
	if err != nil {
		return nil, err
	}

	node.Start = startToken.Start
	node.Length = endToken.Start + endToken.Length - startToken.Start

	return node, nil
}
