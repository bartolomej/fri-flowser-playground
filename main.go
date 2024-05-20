package main

func main() {
	repository := &GitRepository{
		Test: "Hello",
	}
	emulator := &FlowEmulator{}
	index := &BlockchainIndex{}

	project := &Project{
		repository: repository,
		emulator:   emulator,
		index:      index,
	}

	project.Open("https://github.com/findonflow/find")
}
