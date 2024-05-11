package main

import "fmt"

type FlowEmulator struct{}

func (e FlowEmulator) Start() {
	fmt.Println("Starting emulator ...")
}

func (e FlowEmulator) Stop() {
	fmt.Println("Stopping emulator ...")
}
