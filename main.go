//go:build !js

package main

import (
	"flag"
	"fmt"
	"os"

	"git.jaezmien.com/Jaezmien/fim/celestia"
	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/twilight"
)

func main() {
	prettyFlag := flag.Bool("pretty", false, "Prettify output")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		return
	}

	filePath := args[0]
	if stat, err := os.Stat(filePath); err != nil || !stat.Mode().IsRegular() {
		fmt.Printf("Invalid file '%s'\n", filePath)
		return
	}
	
	rawSource, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("An error has occured while trying to load file '%s'\n", filePath)
		fmt.Printf("%s\n", err.Error())
		return
	}
	source := string(rawSource)

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

	if *prettyFlag == true {
		fmt.Printf("┌─ fim (v0.0.0-alpha)\n")
		fmt.Printf("├─ Report Name: %s\n", interpreter.ReportName())
		fmt.Printf("└─ Report Author: %s\n", interpreter.ReportAuthor())
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
