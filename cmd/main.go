package main

import (
	"fmt"
	"fri-flowser-playground/internal/project"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	logger := initLogger()
	p := project.New(logger)

	err := p.Open("https://github.com/nvdtf/flow-nft-scaffold")

	if err != nil {
		fmt.Println(fmt.Errorf("error opening project: %v", err))
		os.Exit(1)
	}

	result, err := p.ExecuteScript([]byte("pub fun main(): Int { return 0 }"), "")

	if err != nil {
		fmt.Println(fmt.Errorf("error executing script: %v", err))
		os.Exit(1)
	}

	fmt.Println(result)

	result, err = p.ExecuteTransaction([]byte("transaction { }"), "")

	if err != nil {
		fmt.Println(fmt.Errorf("error executing transaction: %v", err))
		os.Exit(1)
	}

	fmt.Println(result)

	// Prevent the function from returning
	select {}
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
