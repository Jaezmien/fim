//go:build js

package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"git.jaezmien.com/Jaezmien/fim/celestia"
	"git.jaezmien.com/Jaezmien/fim/spike"
	"git.jaezmien.com/Jaezmien/fim/twilight"
)

type CallbackWriter struct {
	Callback *js.Value
}

func NewCallbackWriter(callback *js.Value) (*CallbackWriter, error) {
	if callback.Type() != js.TypeFunction {
		return nil, errors.New("Expected callback function to be a JS function")
	}

	return &CallbackWriter{
		Callback: callback,
	}, nil
}
func (w *CallbackWriter) Write(p []byte) (n int, err error) {
	w.Callback.Invoke(string(p))
	return len(p), nil
}

func Exists(v js.Value) bool {
	return !v.IsNull() && !v.IsUndefined()
}

func main() {
	c := make(chan struct{}, 0)

	// fim(source: string) => [tokens: string[], error: string | null]
	js.Global().Set("fim", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return []any{nil, "Expected only one argument"}
		}
		if args[0].Type() != js.TypeString {
			return []any{nil, "Expected argument 0 to be type of string"}
		}

		source := args[0]

		tokens := twilight.Parse(source.String())

		result := make([]any, 0)
		for _, token := range tokens {
			if len(token.Value) > 0 {
				result = append(result, token.Value)
			}
		}

		return []any{result, nil}
	}))

	// fim_exec( source: string, output?: (data: string) => void, prompt?: (prompt: string) => string, error?: (info: string) => void) => error?: string
	js.Global().Set("fim_exec", js.FuncOf(func(this js.Value, args []js.Value) any {
		console := js.Global().Get("console")
		if !Exists(console) {
			return "Could not get console"
		}

		alert := js.Global().Get("alert")
		promptCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
			return ""
		}).Value
		if !Exists(alert) {
			promptCallback = alert
		}

		console_log := js.Global().Get("console").Get("log")
		outputCallback, err := NewCallbackWriter(&console_log)
		if err != nil {
			return "Could not get console.log"
		}

		console_error := js.Global().Get("console").Get("error")
		errorCallback, err := NewCallbackWriter(&console_error)
		if err != nil {
			return "Could not get console.error"
		}

		if len(args) < 1 {
			return "Expected at least one argument"
		}

		if args[0].Type() != js.TypeString {
			return "Expected argument 0 to be type of string"
		}
		source := args[0].String()

		if len(args) >= 2 && args[1].Type() == js.TypeFunction {
			outputCallback, err = NewCallbackWriter(&args[1])
			if err != nil {
				return err.Error()
			}
		}

		if len(args) >= 3 && args[2].Type() == js.TypeFunction {
			promptCallback = args[2]
		}

		if len(args) >= 4 && args[3].Type() == js.TypeFunction {
			errorCallback, err = NewCallbackWriter(&args[3])
			if err != nil {
				return err.Error()
			}
		}

		tokens := twilight.Parse(source)

		report, err := spike.CreateReport(tokens, source)
		if err != nil {
			fmt.Fprintln(errorCallback, err)
			return nil
		}

		interpreter, err := celestia.NewInterpreter(report, source)
		if err != nil {
			fmt.Fprintln(errorCallback, err)
			return nil
		}
		interpreter.Writer = outputCallback
		interpreter.ErrorWriter = errorCallback
		interpreter.Prompt = func(prompt string) (string, error) {
			result := promptCallback.Invoke(prompt)
			if result.Type() != js.TypeString {
				panic("PromptCallback returned a non-string value")
			}
			return result.String(), nil
		}

		for _, paragraph := range interpreter.Paragraphs {
			if paragraph.Main {
				if err := paragraph.Execute(); err != nil {
					fmt.Fprintln(errorCallback, err)
					return nil
				}
			}
		}

		return nil
	}))

	<-c
}
