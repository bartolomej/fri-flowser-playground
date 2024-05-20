package main

import (
	"fmt"
	"os"
)

func main() {
	repository := &GitRepository{}
	emulator := &FlowEmulator{}
	index := &BlockchainIndex{}

	project := &Project{
		repository: repository,
		emulator:   emulator,
		index:      index,
	}

	err := project.Open("https://github.com/nvdtf/flow-nft-scaffold")

	if err != nil {
		fmt.Println(fmt.Errorf("error opening project: %v", err))
		os.Exit(1)
	}

	// Prevent the function from returning
	select {}
}
