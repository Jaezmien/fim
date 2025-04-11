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
func (s *StatementsNode) ToNode() Node {
	return Node{
		Start:  s.Start,
		Length: s.Length,
	}
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

		if curAST.Peek().Type == token.TokenType_Identifier && curAST.PeekNext().Type == token.TokenType_Modify {
			node, err := ParseVariableModifyNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.Peek().Type == token.TokenType_FunctionCall {
			node, err := ParseFunctionCallNode(curAST)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if curAST.CheckType(token.TokenType_UnaryIncrementPrefix, token.TokenType_UnaryDecrementPrefix) {
			if curAST.PeekNext().Type == token.TokenType_Identifier {
				node, err := ParsePrefixUnary(curAST)
				if err != nil {
					return nil, err
				}

				statements.Statements = append(statements.Statements, node)
				continue
			}
		}

		if curAST.Peek().Type == token.TokenType_Identifier {
			if curAST.CheckNextType(token.TokenType_UnaryIncrementPostfix, token.TokenType_UnaryDecrementPostfix) {
				node, err := ParsePostfixUnary(curAST)
				if err != nil {
					return nil, err
				}

				statements.Statements = append(statements.Statements, node)
				continue
			}
		}

		if curAST.Peek().Type == token.TokenType_KeywordReturn {
			node, err := ParseFunctionReturnNode(curAST)
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
