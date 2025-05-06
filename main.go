//go:build !js

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"git.jaezmien.com/Jaezmien/fim/celestia"
	"git.jaezmien.com/Jaezmien/fim/luna/aprint"
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

	if *tokenDisplayFlag {
		epf := aprint.New(4, " ", aprint.LEFT_ALIGN)
		epf.SetAlignment(0, aprint.RIGHT_ALIGN)
		epf.SetAlignment(2, aprint.RIGHT_ALIGN)
		epf.SetDelimeter(2, " -> ")

		for idx, token := range tokens {
			epf.Add(
				strconv.FormatInt(int64(idx + 1), 10) + ".",
				strconv.FormatInt(int64(token.Start), 10) + ":" + 
				strconv.FormatInt(int64(token.Start + token.Length), 10),
				token.Value,
				token.Type.String(),
			)
		}

		fmt.Println(epf.String())

		return
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

	if *prettyFlag {
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
