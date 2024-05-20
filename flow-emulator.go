package main

import (
	"fmt"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/rs/zerolog"
	"os"
)

type FlowEmulator struct {
	logger     *zerolog.Logger
	blockchain *emulator.Blockchain
}

func (e *FlowEmulator) Start() error {
	logger := initLogger()
	blockchain, err := emulator.New(
		emulator.WithLogger(*logger),
	)

	if err != nil {
		return err
	}

	e.blockchain = blockchain
	e.logger = logger

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
