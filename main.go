//go:build !js

package main

import (
	"fmt"

	"git.jaezmien.com/Jaezmien/fim/celestia"
	"git.jaezmien.com/Jaezmien/fim/spike"
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

	report, err := spike.CreateReport(tokens.Flatten(), source)
	if err != nil {
		fmt.Println(err)
		return
	}

	interpreter, err := celestia.NewInterpreter(report, source)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, paragraph := range interpreter.Paragraphs {
		if paragraph.Main {
			if err := paragraph.Execute(); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
