package main

import (
	"fmt"
	"fri-flowser-playground/internal/project"
	"github.com/rs/zerolog"
	"os"
	"syscall/js"
)

var p *project.Project
var logger *zerolog.Logger

func main() {
	// Mount the function on the JavaScript global object.
	js.Global().Set("openProject", js.FuncOf(openProject))

	// Prevent the function from returning, which is required in a wasm module
	select {}
}

func openProject(this js.Value, args []js.Value) interface{} {
	logger = initLogger()
	p = project.New(logger)

	if len(args) != 1 {
		panic("Missing project URL arg")
	}

	err := p.Open(args[0].String())

	if err != nil {
		panic("error opening project: " + err.Error())
	}

	return nil
}

func initLogger() *zerolog.Logger {

	level := zerolog.InfoLevel
	zerolog.MessageFieldName = "msg"

	writer := zerolog.MultiLevelWriter(
		NewTextWriter(),
	)

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(level)

	return &logger
}

func NewTextWriter() zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	writer.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("%-44s", i)
	}

	return writer
}
