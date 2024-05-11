package main

import "fmt"

type BlockchainIndex struct{}

func (index BlockchainIndex) StartProcessing() {
	fmt.Println("Starting processing ...")
}

func (index BlockchainIndex) StopProcessing() {
	fmt.Println("Stopping processing ...")
}
