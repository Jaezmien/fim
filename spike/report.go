package spike

import (
	"git.jaezmien.com/Jaezmien/fim/spike/ast"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

// Create a ReportNode based on the generated token array.
func CreateReport(tokens []*token.Token, source string) (*nodes.ReportNode, error) {
	ast := &ast.AST{
		Tokens:     tokens,
		TokenIndex: 0,
		Source:     source,
	}

	return nodes.ParseReportNode(ast)
}
