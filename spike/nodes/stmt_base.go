package nodes

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

type StatementsNode struct {
	Node

	Statements []DynamicNode
}

func ParseStatementsNode(curAST *ast.AST, expectedEndType ...token.TokenType) (*StatementsNode, error) {
	statements := &StatementsNode{}

	checks := []struct {
		Check  func() bool
		Parser func(ast *ast.AST) (DynamicNode, error)
	}{
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_Print) || curAST.CheckType(token.TokenType_PrintNewline)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParsePrintNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_Prompt)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParsePromptNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_Declaration)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseVariableDeclarationNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_Identifier) && curAST.CheckNextType(token.TokenType_Modify)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseVariableModifyNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_FunctionCall)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseFunctionCallNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_IfClause)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseIfStatementsNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_WhileClause)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseWhileStatementNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_UnaryIncrementPrefix, token.TokenType_UnaryDecrementPrefix) &&
					curAST.CheckNextType(token.TokenType_Identifier)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParsePrefixUnary(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_Identifier) &&
					curAST.CheckNextType(token.TokenType_UnaryIncrementPostfix, token.TokenType_UnaryDecrementPostfix)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParsePostfixUnary(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_KeywordReturn)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseFunctionReturnNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_ForEveryClause) &&
					curAST.ContainsWithStop(token.TokenType_KeywordIn, token.TokenType_EndOfFile, token.TokenType_NewLine)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseForEveryArrayStatementNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.CheckType(token.TokenType_ForEveryClause) &&
					curAST.ContainsWithStop(token.TokenType_KeywordFrom, token.TokenType_EndOfFile, token.TokenType_NewLine) &&
					curAST.ContainsWithStop(token.TokenType_KeywordTo, token.TokenType_EndOfFile, token.TokenType_NewLine)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseForEveryRangeStatementNode(ast)
			},
		},
		{
			Check: func() bool {
				return curAST.ContainsWithStop(token.TokenType_KeywordOf, token.TokenType_EndOfFile, token.TokenType_NewLine) &&
					curAST.ContainsWithStop(token.TokenType_OperatorEq, token.TokenType_EndOfFile, token.TokenType_NewLine)
			},
			Parser: func(ast *ast.AST) (DynamicNode, error) {
				return ParseArrayModifyNode(ast)
			},
		},
	}

	for {
		if curAST.CheckType(expectedEndType...) {
			break
		}
		if curAST.CheckType(token.TokenType_EndOfFile) {
			return nil, curAST.Peek().CreateError(token.TokenType_FunctionFooter.Message("Could not find %s"), curAST.Source)
		}

		if curAST.CheckType(token.TokenType_NewLine) {
			curAST.Consume()
			continue
		}

		if curAST.CheckType(token.TokenType_Punctuation) {
			curAST.Consume()
			continue
		}

		foundStatement := false
		for _, check := range checks {
			if check.Check() {
				node, err := check.Parser(curAST)
				if err != nil {
					return nil, err
				}
				statements.Statements = append(statements.Statements, node)
				foundStatement = true
				break
			}
		}

		if foundStatement {
			continue
		}

		return nil, curAST.Peek().CreateError(fmt.Sprintf("Unsupported statement token: %s", curAST.Peek().Type), curAST.Source)
	}

	return statements, nil
}

// --- //

type ConditionStatementNode struct {
	StatementsNode

	Condition *DynamicNode
}
