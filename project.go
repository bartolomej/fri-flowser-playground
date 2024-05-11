package main

import "fmt"

type Project struct {
	emulator   FlowEmulator
	repository GitRepository
	index      BlockchainIndex
}

func (p Project) Open(projectUrl string) {
	fmt.Println("Opening project: " + projectUrl)
	p.repository.Clone(projectUrl)
	p.emulator.Start()
	p.index.StartProcessing()
}
