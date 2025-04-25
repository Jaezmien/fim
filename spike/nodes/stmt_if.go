package nodes

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type IfStatementNode struct {
	Node

	Conditions []ConditionStatementNode
}

func ParseIfStatementsNode(curAST *ast.AST, expectedEndType ...token.TokenType) (*IfStatementNode, error) {
	node := &IfStatementNode{}
	node.Conditions = make([]ConditionStatementNode, 0)

	// Required: IF clause
	startToken, err := curAST.ConsumeToken(token.TokenType_IfClause, token.TokenType_IfClause.Message("Expected %s"))
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

	if curAST.CheckType(token.TokenType_KeywordThen) {
		curAST.ConsumeToken(token.TokenType_KeywordThen, token.TokenType_KeywordThen.Message("Expected %s"))
	}
	_, err = curAST.ConsumeFunc(func(t *token.Token) bool {
		return t.Type == token.TokenType_Punctuation && t.Value == ","
	}, token.TokenType_Punctuation.Message("Expected %s (comma)"))
	if err != nil {
		return nil, err
	}

	statements, err := ParseStatementsNode(curAST, token.TokenType_ElseClause, token.TokenType_IfEndClause)
	if err != nil {
		return nil, err
	}
	node.Conditions = append(node.Conditions, ConditionStatementNode{
		Condition:      &conditionNode,
		StatementsNode: *statements,
	})
	hasElseClause := false

	// Optional: ELSE clause
	for curAST.CheckType(token.TokenType_ElseClause) {
		elseToken, err := curAST.ConsumeToken(token.TokenType_ElseClause, token.TokenType_ElseClause.Message("Expected %s"))
		if err != nil {
			return nil, err
		}

		clause := ConditionStatementNode{
			Condition: nil,
		}

		conditionTokens, err := curAST.ConsumeUntilFuncMatch(func(t *token.Token) bool {
			return t.Type == token.TokenType_Punctuation || t.Type == token.TokenType_KeywordThen
		}, "Expected token for statement condition ending")
		if err != nil {
			return nil, err
		}

		if len(conditionTokens) > 0 {
			conditionNode, err := CreateValueNode(conditionTokens, CreateValueNodeOptions{})
			if err != nil {
				return nil, err
			}

			clause.Condition = &conditionNode
		} else {
			if hasElseClause {
				return nil, elseToken.CreateError("Else condition already exists", curAST.Source)
			}

			hasElseClause = true
		}

		if curAST.CheckType(token.TokenType_KeywordThen) {
			curAST.ConsumeToken(token.TokenType_KeywordThen, token.TokenType_KeywordThen.Message("Expected %s"))
		}
		_, err = curAST.ConsumeFunc(func(t *token.Token) bool {
			return t.Type == token.TokenType_Punctuation && t.Value == ","
		}, token.TokenType_Punctuation.Message("Expected %s (comma)"))
		if err != nil {
			return nil, err
		}

		statements, err := ParseStatementsNode(curAST, token.TokenType_ElseClause, token.TokenType_IfEndClause)
		if err != nil {
			return nil, err
		}
		clause.StatementsNode = *statements

		node.Conditions = append(node.Conditions, clause)
	}

	// Required: IF_END clause
	_, err = curAST.ConsumeToken(token.TokenType_IfEndClause, token.TokenType_IfEndClause.Message("Expected %s"))
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
