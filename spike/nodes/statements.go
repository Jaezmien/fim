package nodes

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
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

func ParseStatementsNode(ast *ast.AST, expectedEndType ...token.TokenType) (*StatementsNode, error) {
	statements := &StatementsNode{}

	for {
		if ast.CheckType(expectedEndType...) {
			break
		}
		if ast.CheckType(token.TokenType_EndOfFile) {
			return nil, ast.CreateErrorFromToken(ast.Peek(), token.TokenType_FunctionFooter.Message("Could not find %s"))
		}

		if ast.CheckType(token.TokenType_NewLine) {
			continue
		}

		if ast.CheckType(token.TokenType_Print) || ast.CheckType(token.TokenType_PrintNewline) {
			node, err := ParsePrintNode(ast)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}
		if ast.CheckType(token.TokenType_Prompt)  {
			node, err := ParsePromptNode(ast)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if ast.CheckType(token.TokenType_Declaration) {
			node, err := ParseVariableDeclarationNode(ast)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
		}

		if ast.ContainsWithStop(token.TokenType_Modify, token.TokenType_EndOfFile, token.TokenType_Punctuation) {
			node, err := ParseVariableModifyNode(ast)
			if err != nil {
				return nil, err
			}

			statements.Statements = append(statements.Statements, node)
			continue
			
		}

		return nil, ast.CreateErrorFromToken(ast.Peek(), fmt.Sprintf("Unsupported statement token: %s", ast.Peek().Type))
	}

	return statements, nil
}
