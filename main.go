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
	tokenDisplayFlag := flag.Bool("tokens", false, "Display tokens")
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

	if *tokenDisplayFlag == true {
		for idx, token := range tokens {
			fmt.Printf(
				"%d.\t%d:%d\t-> %s (%s)\n",
				idx,
				token.Start, token.Start+token.Length,
				token.Value,
				token.Type,
			)
		}
	}

	report, err := spike.CreateReport(tokens, source)
	if err != nil {
		fmt.Println("Spike noticed something unusual in your report...")
		fmt.Println(err)
		return
	}

	interpreter, err := celestia.NewInterpreter(report, source)
	if err != nil {
		fmt.Println("Princess Celestia noticed something unusual in your report...")
		fmt.Println(err)
		return
	}

	if *prettyFlag == true {
		fmt.Printf("┌─ fim (v0.0.0-alpha)\n")
		fmt.Printf("├─ Report Name: %s\n", interpreter.ReportTitle())
		fmt.Printf("└─ Report Author: %s\n", interpreter.ReportAuthor())
	}

	for _, paragraph := range interpreter.Paragraphs {
		if paragraph.Main {
			if _, err := paragraph.Execute(); err != nil {
				fmt.Println("Princess Celestia caught something unusual in your report!")
				fmt.Println(err)
				return
			}
		}
	}
}
