package main

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/twilight"
)

func main() {
	source :=
	`Dear Princess Celestia: Hello World!
	Your faithful student, Twilight Sparkle.
	`

	tokens := twilight.Parse(source)
	for tokens.Len() > 0 {
		fmt.Printf("%+v\n", tokens.Dequeue().Value)
	}
}
