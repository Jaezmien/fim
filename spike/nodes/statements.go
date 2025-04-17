package nodes

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type StatementsNode struct {
	Node

	Statements []INode
}

func (s *StatementsNode) Type() NodeType {
	return TYPE_STATEMENTS
}

func ParseStatementsNode(curAST *ast.AST, expectedEndType ...token.TokenType) (*StatementsNode, error) {
	statements := &StatementsNode{}

	for {
		if curAST.CheckType(expectedEndType...) {
			break
		}
		if curAST.CheckType(token.TokenType_EndOfFile) {
			return nil, curAST.Peek().CreateError(token.TokenType_FunctionFooter.Message("Could not find %s"), curAST.Source)
		}

		if curAST.CheckType(token.TokenType_NewLine) {
			continue
		}

		if curAST.CheckType(token.TokenType_Print) || curAST.CheckType(token.TokenType_PrintNewline) {
			node, err := ParsePrintNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}
		if curAST.CheckType(token.TokenType_Prompt) {
			node, err := ParsePromptNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.CheckType(token.TokenType_Declaration) {
			node, err := ParseVariableDeclarationNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.CheckType(token.TokenType_Identifier) && curAST.CheckNextType(token.TokenType_Modify) {
			node, err := ParseVariableModifyNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.CheckType(token.TokenType_FunctionCall) {
			node, err := ParseFunctionCallNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.CheckType(token.TokenType_IfClause) {
			node, err := ParseIfStatementsNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.CheckType(token.TokenType_UnaryIncrementPrefix, token.TokenType_UnaryDecrementPrefix) {
			if curAST.CheckNextType(token.TokenType_Identifier) {
				node, err := ParsePrefixUnary(curAST)
				if err != nil {
					return nil, err
				}

				statements.Statements = append(statements.Statements, node)
				continue
			}
		}

		if curAST.CheckType(token.TokenType_Identifier) {
			if curAST.CheckNextType(token.TokenType_UnaryIncrementPostfix, token.TokenType_UnaryDecrementPostfix) {
				node, err := ParsePostfixUnary(curAST)
				if err != nil {
					return nil, err
				}

				statements.Statements = append(statements.Statements, node)
				continue
			}
		}

		if curAST.CheckType(token.TokenType_KeywordReturn) {
			node, err := ParseFunctionReturnNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.ContainsWithStop(token.TokenType_KeywordOf, token.TokenType_EndOfFile, token.TokenType_NewLine) &&
			curAST.ContainsWithStop(token.TokenType_OperatorEq, token.TokenType_EndOfFile, token.TokenType_NewLine) {
			node, err := ParseArrayModifyNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		return nil, curAST.Peek().CreateError(fmt.Sprintf("Unsupported statement token: %s", curAST.Peek().Type), curAST.Source)
	}

	return statements, nil
}

// --- //

type ConditionStatementNode struct {
	StatementsNode

	Condition *INode
}
func (s *ConditionStatementNode) Type() NodeType {
	return TYPE_STATEMENTS_CONDITION
}
func (s *ConditionStatementNode) ToNode() Node {
	return Node{
		Start:  s.Start,
		Length: s.Length,
	}
}

// --- //

type IfStatementNode struct {
	Node

	Conditions []ConditionStatementNode
}

func (s *IfStatementNode) Type() NodeType {
	return TYPE_STATEMENTS_IF
}
func (s *IfStatementNode) ToNode() Node {
	return Node{
		Start:  s.Start,
		Length: s.Length,
	}
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
		Condition: &conditionNode,
		StatementsNode: *statements,
	})

	// Optional: ELSE clause
	for curAST.CheckType(token.TokenType_ElseClause) {
		_, err := curAST.ConsumeToken(token.TokenType_ElseClause, token.TokenType_ElseClause.Message("Expected %s"))
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
