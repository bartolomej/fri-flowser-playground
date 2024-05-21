package main

import (
	"fmt"
	"fri-flowser-playground/internal/project"
	"os"
)

func main() {
	p := project.New()

	err := p.Open("https://github.com/nvdtf/flow-nft-scaffold")

	if err != nil {
		fmt.Println(fmt.Errorf("error opening project: %v", err))
		os.Exit(1)
	}

	// Prevent the function from returning
	select {}
}
