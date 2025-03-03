package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/spike/nodes"
	"git.jaezmien.com/Jaezmien/fim/twilight"
)

func main() {
	source :=
	`Dear Princess Celestia: Hello World!
	Today I learned how to say hello world!
	I said "Hello World"!
	That's all about how to say hello world.
	Your faithful student, Twilight Sparkle.
	`

	tokens := twilight.Parse(source)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	for idx := range tokens.Len() {
		token := tokens.Peek(idx).Value
		fmt.Fprintf(w, "%s\t\t%s\n", token.Type.String(), token.Value)
	}
	w.Flush()

	ast := spike.NewAST(tokens.Flatten(), source)
	report, err := nodes.ParseReportNode(ast)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", report)
}
